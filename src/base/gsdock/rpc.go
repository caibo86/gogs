// -------------------------------------------
// @file      : rpc.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/3 下午8:31
// -------------------------------------------

package gsdock

import (
	"gogs/base/gsnet"
	"runtime"
	"sync"
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
func (rpc *rpcService) wait(lock *sync.Mutex, session gsnet.ISession, call *gsnet.Call, timeout time.Duration) (Future, error) {
	monitor := &rpcMonitor{
		future: make(Future, 1),
	}
	call.ID = atomic.AddUint32(&rpc.idgen, 1)
	lock.Lock()
	rpc.monitors[call.ID] = monitor
	lock.Unlock()
	id := call.ID
	monitor.timer = time.AfterFunc(timeout, func() {
		lock.Lock()
		defer lock.Unlock()
		// 调用超时
		if monitor1, ok := rpc.monitors[id]; ok {
			delete(rpc.monitors, id)
			monitor1.future <- &ReturnVal{
				Timeout: true,
			}
		}
	})
	data := call.Marshal()
	msg := &gsnet.Message{
		Type: gsnet.MessageTypeCall,
		Data: data,
	}
	if err := session.Write(msg); err != nil {
		return nil, err
	}
	return monitor.future, nil
}

// notify 异步调用的返回通知,找到对应的监控器,将结果写入监控器的结果通道中
func (rpc *rpcService) notify(lock *sync.Mutex, callReturn *gsnet.Return) bool {
	lock.Lock()
	defer lock.Unlock()
	// 查找对应id的结果监控器,如果存在则将结果写入监控器的结果通道中,非超时
	if monitor, ok := rpc.monitors[callReturn.ID]; ok {
		delete(rpc.monitors, callReturn.ID)
		monitor.timer.Stop()
		monitor.future <- &ReturnVal{
			CallReturn: callReturn,
		}
		return true
	}
	return false
}

// RPC 远程调用集中管理器
type RPC struct {
	locks []sync.Mutex  // 预分配的互斥锁列表,与rpc服务器一一对应
	group []*rpcService // rpc服务器列表,多个服务按照其id取模后取对应的rpc服务器
}

// NewRPC 新建远程调用集中管理器
func NewRPC() *RPC {
	// 单个管理器下的远程调用服务数量等于CPU数量的8倍
	groups := runtime.NumCPU() * 8
	rpc := &RPC{
		locks: make([]sync.Mutex, groups),
		group: make([]*rpcService, groups),
	}
	for i := 0; i < groups; i++ {
		rpc.group[i] = newRPCService()
	}
	return rpc
}

// Post 简单取模hash,获取对应的rpc服务器,并调用其post方法
func (rpc *RPC) Post(session gsnet.ISession, call *gsnet.Call) error {
	group := ID(call.ServiceID) % ID(len(rpc.locks))
	return rpc.group[group].post(session, call)
}

// Wait 简单取模hash,获取对应的rpc服务器,并调用其wait方法
func (rpc *RPC) Wait(session gsnet.ISession, call *gsnet.Call, timeout time.Duration) (Future, error) {
	group := ID(call.ServiceID) % ID(len(rpc.locks))
	return rpc.group[group].wait(&rpc.locks[group], session, call, timeout)
}

// Notify 简单取模hash,获取对应的rpc服务器,并调用其notify方法
func (rpc *RPC) Notify(callReturn *gsnet.Return) bool {
	group := ID(callReturn.ServiceID) % ID(len(rpc.locks))
	return rpc.group[group].notify(&rpc.locks[group], callReturn)
}
