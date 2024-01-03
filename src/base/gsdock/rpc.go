// -------------------------------------------
// @file      : rpc.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/3 下午8:31
// -------------------------------------------

package gsdock

import (
	"gogs/base/gsnet"
	"sync/atomic"
	"time"
)

// rpcMonitor 远程调用返回值监控器
type rpcMonitor struct {
	future Future
	timer  *time.Timer
}

// rpcService 远程调用服务
type rpcService struct {
	idgen    uint32                 // id生成器
	monitors map[uint32]*rpcMonitor // 返回值监控器
}

// newRPCService 新建远程调用服务
func newRPCService() *rpcService {
	return &rpcService{
		monitors: make(map[uint32]*rpcMonitor),
	}
}

// post 远程调用,无返回值
func (rpc *rpcService) post(session gsnet.ISession, call *gsnet.Call) error {
	call.ID = atomic.AddUint32(&rpc.idgen, 1)
	data := call.Marshal()
	msg := &gsnet.Message{
		Type: gsnet.MessageTypeCall,
		Data: data,
	}
	if err := session.Write(msg); err != nil {
		return err
	}
	return nil
}

// wait 远程调用,有返回值,使用监控器处理超时
func (rpc *rpcService) wait(session gsnet.ISession, call *gsnet.Call, timeout time.Duration) (Future, error) {
	monitor := &rpcMonitor{
		future: make(Future, 1),
	}
	return monitor.future, nil
}
