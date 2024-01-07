// -------------------------------------------
// @file      : gate_agent.go
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

// GateAgent 网关远程代理
// Gate和GateSession的中间层,每一个GateSession都有一个GateRemote
// 实现了 IAgent 和 network.ISessionHandler 接口
type GateAgent struct {
	Gate       *Gate            // 所属网关
	host       *Host            // 网关挂载的集群服务器
	service    IService         // 承载的服务
	gameServer IGameServer      // 游戏服务器远程服务
	session    network.ISession // 网关会话
	sessionID  int64            // 网关会话ID
	userID     int64            // 用户id
}

// newGateAgent 新建网关远程代理
func newGateAgent(gate *Gate, session network.ISession, sessionID int64) (*GateAgent, error) {
	remote := &GateAgent{
		Gate:      gate,
		host:      gate.host,
		session:   session,
		sessionID: sessionID,
	}
	// 网关本地服务的上下文存的这个 GateAgent
	var err error
	remote.service, err = gate.builder.NewService(session.Name(), gateID, remote)
	return remote, err
}

// rProxy 网关代理
func (agent *GateAgent) rProxy(userID int64, service IService) (Err, error) {
	gameServer := service.(IGameServer)
	rProxyMsg := &UserLoginMsg{
		UserID:    userID,
		SessionID: agent.sessionID,
		Gate:      agent.Gate.name,
	}
	id, status, err := gameServer.Login(rProxyMsg)
	if err != nil || status != ErrOK {
		return status, cberrors.New("call GameServer#Login(%s) status: %s err: %s", gameServer, status, err)
	}
	agent.gameServer = gameServer
	agent.userID = id
	return ErrOK, nil
}

// GameServer .
func (agent *GateAgent) GameServer() IGameServer {
	return agent.gameServer
}

// ActorName 绑定的角色唯一Name
func (agent *GateAgent) ActorName() string {
	return agent.actorName
}

// Close implements IAgent
func (agent *GateAgent) Close() {
}

// Session implements IAgent
func (agent *GateAgent) Session() network.ISession {
	return agent.session
}

// Post implements IAgent
func (agent *GateAgent) Post(service IService, call *network.Call) error {
	return agent.Gate.Post(agent.session, call)
}

// Wait implements IAgent
func (agent *GateAgent) Wait(service IService, call *network.Call, timeout time.Duration) (Future, error) {
	return agent.Gate.Wait(agent.session, call, timeout)
}

// Write implements IAgent
func (agent *GateAgent) Write(msg *network.Message) error {
	return agent.session.Write(msg)
}

// SessionStatusChanged implements network.ISessionHandler
func (agent *GateAgent) SessionStatusChanged(status network.SessionStatus) {
	if status == network.SessionStatusClosed && agent.gameServer != nil {
		agent.Gate.sessionStatusChanged(agent, status)
		rProxyMsg := &RProxyMsg{
			UserID:    agent.userID,
			SessionID: agent.sessionID,
			Gate:      agent.Gate.name,
		}
		_ = agent.gameServer.Logout(rProxyMsg)
	}
}

// Read implements network.ISessionHandler
func (agent *GateAgent) Read(session network.ISession, msg *network.Message) {
	switch msg.Type {
	case network.MessageTypeCall:
		go agent.handleCall(msg.Data)
	case network.MessageTypeReturn:
		go agent.handleReturn(msg.Data)
	}
}

// handleCall 处理对本地网关服务的调用
func (agent *GateAgent) handleCall(data []byte) {
	call, err := network.UnmarshalCall(data)
	if err != nil {
		log.Warnf("unmarshal call from %s err: %s", agent.session, err)
		return
	}
	log.Infof("handle rpc call id: %d serviceID: %d methodID: %d from %s",
		call.ID, call.ServiceID, call.MethodID, agent.session)
	switch ID(call.ServiceID) {
	case gateID:
		var callReturn *network.Return
		callReturn, err = agent.service.Call(call)
		if err != nil {
			log.Warnf("handle rpc call id: %d serviceID: %d methodID: %d from %s err: %s",
				call.ID, call.ServiceID, call.MethodID, agent.session, err)
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
		err = agent.session.Write(msg)
		if err != nil {
			log.Warnf("handle rpc call userID: %d serviceID: %d methodID: %d from %s err: %s",
				call.ID, call.ServiceID, call.MethodID, agent.session, err)
			return
		}
	case gameID:
		if agent.gameServer != nil {
			tunnelMsg := &TunnelMsg{
				UserID: agent.userID,
				Type:   network.MessageTypeCall,
				Data:   data,
			}
			err = agent.gameServer.Tunnel(tunnelMsg)
			if err != nil {
				log.Warnf("rProxy call from %s to %s err: %s",
					agent.session, agent.gameServer, err)
				return
			}
		} else {
			log.Warnf("call %v from a non logged in user: %s", call, agent.session)
		}
	default:
		log.Errorf("call %v from %s, unknown serviceID: %d", call, agent.session, call.ServiceID)
	}
}

// handleReturn 处理对远程服务的调用返回
func (agent *GateAgent) handleReturn(data []byte) {
	gameServer := agent.gameServer
	if gameServer == nil {
		return
	}
	tunnelMsg := &TunnelMsg{
		UserID: agent.userID,
		Type:   network.MessageTypeReturn,
		Data:   data,
	}
	err := gameServer.Tunnel(tunnelMsg)
	if err != nil {
		log.Errorf("send tunnel msg to %s err: %s", gameServer, err)
	}
}
