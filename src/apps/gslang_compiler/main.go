// -------------------------------------------
// @file      : main.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/14 下午3:03
// -------------------------------------------

package main

import (
	"fmt"
	log "gogs/base/logger"
)

func main() {
	log.Init()
	fmt.Println(11)
	log.Panic("我是panic")
	fmt.Println(22)
}
