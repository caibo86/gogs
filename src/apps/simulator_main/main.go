// -------------------------------------------
// @file      : main.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/7 下午6:13
// -------------------------------------------

package main

import (
	"gogs/base/config"
	"gogs/base/etcd"
	"gogs/simulator"
)

func main() {
	config.ServerType = etcd.ServerTypeGame
	config.With(
		config.KeyLog,
		config.KeyRPC,
		config.KeySimulator,
	)
	config.LoadGlobalConfig("simulator.yml")
	simulator.Main()
}
