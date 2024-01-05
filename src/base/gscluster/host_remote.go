// -------------------------------------------
// @file      : gscluster.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 上午12:40
// -------------------------------------------

package gscluster

import (
	"gogs/base/gsnet"
	log "gogs/base/logger"
	"time"
)

// ClusterRemote 集群远程服务
// 实现了gsnet.ISessionHandler接口和IRemote接口
// 将Host与Session连接起来 在两者中间实现消息传递
// 负责识别Message结构 根据不同Code调用上层函数
type ClusterRemote struct {
	session gsnet.ISession // 远程服务 传输层通道
	Host    *Host          // 所属Host
}

// NewClusterRemote 新建集群远程服务
func NewClusterRemote(host *Host, session gsnet.ISession) *ClusterRemote {
	return &ClusterRemote{
		session: session,
		Host:    host,
	}
}

// Name 其实就是session的名字
func (remote *ClusterRemote) Name() string {
	return remote.session.Name()
}

// Close implements IRemote
func (remote *ClusterRemote) Close() {
	remote.session.Close()
}

// Session implements IRemote
func (remote *ClusterRemote) Session() gsnet.ISession {
	return remote.session
}

// Post implements IRemote
func (remote *ClusterRemote) Post(service IService, call *gsnet.Call) error {
	return remote.Host.Post(remote.session, call)
}

// Wait implements IRemote
func (remote *ClusterRemote) Wait(service IService, call *gsnet.Call, timeout time.Duration) (Future, error) {
	return remote.Host.Wait(remote.session, call, timeout)
}

// Write implements IRemote
func (remote *ClusterRemote) Write(msg *gsnet.Message) error {
	return remote.session.Write(msg)
}

// SessionStatusChanged implements gsnet.ISessionHandler
func (remote *ClusterRemote) SessionStatusChanged(status gsnet.SessionStatus) {
	remote.Host.sessionStatusChanged(remote, status)
}

// Read implements gsnet.ISessionHandler
func (remote *ClusterRemote) Read(session gsnet.ISession, msg *gsnet.Message) {
	switch msg.Type {
	case gsnet.MessageTypeRegistry:
		// 处理来自邻居节点的服务状态变更通知
		go remote.handleServiceRegistry(msg.Data)
	case gsnet.MessageTypeCall:
		go remote.handleCall(msg.Data)
	case gsnet.MessageTypeReturn:
		go remote.handleReturn(msg.Data)
	}
}

// handleServiceRegistry 处理来自邻居节点的服务注册消息
func (remote *ClusterRemote) handleServiceRegistry(data []byte) {
	srd, err := gsnet.UnmarshalServiceRegistryData(data)
	if err != nil {
		log.Warnf("unmarshal service registry data from %s err: %s", remote.session, err)
	}
	remote.Host.handleServiceRegistry(remote, srd)
}

// handleCall 处理来自对本地服务的调用
func (remote *ClusterRemote) handleCall(data []byte) {
	call, err := gsnet.UnmarshalCall(data)
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
	msg := &gsnet.Message{
		Type: gsnet.MessageTypeReturn,
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
func (remote *ClusterRemote) handleReturn(data []byte) {
	callReturn, err := gsnet.UnmarshalReturn(data)
	if err != nil {
		log.Warn("%s read return err: %s", remote.session, err)
		return
	}
	remote.Host.Notify(callReturn)
}
