// -------------------------------------------
// @file      : gate_api.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/16 下午1:37
// -------------------------------------------

package logger

func Debug(args ...interface{}) {
	if global.zapLogger != nil {
		global.zapLogger.Debug(args...)
	}

}

func Debugf(template string, args ...interface{}) {
	if global.zapLogger != nil {
		global.zapLogger.Debugf(template, args...)
	}

}

func Debugw(msg string, keysAndValues ...interface{}) {
	if global.zapLogger != nil {
		global.zapLogger.Debugw(msg, keysAndValues...)
	}
}

func Info(args ...interface{}) {
	if global.zapLogger != nil {
		global.zapLogger.Info(args...)
	}
}

func Infof(template string, args ...interface{}) {
	if global.zapLogger != nil {
		global.zapLogger.Infof(template, args...)
	}
}

func Infow(msg string, keysAndValues ...interface{}) {
	if global.zapLogger != nil {
		global.zapLogger.Infow(msg, keysAndValues...)
	}
}

func Warn(args ...interface{}) {
	if global.zapLogger != nil {
		global.zapLogger.Warn(args...)
	}

}

func Warnf(template string, args ...interface{}) {
	if global.zapLogger != nil {
		global.zapLogger.Warnf(template, args...)
	}
}

func Warnw(msg string, keysAndValues ...interface{}) {
	if global.zapLogger != nil {
		global.zapLogger.Warnw(msg, keysAndValues...)
	}
}

func Error(args ...interface{}) {
	if global.zapLogger != nil {
		global.zapLogger.Error(args...)
	}

}

func Errorf(template string, args ...interface{}) {
	if global.zapLogger != nil {
		global.zapLogger.Errorf(template, args...)
	}
}

func Errorw(msg string, keysAndValues ...interface{}) {
	if global.zapLogger != nil {
		global.zapLogger.Errorw(msg, keysAndValues...)
	}
}

func Panic(args ...interface{}) {
	if global.zapLogger != nil {
		global.zapLogger.Panic(args...)
	}
}

func Panicf(template string, args ...interface{}) {
	if global.zapLogger != nil {
		global.zapLogger.Panicf(template, args...)
	}
}

func Panicw(msg string, keysAndValues ...interface{}) {
	if global.zapLogger != nil {
		global.zapLogger.Panicw(msg, keysAndValues...)
	}
}

func Fatal(args ...interface{}) {
	if global.zapLogger != nil {
		global.zapLogger.Fatal(args...)
	}

}

func Fatalf(template string, args ...interface{}) {
	if global.zapLogger != nil {
		global.zapLogger.Fatalf(template, args...)
	}
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	if global.zapLogger != nil {
		global.zapLogger.Fatalw(msg, keysAndValues...)
	}
}

// func InfoMsg(level zapcore.Level, msg proto.Message, template string, args ...interface{}) {
// 	if level < global.atom.Level() {
// 		return
// 	}
// 	var buf strings.Builder
// 	buf.WriteString(fmt.Sprintf(template, args...))
// 	buf.WriteString(" ")
// 	jm := &jsonpb.Marshaler{}
// 	if err := jm.Marshal(&buf, msg); err != nil {
// 		Errorf("json pb marshal err:%s", err)
// 		return
// 	}
// 	global.zapLogger.Infof(buf.String())
// }
//
