// -------------------------------------------
// @file      : service.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/29 下午6:23
// -------------------------------------------

package cluster

import (
	"gogs/base/cluster/network"
	"time"
)

// ID 服务ID.自增
type ID uint32

// IService 服务
type IService interface {
	Type() string                                     // 获取服务类型
	Name() string                                     // 获取服务名称
	ID() ID                                           // 服务ID
	Call(call *network.Call) (*network.Return, error) // 调用服务
	Context() interface{}                             // 服务上下文
}

// IAgent 会话代理
type IAgent interface {
	Post(service IService, call *network.Call) error                                  // 远程调用
	Wait(service IService, call *network.Call, timeout time.Duration) (Future, error) // 远程调用,需要返回结果
	Write(msg *network.Message) error                                                 // 写入消息
	Session() network.ISession                                                        // 代理的会话
	Close()                                                                           // 关闭
}

// IRemoteService 远程服务
type IRemoteService interface {
	IService
	RemoteID() ID  // 远程服务在其本地的ID
	Agent() IAgent // 获取Remote
}

// IServiceBuilder 服务类型builder
type IServiceBuilder interface {
	ServiceType() string                                                                             // 服务类型,service.typename
	NewService(name string, id ID, context interface{}) (IService, error)                            // 新建本地服务
	NewRemoteService(remote IAgent, name string, lid ID, rid ID, context interface{}) IRemoteService // 新建远程服务
}

// ReturnVal RPC调用返回值结构
type ReturnVal struct {
	Timeout    bool            // 结果是否超时
	CallReturn *network.Return // 调用结果
}

// Future RPC调用结果返回chan
type Future chan *ReturnVal
