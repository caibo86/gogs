// -------------------------------------------
// @file      : host_agent.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 上午12:40
// -------------------------------------------

package cluster

import (
	"gogs/base/cluster/network"
	log "gogs/base/logger"
	"time"
)

// HostAgent 集群会话代理
// 实现了 network.ISessionHandler 接口和 IAgent 接口
// 将 Host 与 network.HostSession 连接起来 在两者中间实现消息传递
// 负责识别 network.Message 结构 根据不同Code调用上层函数
type HostAgent struct {
	session network.ISession // 远程服务 传输层通道
	Host    *Host            // 所属Host
}

// NewClusterRemote 新建集群远程服务
func NewClusterRemote(host *Host, session network.ISession) *HostAgent {
	return &HostAgent{
		session: session,
		Host:    host,
	}
}

// Name 其实就是session的名字
func (agent *HostAgent) Name() string {
	return agent.session.Name()
}

// Close implements IAgent
func (agent *HostAgent) Close() {
	agent.session.Close()
}

// Session implements IAgent
func (agent *HostAgent) Session() network.ISession {
	return agent.session
}

// Post implements IAgent
func (agent *HostAgent) Post(service IService, call *network.Call) error {
	return agent.Host.Post(agent.session, call)
}

// Wait implements IAgent
func (agent *HostAgent) Wait(service IService, call *network.Call, timeout time.Duration) (Future, error) {
	return agent.Host.Wait(agent.session, call, timeout)
}

// Write implements IAgent
func (agent *HostAgent) Write(msg *network.Message) error {
	return agent.session.Write(msg)
}

// SessionStatusChanged implements network.ISessionHandler
func (agent *HostAgent) SessionStatusChanged(status network.SessionStatus) {
	agent.Host.sessionStatusChanged(agent, status)
}

// Read implements network.ISessionHandler
func (agent *HostAgent) Read(session network.ISession, msg *network.Message) {
	switch msg.Type {
	case network.MessageTypeRegistry:
		// 处理来自邻居节点的服务状态变更通知
		go agent.handleServiceRegistry(msg.Data)
	case network.MessageTypeCall:
		go agent.handleCall(msg.Data)
	case network.MessageTypeReturn:
		go agent.handleReturn(msg.Data)
	}
}

// handleServiceRegistry 处理来自邻居节点的服务注册消息
func (agent *HostAgent) handleServiceRegistry(data []byte) {
	srd, err := network.UnmarshalServiceRegistryData(data)
	if err != nil {
		log.Warnf("unmarshal service registry data from %s err: %s", agent.session, err)
	}
	agent.Host.handleServiceRegistry(agent, srd)
}

// handleCall 处理来自对本地服务的调用
func (agent *HostAgent) handleCall(data []byte) {
	call, err := network.UnmarshalCall(data)
	if err != nil {
		log.Warnf("unmarshal call from %s err: %s", agent.session, err)
	}
	log.Infof("start handle rpc call userID: %d serviceID: %d methodID: %d from %s",
		call.ID, call.ServiceID, call.MethodID, agent.session)
	callReturn, err := agent.Host.handleCall(call)
	if err != nil {
		log.Warnf("handle rpc call userID: %d serviceID: %d methodID: %d from %s err: %s",
			call.ID, call.ServiceID, call.MethodID, agent.session, err)
	}
	if callReturn == nil {
		return
	}
	data = callReturn.Marshal()
	msg := &network.Message{
		Type: network.MessageTypeReturn,
		Data: data,
	}
	err = agent.session.Write(msg)
	if err != nil {
		log.Warnf("handle rpc call userID: %d serviceID: %d methodID: %d from %s err: %s",
			call.ID, call.ServiceID, call.MethodID, agent.session, err)
	}
	log.Infof("finish handle rpc call userID: %d serviceID: %d methodID: %d from %s",
		call.ID, call.ServiceID, call.MethodID, agent.session)
}

// handleReturn 处理对远程服务的调用返回
func (agent *HostAgent) handleReturn(data []byte) {
	callReturn, err := network.UnmarshalReturn(data)
	if err != nil {
		log.Warn("%s read return err: %s", agent.session, err)
		return
	}
	agent.Host.Notify(callReturn)
}
