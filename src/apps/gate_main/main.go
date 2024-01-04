// -------------------------------------------
// @file      : main.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午9:21
// -------------------------------------------

package main

import (
	"gogs/base/config"
	"gogs/base/etcd"
	"gogs/gate"
)

func main() {
	config.Preload()
	config.ServerType = etcd.ServerTypeGate
	config.With(
		config.KeyEtcd,
		config.KeyLog,
		config.KeyRPC,
		config.KeyGate,
	)
	config.LoadGlobalConfig("gate.yml")
	gate.Main()
}
