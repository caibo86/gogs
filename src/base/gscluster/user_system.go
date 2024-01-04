// -------------------------------------------
// @file      : user_system.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午2:56
// -------------------------------------------

package gscluster

import (
	"gogs/base/config"
	"sync"
)

// UserSystem 用户系统
type UserSystem struct {
	*RPC                                 // RPC管理器
	name         string                  // 名字
	host         *Host                   // 集群服务器
	builders     map[string]ITypeBuilder // 服务建造者集合
	actors       map[string]IUser        // 用户集合
	actorLock    sync.RWMutex            // 用户集合读写锁
	neighbors    map[string]IUserSystem  // 邻居用户系统集合
	neighborLock sync.RWMutex            // 邻居用户系统集合读写锁
	idgen        int64                   // id generator
	groupLocks   []sync.Mutex            // 分组锁
}

// NewUserSystem 新建角色系统
func NewUserSystem(name string, builders map[string]ITypeBuilder, localAddr string) (*UserSystem, error) {
	system := &UserSystem{
		RPC:        NewRPC(),
		name:       name,
		host:       NewHost(localAddr),
		builders:   builders,
		actors:     make(map[string]IUser),
		neighbors:  make(map[string]IUserSystem),
		groupLocks: make([]sync.Mutex, config.ActorLocks()),
	}
	system.host.RegisterBuilder(NewAc)
}
