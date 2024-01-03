// -------------------------------------------
// @file      : gsdocker.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/29 下午6:28
// -------------------------------------------

package gsdock

import (
	"sync"
)

// GSDock 集群服务节点
type GSDock struct {
	sync.WaitGroup
	*RPC
}
