// -------------------------------------------
// @file      : cluster.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 上午12:40
// -------------------------------------------

package gsdock

import (
	"gogs/base/gsnet"
	log "gogs/base/logger"
	"time"
)

// Neighbor 集群中的邻居
type Neighbor struct {
	remote       *ClusterRemote
	services     map[string]IRemoteService // 按名字索引的邻居节点上的远程服务
	servicesByID map[ID]IRemoteService     // 按ID索引的邻居节点上的远程服务
}

// NewNeighbor 新建邻居节点
func NewNeighbor(remote *ClusterRemote) *Neighbor {
	return &Neighbor{
		remote:       remote,
		services:     make(map[string]IRemoteService),
		servicesByID: make(map[ID]IRemoteService),
	}
}

// ClusterRemote 集群远程服务
// 实现了gsnet.ISessionHandler接口和IRemote接口
// 将Dock与Session连接起来 在两者中间实现消息传递
// 负责识别Message结构 根据不同Code调用上层函数
type ClusterRemote struct {
	session gsnet.ISession // 远程服务 传输层通道
	Dock    *GSDock        // 所属Dock
}

// NewClusterRemote 新建集群远程服务
func NewClusterRemote(dock *GSDock, session gsnet.ISession) *ClusterRemote {
	return &ClusterRemote{
		session: session,
		Dock:    dock,
	}
}

// String implements fmt.Stringer
func (remote *ClusterRemote) String() string {
	return remote.session.String()
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
	return remote.Dock.Post(remote.session, call)
}

// Wait implements IRemote
func (remote *ClusterRemote) Wait(service IService, call *gsnet.Call, timeout time.Duration) (Future, error) {
	return remote.Dock.Wait(remote.session, call, timeout)
}

// Write implements IRemote
func (remote *ClusterRemote) Write(msg *gsnet.Message) error {
	return remote.session.Write(msg)
}

// StatusChanged implements gsnet.ISessionHandler
func (remote *ClusterRemote) StatusChanged(status gsnet.SessionStatus) {
	remote.Dock.statusChanged(remote, status)
}

// Read implements gsnet.ISessionHandler
func (remote *ClusterRemote) Read(session gsnet.ISession, msg *gsnet.Message) {
	switch msg.Type {
	case gsnet.MessageTypeCER:
		// 处理来自邻居节点的服务状态变更通知
		go remote.handleCERs(msg.Data)
	case gsnet.MessageTypeCall:
		go remote.handleCall(msg.Data)
	case gsnet.MessageTypeReturn:
		go remote.handleReturn(msg.Data)
	}
}

// handleCERs 处理来自邻居节点的服务状态变更通知
func (remote *ClusterRemote) handleCERs(data []byte) {
	cers, err := gsnet.UnmarshalCERs(data)
	if err != nil {
		log.Warn("%s read cers err: %s", remote.session, err)
	}
	remote.Dock.handleCERs(remote, cers)
}

// handleCall 处理来自对本地服务的调用
func (remote *ClusterRemote) handleCall(data []byte) {
	call, err := gsnet.UnmarshalCall(data)
	if err != nil {
		log.Warn("%s read call err: %s", remote.session, err)
	}
	callReturn, err := remote.Dock.handleCall(call)
	if err != nil {
		log.Warn("handle rpc call id: %d serviceID: %d methodID: %d from %s err: %s",
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
		log.Warn("handle rpc call id: %d serviceID: %d methodID: %d from %s err: %s",
			call.ID, call.ServiceID, call.MethodID, remote.session, err)
	}
}

// handleReturn 处理对远程服务的调用返回
func (remote *ClusterRemote) handleReturn(data []byte) {
	callReturn, err := gsnet.UnmarshalReturn(data)
	if err != nil {
		log.Warn("%s read return err: %s", remote.session, err)
		return
	}
	remote.Dock.Notify(callReturn)
}
