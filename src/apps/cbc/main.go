// -------------------------------------------
// @file      : main.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/16 下午3:03
// -------------------------------------------

package main

import (
	"fmt"
	log "gogs/base/logger"
)

func main() {
	defer func() {
		fmt.Println(log.Close())
	}()
	log.Init(
		log.SetIsOpenFile(true),
		log.SetFilename("./log/cbc.log"),
	)
	fmt.Println(11)
	log.Panic("我是panic")
	fmt.Println(22)
}
