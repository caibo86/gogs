// -------------------------------------------
// @file      : game.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午2:54
// -------------------------------------------

package gscluster

import (
	"sync"
)

// Game 游戏服务器
type Game struct {
	*RPC         // RPC管理器
	sync.RWMutex // 读写锁

}
