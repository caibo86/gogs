// -------------------------------------------
// @file      : redirect_err_unix.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/16 下午1:40
// -------------------------------------------

//go:build !windows

package logger

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"
)

func redirectStdErrLog() error {
	// dir, err := os.Stat("log")
	// if err != nil {
	// 	if os.IsNotExist(err) {
	// 		err = os.Mkdir("log", 0755)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	} else {
	// 		return err
	// 	}
	// } else {
	// 	if !dir.IsDir() {
	// 		return fmt.Errorf("log is not dir")
	// 	}
	// }
	panicFile := strings.Replace(global.fileName, ".log", ".panic", -1)
	fd, err := os.OpenFile(panicFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	err = syscall.Dup2(int(fd.Fd()), int(os.Stderr.Fd()))
	if err != nil {
		return err
	}
	Info("Redirect panic log successfully")
	// 保活文件 避免删除
	go func() {
		keepaliveFile := func() {
			CheckStdErrLogFile()
			_, err = fmt.Fprintln(os.Stderr, "no panic:"+time.Now().Format(time.RFC3339))
		}
		keepaliveFile()
		hourTimer := time.NewTicker(1 * time.Hour)
		defer hourTimer.Stop()
		for {
			select {
			case <-hourTimer.C:
				keepaliveFile()
				hourTimer.Reset(1 * time.Hour)
			}
		}
	}()
	return nil
}

func CheckStdErrLogFile() {
	panicFile := strings.Replace(global.fileName, ".log", ".panic", -1)
	_, err := os.Stat(panicFile)
	if !os.IsNotExist(err) {
		return
	}
	fd, err := os.OpenFile(panicFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	err = syscall.Dup2(int(fd.Fd()), int(os.Stderr.Fd()))
	if err != nil {
		return
	}
}
