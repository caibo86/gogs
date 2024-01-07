// -------------------------------------------
// @file      : builders.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午11:55
// -------------------------------------------

package game

import (
	"gogs/base/cluster"
	"gogs/idl"
)

var builders map[string]cluster.IServiceBuilder

func init() {
	builders = make(map[string]cluster.IServiceBuilder)
}

func RegisterBuilders() {
	builders["UserAPI"] = idl.NewUserBuilder(func(service cluster.IService) (idl.IUser, error) {
		return NewUser(service.Context().(*cluster.ClientAgent))
	})
}
