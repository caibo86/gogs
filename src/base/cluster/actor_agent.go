// -------------------------------------------
// @file      : actor_agent.go
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

type ActorAgent struct {
	system   *ActorSystem // 本地角色系统引用
	name     string       // 角色名字
	neighbor IActorSystem
}

// newActorAgent 新建角色会话代理
func newActorAgent(system *ActorSystem, name string, neighbor IActorSystem) *ActorAgent {
	return &ActorAgent{
		system:   system,
		name:     name,
		neighbor: neighbor,
	}
}

// Name 获取角色名字
func (agent *ActorAgent) Name() string {
	return agent.name
}

// Post implement IAgent
func (agent *ActorAgent) Post(service IService, call *network.Call) error {
	return agent.system.Post(agent, call)
}

// Wait implement IAgent
func (agent *ActorAgent) Wait(service IService, call *network.Call, timeout time.Duration) (Future, error) {
	return agent.system.Wait(agent, call, timeout)
}

// Write implement IAgent
func (agent *ActorAgent) Write(msg *network.Message) error {
	if msg.Type == network.MessageTypeCall {
		return cberrors.New("check actor remote service implement")
	}
	actorMsg := &ActorMsg{
		ActorName: agent.name,
		Data:      msg.Data,
	}
	callReturn, status, err := agent.neighbor.ActorInvoke(actorMsg)
	if err != nil {
		return err
	}
	if err != ErrOK {
		return status
	}
	if len(callReturn.Params) != 0 {
		agent.system.Notify(callReturn)
	}
	return nil
}

// Close implement IAgent
func (agent *ActorAgent) Close() {
}

// Session implement IAgent
func (agent *ActorAgent) Session() network.ISession {
	cberrors.Panic("not here")
	return nil
}

// Status implement network.ISession
func (agent *ActorAgent) Status() network.SessionStatus {
	return network.SessionStatusInConnected
}

// DriverType implement network.ISession
func (agent *ActorAgent) DriverType() network.DriverType {
	return network.DriverTypeActor
}

// Handler implement network.ISession
func (agent *ActorAgent) Handler() network.ISessionHandler {
	return nil
}

// RemoteAddr implement network.ISession
func (agent *ActorAgent) RemoteAddr() string {
	return ""
}
