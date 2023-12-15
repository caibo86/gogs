// -------------------------------------------
// @file      : default.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/16 下午1:38
// -------------------------------------------

package logger

import "go.uber.org/zap"

const (
	LogFormatJson    = "json"
	LogFormatConsole = "console"
)

const (
	AsyncChanSize = 4096
)

var DefaultOptions = Options{
	Filename:        "./log/default.log",
	Level:           zap.DebugLevel,
	MaxFileSize:     128,
	MaxAge:          30,
	MaxBackups:      1024,
	Stacktrace:      zap.ErrorLevel,
	FormatType:      "console",
	CallerSkip:      1,
	IsAsync:         false,
	IsCompress:      true,
	IsOpenPprof:     false,
	IsOpenConsole:   true,
	IsOpenFile:      false,
	IsOpenErrorFile: false,
}
