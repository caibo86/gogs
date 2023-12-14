// -------------------------------------------
// @file      : logger_test.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/14 下午1:39
// -------------------------------------------

package logger

import (
	"go.uber.org/zap"
	"runtime"
	"testing"
)

func BenchmarkInfo(b *testing.B) {
	defer func() {
		_ = Close()
	}()
	runtime.GOMAXPROCS(runtime.NumCPU())
	Init(
		SetIsAsync(true),
		SetIsOpenFile(true),
		SetIsOpenErrorFile(true),
	)
	for i := 1; i <= b.N; i++ {
		Info("我来测试一下Info", zap.Int("num", i))
	}
}
