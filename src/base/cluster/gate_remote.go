// -------------------------------------------
// @file      : gate_remote.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 上午11:19
// -------------------------------------------

package cluster

import (
	"gogs/base/cberrors"
	"gogs/base/cluster/network"
	log "gogs/base/logger"
	"time"
)

const (
	gateID = ID(1)
	gameID = ID(2)
)

// GateRemote 网关远程代理
// Gate和GateSession的中间层,每一个GateSession都有一个GateRemote
// 实现了 IAgent 和 network.ISessionHandler 接口
type GateRemote struct {
	gate       *Gate    // 所属网关
	host       *Host    // 网关挂载的集群服务器
	service    IService // 承载的服务
	gameServer IGameServer
	session    network.ISession // 网关会话
	sessionID  int64            // 网关会话ID
	userID     int64            // 用户ID
}

// newGateRemote 新建网关远程代理
func newGateRemote(gate *Gate, session network.ISession, sessionID int64) (*GateRemote, error) {
	remote := &GateRemote{
		gate:      gate,
		host:      gate.host,
		session:   session,
		sessionID: sessionID,
	}
	// 网关本地服务的上下文存的这个GateRemote
	var err error
	remote.service, err = gate.builder.NewService(session.Name(), gateID, remote)
	return remote, err
}

// rProxy 网关代理
func (remote *GateRemote) rProxy(userID int64, service IService) (Err, error) {
	gameServer := service.(IGameServer)
	rProxyMsg := &RProxyMsg{
		UserID:    userID,
		SessionID: remote.sessionID,
		Gate:      remote.gate.name,
	}
	id, status, err := gameServer.Login(rProxyMsg)
	if err != nil || status != ErrOK {
		return status, cberrors.New("call GameServer#Login(%s) status: %s err: %s", gameServer, status, err)
	}
	remote.gameServer = gameServer
	remote.userID = id
	return ErrOK, nil
}

// GameServer .
func (remote *GateRemote) GameServer() IGameServer {
	return remote.gameServer
}

// UserID .
func (remote *GateRemote) UserID() int64 {
	return remote.userID
}

// Session .
func (remote *GateRemote) Session() network.ISession {
	return remote.session
}

// Post implements IAgent
func (remote *GateRemote) Post(service IService, call *network.Call) error {
	return remote.gate.Post(remote.session, call)
}

// Wait implements IAgent
func (remote *GateRemote) Wait(service IService, call *network.Call, timeout time.Duration) (Future, error) {
	return remote.gate.Wait(remote.session, call, timeout)
}

// Write implements IAgent
func (remote *GateRemote) Write(msg *network.Message) error {
	return remote.session.Write(msg)
}

// SessionStatusChanged implements network.ISessionHandler
func (remote *GateRemote) SessionStatusChanged(status network.SessionStatus) {
	if status == network.SessionStatusClosed && remote.gameServer != nil {
		remote.gate.sessionStatusChanged(remote, status)
		rProxyMsg := &RProxyMsg{
			UserID:    remote.userID,
			SessionID: remote.sessionID,
			Gate:      remote.gate.name,
		}
		_ = remote.gameServer.Logout(rProxyMsg)
	}
}

// Read implements network.ISessionHandler
func (remote *GateRemote) Read(session network.ISession, msg *network.Message) {
	switch msg.Type {
	case network.MessageTypeCall:
		go remote.handleCall(msg.Data)
	case network.MessageTypeReturn:
		go remote.handleReturn(msg.Data)
	}
}

// handleCall 处理对本地网关服务的调用
func (remote *GateRemote) handleCall(data []byte) {
	call, err := network.UnmarshalCall(data)
	if err != nil {
		log.Warnf("unmarshal call from %s err: %s", remote.session, err)
		return
	}
	switch ID(call.ServiceID) {
	case gateID:
		var callReturn *network.Return
		callReturn, err = remote.service.Call(call)
		if err != nil {
			log.Warnf("handle rpc call userID: %d serviceID: %d methodID: %d from %s err: %s",
				call.ID, call.ServiceID, call.MethodID, remote.session, err)
			return
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
			return
		}
	case gameID:
		if remote.gameServer != nil {
			tunnelMsg := &TunnelMsg{
				UserID: remote.userID,
				Type:   network.MessageTypeCall,
				Data:   data,
			}
			err = remote.gameServer.Tunnel(tunnelMsg)
			if err != nil {
				log.Warnf("rProxy call from %s to %s err: %s",
					remote.session, remote.gameServer, err)
				return
			}
		} else {
			log.Warnf("call %v from a non logged in user: %s", call, remote.session)
		}
	}
}

// handleReturn 处理对远程服务的调用返回
func (remote *GateRemote) handleReturn(data []byte) {
	gameServer := remote.gameServer
	if gameServer == nil {
		return
	}
	tunnelMsg := &TunnelMsg{
		UserID: remote.userID,
		Type:   network.MessageTypeReturn,
		Data:   data,
	}
	err := gameServer.Tunnel(tunnelMsg)
	if err != nil {
		log.Errorf("send tunnel msg to %s err: %s", gameServer, err)
	}
}
