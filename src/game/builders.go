// -------------------------------------------
// @file      : builders.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午11:55
// -------------------------------------------

package game

import (
	"gogs/base/gscluster"
	"gogs/game/model"
	"gogs/idl"
)

var builders map[string]gscluster.IServiceBuilder

func init() {
	builders = make(map[string]gscluster.IServiceBuilder)
}

func RegisterBuilders() {
	builders["UserAPI"] = idl.NewPlayerBuilder(func(service gscluster.IService) (idl.IPlayer, error) {
		return model.NewUser(service.Context().(*gscluster.ClientAgent))
	})
}
