// -------------------------------------------
// @file      : builders.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午11:55
// -------------------------------------------

package login

import (
	"gogs/base/cluster"
)

var builders map[string]cluster.IServiceBuilder

func init() {
	builders = make(map[string]cluster.IServiceBuilder)
}

func RegisterBuilders() {
	builders =
}
