// -------------------------------------------
// @file      : normal.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/8 下午5:17
// -------------------------------------------

package cluster

import (
	"gogs/base/etcd"
	log "gogs/base/logger"
	"sync"
	"sync/atomic"
)

// Normal 普通功能服务器
type Normal struct {
	*RPC                                    // RPC管理器
	sync.RWMutex                            // 读写锁
	Host         *Host                      // 集群服务器
	builders     map[string]IServiceBuilder // 服务构造器集合
	idgen        uint32                     // service userID generator
	serverName   string                     // 服务器名字
}

// NewNormal 新建普通功能服务器
func NewNormal(name string, builders map[string]IServiceBuilder, localAddr string) *Normal {
	normal := &Normal{
		RPC:        NewRPC(),
		Host:       NewHost(localAddr),
		builders:   builders,
		serverName: name,
	}
	return normal
}

// Shutdown 关闭服务器
func (normal *Normal) Shutdown() {
	log.Infof("%s shutdown start:", normal.serverName)
	log.Infof("%s:Host closing...", normal.serverName)
	normal.Host.Close()
	log.Infof("%s:Etcd closing....", normal.serverName)
	etcd.Exit()
	log.Infof("%s shutdown finished.", normal.serverName)
}

// newServiceID 生成一个唯一ID
func (normal *Normal) newServiceID() ID {
	return ID(atomic.AddUint32(&normal.idgen, 1))
}

// Name 服务器名字
func (normal *Normal) Name() string {
	return normal.serverName
}
