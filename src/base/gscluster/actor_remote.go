// -------------------------------------------
// @file      : actor_remote.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午4:22
// -------------------------------------------

package gscluster

import (
	"gogs/base/gserrors"
	"gogs/base/gsnet"
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

// Post implement IRemote
func (remote *ActorRemote) Post(service IService, call *gsnet.Call) error {
	return remote.system.Post(remote, call)
}

// Wait implement IRemote
func (remote *ActorRemote) Wait(service IService, call *gsnet.Call, timeout time.Duration) (Future, error) {
	return remote.system.Wait(remote, call, timeout)
}

// Write implement IRemote
func (remote *ActorRemote) Write(msg *gsnet.Message) error {
	if msg.Type == gsnet.MessageTypeCall {
		return gserrors.Newf("check actor remote service implement")
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

// Close implement IRemote
func (remote *ActorRemote) Close() {
}

// Session implement IRemote
func (remote *ActorRemote) Session() gsnet.ISession {
	gserrors.Panic("not here")
	return nil
}

// Status implement gsnet.ISession
func (remote *ActorRemote) Status() gsnet.SessionStatus {
	return gsnet.SessionStatusInConnected
}

// DriverType implement gsnet.ISession
func (remote *ActorRemote) DriverType() gsnet.DriverType {
	return gsnet.DriverTypeActor
}

// Handler implement gsnet.ISession
func (remote *ActorRemote) Handler() gsnet.ISessionHandler {
	return nil
}

// RemoteAddr implement gsnet.ISession
func (remote *ActorRemote) RemoteAddr() string {
	return ""
}
