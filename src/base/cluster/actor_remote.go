// -------------------------------------------
// @file      : actor_remote.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午4:22
// -------------------------------------------

package cluster

import (
	"gogs/base/cberrors"
	"gogs/base/cluster/network"
	"time"
)

type ActorRemote struct {
	system   *ActorSystem // 本地角色系统引用
	name     string       // 角色名字
	neighbor IActorSystem
}

// newActorRemote 新建角色远程代理
func newActorRemote(system *ActorSystem, name string, neighbor IActorSystem) *ActorRemote {
	return &ActorRemote{
		system:   system,
		name:     name,
		neighbor: neighbor,
	}
}

// Name 获取角色名字
func (remote *ActorRemote) Name() string {
	return remote.name
}

// Post implement IAgent
func (remote *ActorRemote) Post(service IService, call *network.Call) error {
	return remote.system.Post(remote, call)
}

// Wait implement IAgent
func (remote *ActorRemote) Wait(service IService, call *network.Call, timeout time.Duration) (Future, error) {
	return remote.system.Wait(remote, call, timeout)
}

// Write implement IAgent
func (remote *ActorRemote) Write(msg *network.Message) error {
	if msg.Type == network.MessageTypeCall {
		return cberrors.New("check actor remote service implement")
	}
	actorMsg := &ActorMsg{
		ActorName: remote.name,
		Data:      msg.Data,
	}
	callReturn, status, err := remote.neighbor.ActorInvoke(actorMsg)
	if err != nil {
		return err
	}
	if err != ErrOK {
		return status
	}
	if len(callReturn.Params) != 0 {
		remote.system.Notify(callReturn)
	}
	return nil
}

// Close implement IAgent
func (remote *ActorRemote) Close() {
}

// Session implement IAgent
func (remote *ActorRemote) Session() network.ISession {
	cberrors.Panic("not here")
	return nil
}

// Status implement network.ISession
func (remote *ActorRemote) Status() network.SessionStatus {
	return network.SessionStatusInConnected
}

// DriverType implement network.ISession
func (remote *ActorRemote) DriverType() network.DriverType {
	return network.DriverTypeActor
}

// Handler implement network.ISession
func (remote *ActorRemote) Handler() network.ISessionHandler {
	return nil
}

// RemoteAddr implement network.ISession
func (remote *ActorRemote) RemoteAddr() string {
	return ""
}
