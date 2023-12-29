// -------------------------------------------
// @file      : option.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/16 下午1:40
// -------------------------------------------

package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gogs/base/gserrors"
	"os"
	"strings"
)

// Option 日志配置项
type Option func(options *Options)

// Options 日志配置
type Options struct {
	Filename        string        // 日志文件路径
	Level           zapcore.Level // 日志级别
	MaxFileSize     int           // 日志分割的尺寸
	MaxAge          int           // 日志保存的时间 天
	MaxBackups      int           // 最大日志数量
	Stacktrace      zapcore.Level // 记录堆栈的日志级别
	FormatType      string        // 日志格式
	CallerSkip      int           // 堆栈的跳过层数
	IsAsync         bool          // 异步日志
	IsCompress      bool          // 是否压缩
	IsOpenPprof     bool          // 是否打开pprof
	IsOpenConsole   bool          // 是否打开终端标准输出
	IsOpenFile      bool          // 是否打开文件日志
	IsOpenErrorFile bool          // 是否打开高级别错误文件日志
}

func SetFilename(filename string) Option {
	return func(options *Options) {
		if filename != "" {
			options.Filename = filename
		}
	}
}

func SetLevel(level zapcore.Level) Option {
	return func(options *Options) {
		options.Level = level
	}
}

func SetMaxFileSize(size int) Option {
	return func(options *Options) {
		if size > 0 {
			options.MaxFileSize = size
		}
	}
}

func SetMaxAge(age int) Option {
	return func(options *Options) {
		options.MaxAge = age
	}
}

func SetMaxBackups(backups int) Option {
	return func(options *Options) {
		options.MaxBackups = backups
	}
}

func SetStacktrace(level zapcore.Level) Option {
	return func(options *Options) {
		options.Stacktrace = level
	}
}

func SetIsOpenConsole(console bool) Option {
	return func(options *Options) {
		options.IsOpenConsole = console
	}
}

func SetFormatType(format string) Option {
	return func(options *Options) {
		options.FormatType = format
	}
}

func SetCallerSkip(callerSkip int) Option {
	return func(options *Options) {
		options.CallerSkip = callerSkip
	}
}

func SetIsAsync(async bool) Option {
	return func(options *Options) {
		options.IsAsync = async
	}
}

func SetIsCompress(compress bool) Option {
	return func(options *Options) {
		options.IsCompress = compress
	}
}

func SetIsOpenPprof(pprof bool) Option {
	return func(options *Options) {
		options.IsOpenPprof = pprof
	}
}

func SetIsOpenFile(file bool) Option {
	return func(options *Options) {
		options.IsOpenFile = file
	}
}

func SetIsOpenErrorFile(errorFile bool) Option {
	return func(options *Options) {
		options.IsOpenErrorFile = errorFile
	}
}

func (o *Options) getLogFilename() string {
	ret := "Unknown"
	arr := strings.Split(o.Filename, "/")
	if len(arr) <= 0 {
		return ret
	}
	ret = arr[len(arr)-1]
	if tmp := strings.Split(ret, "."); len(tmp) > 0 {
		return tmp[0]
	}
	return ret
}

func (o *Options) getConsoleCore() zapcore.Core {
	var consoleWS zapcore.WriteSyncer
	if o.IsAsync {
		consoleWS = &zapcore.BufferedWriteSyncer{
			WS:   zapcore.AddSync(os.Stdout),
			Size: AsyncChanSize,
		}
	} else {
		consoleWS = zapcore.AddSync(os.Stdout)
	}
	encoder := GetEncoder(o.FormatType, o.getLogFilename())
	var enabler zap.LevelEnablerFunc
	enabler = func(level zapcore.Level) bool {
		return level >= global.atom.Level()
	}
	return zapcore.NewCore(encoder, consoleWS, enabler)
}

func (o *Options) getFileCore() zapcore.Core {
	var fileLogger zapcore.WriteSyncer
	fileLogger = &FileLogger{
		Logger: &lumberjack.Logger{
			Filename:   o.Filename,
			MaxSize:    o.MaxFileSize,
			MaxBackups: o.MaxBackups,
			MaxAge:     o.MaxAge,
			Compress:   o.IsCompress,
		},
	}
	if o.IsAsync {
		fileLogger = &zapcore.BufferedWriteSyncer{
			WS:   fileLogger,
			Size: AsyncChanSize,
		}
	}
	encoder := GetEncoder(o.FormatType, o.getLogFilename())
	var enabler zap.LevelEnablerFunc
	enabler = func(level zapcore.Level) bool {
		return level >= global.atom.Level() // && level <= zap.ErrorLevel
	}
	return zapcore.NewCore(encoder, fileLogger, enabler)
}

func (o *Options) getErrorFileCore() zapcore.Core {
	// 高级别错误文件日志直接采用同步写
	errFilename := strings.Replace(o.Filename, ".log", ".err", 1)
	fileLogger := &FileLogger{
		Logger: &lumberjack.Logger{
			Filename:   errFilename,
			MaxSize:    o.MaxFileSize,
			MaxBackups: o.MaxBackups,
			MaxAge:     o.MaxAge,
			Compress:   o.IsCompress,
		},
	}
	encoder := GetEncoder(o.FormatType, o.getLogFilename())
	var enabler zap.LevelEnablerFunc
	enabler = func(level zapcore.Level) bool {
		return level >= zap.ErrorLevel
	}
	return zapcore.NewCore(encoder, fileLogger, enabler)
}

func (o *Options) getCore() zapcore.Core {
	var cores []zapcore.Core
	if o.IsOpenConsole {
		cores = append(cores, o.getConsoleCore())
	}
	if o.IsOpenFile {
		cores = append(cores, o.getFileCore())
	}
	if o.IsOpenErrorFile {
		cores = append(cores, o.getErrorFileCore())
	}
	if len(cores) == 0 {
		gserrors.Panic("At least one log output needs to be opened")
	}
	return zapcore.NewTee(cores...)
}

func (o *Options) GetZapLogger() *zap.SugaredLogger {
	var options []zap.Option
	options = append(options, zap.AddStacktrace(o.Stacktrace))
	options = append(options, zap.AddCaller())
	options = append(options, zap.AddCallerSkip(o.CallerSkip))

	logger := zap.New(o.getCore(), options...).Sugar()
	if logger == nil {
		gserrors.Panic("get zap logger failed.")
	}
	return logger
}

type FileLogger struct {
	*lumberjack.Logger
}

func (fl *FileLogger) Sync() error {
	if fl.Logger != nil {
		return fl.Logger.Close()
	}
	return nil
}
