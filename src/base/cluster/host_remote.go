// -------------------------------------------
// @file      : cluster.go
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
// 将 Host 与 ClusterSession 连接起来 在两者中间实现消息传递
// 负责识别Message结构 根据不同Code调用上层函数
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
func (remote *HostAgent) Name() string {
	return remote.session.Name()
}

// Close implements IAgent
func (remote *HostAgent) Close() {
	remote.session.Close()
}

// Session implements IAgent
func (remote *HostAgent) Session() network.ISession {
	return remote.session
}

// Post implements IAgent
func (remote *HostAgent) Post(service IService, call *network.Call) error {
	return remote.Host.Post(remote.session, call)
}

// Wait implements IAgent
func (remote *HostAgent) Wait(service IService, call *network.Call, timeout time.Duration) (Future, error) {
	return remote.Host.Wait(remote.session, call, timeout)
}

// Write implements IAgent
func (remote *HostAgent) Write(msg *network.Message) error {
	return remote.session.Write(msg)
}

// SessionStatusChanged implements network.ISessionHandler
func (remote *HostAgent) SessionStatusChanged(status network.SessionStatus) {
	remote.Host.sessionStatusChanged(remote, status)
}

// Read implements network.ISessionHandler
func (remote *HostAgent) Read(session network.ISession, msg *network.Message) {
	switch msg.Type {
	case network.MessageTypeRegistry:
		// 处理来自邻居节点的服务状态变更通知
		go remote.handleServiceRegistry(msg.Data)
	case network.MessageTypeCall:
		go remote.handleCall(msg.Data)
	case network.MessageTypeReturn:
		go remote.handleReturn(msg.Data)
	}
}

// handleServiceRegistry 处理来自邻居节点的服务注册消息
func (remote *HostAgent) handleServiceRegistry(data []byte) {
	srd, err := network.UnmarshalServiceRegistryData(data)
	if err != nil {
		log.Warnf("unmarshal service registry data from %s err: %s", remote.session, err)
	}
	remote.Host.handleServiceRegistry(remote, srd)
}

// handleCall 处理来自对本地服务的调用
func (remote *HostAgent) handleCall(data []byte) {
	call, err := network.UnmarshalCall(data)
	if err != nil {
		log.Warnf("unmarshal call from %s err: %s", remote.session, err)
	}
	log.Infof("start handle rpc call userID: %d serviceID: %d methodID: %d from %s",
		call.ID, call.ServiceID, call.MethodID, remote.session)
	callReturn, err := remote.Host.handleCall(call)
	if err != nil {
		log.Warnf("handle rpc call userID: %d serviceID: %d methodID: %d from %s err: %s",
			call.ID, call.ServiceID, call.MethodID, remote.session, err)
	}
	if callReturn == nil {
		return
	}
	data = callReturn.Marshal()
	msg := &network.Message{
		Type: network.MessageTypeReturn,
		Data: data,
	}
	err = remote.session.Write(msg)
	if err != nil {
		log.Warnf("handle rpc call userID: %d serviceID: %d methodID: %d from %s err: %s",
			call.ID, call.ServiceID, call.MethodID, remote.session, err)
	}
	log.Infof("finish handle rpc call userID: %d serviceID: %d methodID: %d from %s",
		call.ID, call.ServiceID, call.MethodID, remote.session)
}

// handleReturn 处理对远程服务的调用返回
func (remote *HostAgent) handleReturn(data []byte) {
	callReturn, err := network.UnmarshalReturn(data)
	if err != nil {
		log.Warn("%s read return err: %s", remote.session, err)
		return
	}
	remote.Host.Notify(callReturn)
}
