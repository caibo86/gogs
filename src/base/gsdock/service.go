// -------------------------------------------
// @file      : service.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/29 下午6:23
// -------------------------------------------

package gsdock

import (
	"fmt"
	"gogs/base/gsnet"
	"time"
)

// ID 服务ID
type ID uint32

// IService 服务
type IService interface {
	fmt.Stringer
	Type() string                                 // 获取服务类型
	ID() ID                                       // 服务ID
	Call(call *gsnet.Call) (*gsnet.Return, error) // 调用服务
	Context() interface{}                         // 服务上下文
}

// IRemote 远程服务调用接口
type IRemote interface {
	Post(service IService, call *gsnet.Call) error                                  // 通知消息
	Wait(service IService, call *gsnet.Call, timeout time.Duration) (Future, error) // 请求消息
	Write(msg *gsnet.Message) error                                                 // 写入消息
	Session() gsnet.ISession                                                        // 句柄对应的通道
	Close()                                                                         // 关闭
}

// IRemoteService 远程服务
type IRemoteService interface {
	IService
	RemoteID() ID    // 远程服务在其本地的ID
	Remote() IRemote // 获取Remote
}

// ITypeBuilder 服务类型builder
type ITypeBuilder interface {
	fmt.Stringer
	NewService(name string, id ID, context interface{}) (IService, error)                             // 新建本地服务
	NewRemoteService(remote IRemote, name string, lid ID, rid ID, context interface{}) IRemoteService // 新建远程服务
}

// ReturnVal RPC调用返回值结构
type ReturnVal struct {
	Timeout    bool          // 结果是否超时
	CallReturn *gsnet.Return // 调用结果
}

// Future RPC调用结果返回chan
type Future chan *ReturnVal
