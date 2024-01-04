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
	"gogs/game"
)

func main() {
	config.Preload()
	config.ServerType = etcd.ServerTypeGame
	config.With(
		config.KeyEtcd,
		config.KeyLog,
		config.KeyRPC,
		config.KeyGame,
	)
	config.LoadGlobalConfig("game.yml")
	config.Adjust(
		config.SetEtcdServiceType(config.ServerType),
		config.SetEtcdServiceID(config.ServerID),
	)
	game.Main()
}
