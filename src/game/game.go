// -------------------------------------------
// @file      : game.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 上午11:13
// -------------------------------------------

package game

import (
	"gogs/base/gscluster"
)

// Game 游戏服务器
type Game struct {
	*gscluster.RPC // 远程
}
