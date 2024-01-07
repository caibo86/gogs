// -------------------------------------------
// @file      : builders.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午11:55
// -------------------------------------------

package game

import (
	"gogs/base/cluster"
	"gogs/cb"
)

var builders map[string]cluster.IServiceBuilder

func init() {
	builders = make(map[string]cluster.IServiceBuilder)
}

func RegisterBuilders() {
	builders["UserAPI"] = cb.NewUserBuilder(func(service cluster.IService) (cb.IUser, error) {
		return NewUser(service.Context().(*cluster.ClientAgent))
	})
}
