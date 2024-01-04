// -------------------------------------------
// @file      : signal.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午10:30
// -------------------------------------------

package gate

import (
	log "gogs/base/logger"
	"os"
	"os/signal"
	"syscall"
)

// ProcessSignal 监听系统信号 实现优雅退出
func ProcessSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		for {
			sig := <-ch
			switch sig {
			case syscall.SIGHUP:
				log.Warn("\033[043;1m[SIGHUP]\033[0m")
			case syscall.SIGTERM, syscall.SIGINT:
				appExit(sig)
			default:
				log.Warnf("unhandled signal:%v", sig)
			}
		}
	}()
}

// appExit 应用退出
func appExit(sig os.Signal) {
	log.Warnf("\033[043;1m[%v, quit]\033[0m", sig)
	exitChan <- struct{}{}
}
