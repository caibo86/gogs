// -------------------------------------------
// @file      : logger.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/16 下午1:39
// -------------------------------------------

package logger

import (
	"sync"
	"time"

	"gogs/base/cberrors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var global *Logger
var once sync.Once

func init() {
	global = &Logger{}
}

// Logger 日志封装
type Logger struct {
	zapLogger *zap.SugaredLogger
	atom      zap.AtomicLevel
	fileName  string
	options   *Options
}

func (logger *Logger) Close() error {
	if logger.zapLogger == nil {
		return nil
	}
	return logger.zapLogger.Sync()
}

// Init 日志初始化
func (logger *Logger) Init(option ...Option) {
	options := DefaultOptions
	for _, o := range option {
		o(&options)
	}
	logger.options = &options
	logger.fileName = options.Filename
	logger.atom = zap.NewAtomicLevel()
	logger.atom.SetLevel(options.Level)
	logger.zapLogger = options.GetZapLogger()
	if options.IsOpenPprof {
		StartPprofTask()
	}
}

// Init 全局日志初始化 每个app必须调用一次
func Init(options ...Option) {
	once.Do(func() {
		global.Init(options...)
	})
	// 触发创建目录
	Info("Init logger successfully")
	err := redirectStdErrLog()
	if err != nil {
		Errorf("Redirect panic log err: %s", err)
	}
}

// Close 关闭服务时调用
func Close() error {
	if global != nil {
		err := global.Close()
		if err.Error() == "sync /dev/stdout: invalid argument" ||
			err.Error() == "sync /dev/stdout: The handle is invalid." {
			return nil
		}
		return err
	}
	return nil
}

func Sync() {
	_ = global.zapLogger.Sync()
}

// GetEncoder 获取encoder
func GetEncoder(t, filename string) zapcore.Encoder {
	rfc339UTC := func(t time.Time, pe zapcore.PrimitiveArrayEncoder) {
		t = t.UTC()
		format := "2006-01-02 15:04:05.000"
		type appendTimeEncoder interface {
			AppendTimeLayout(time.Time, string)
		}
		if encoder, ok := pe.(appendTimeEncoder); ok {
			encoder.AppendTimeLayout(t, format)
			return
		}
		if filename != "" {
			pe.AppendString(t.Format(format) + "\t" + filename)
			return
		}
		pe.AppendString(t.Format(format))
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = rfc339UTC // TODO

	if t == LogFormatJson {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// GetZapLogger 获取zap logger
func GetZapLogger(core zapcore.Core, callerSkip int) *zap.SugaredLogger {
	var options []zap.Option
	options = append(options, zap.AddStacktrace(zapcore.ErrorLevel))
	options = append(options, zap.AddCaller())
	options = append(options, zap.AddCallerSkip(callerSkip))

	logger := zap.New(core, options...).Sugar()
	if logger == nil {
		cberrors.Panic("get zap logger failed.")
	}
	return logger
}
