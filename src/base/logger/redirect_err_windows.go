// -------------------------------------------
// @file      : redirect_err_windows.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/16 下午1:41
// -------------------------------------------

//go:build windows

package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/sys/windows"
)

func redirectStdErrLog() error {
	panicFile := strings.Replace(global.fileName, ".log", ".panic", -1)
	file, err := os.OpenFile(panicFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	if err := windows.SetStdHandle(windows.STD_ERROR_HANDLE, windows.Handle(file.Fd())); err != nil {
		return err
	}
	os.Stderr = file
	Debug("redirect str err log success")
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
	_, err := os.Stat(panicFile)
	if !os.IsNotExist(err) {
		return
	}
	_ = redirectStdErrLog()
}
