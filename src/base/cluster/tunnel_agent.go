// -------------------------------------------
// @file      : tunnel_agent.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午6:44
// -------------------------------------------

package cluster

import (
	"gogs/base/cluster/network"
	"time"
)

// TunnelAgent game和GateServerRemoteService的中间层
// 用于game调用client的接口
// 通过GateServer的Tunnel接口将消息发送给client
type TunnelAgent struct {
	game       *Game
	userID     int64
	gateServer IGateServer
}

func newTunnelAgent(game *Game, userID int64, gateServer IGateServer) *TunnelAgent {
	return &TunnelAgent{
		game:       game,
		userID:     userID,
		gateServer: gateServer,
	}
}

// Session implement IAgent
func (agent *TunnelAgent) Session() network.ISession {
	return nil
}

// Post implement IAgent
func (agent *TunnelAgent) Post(service IService, call *network.Call) error {
	return agent.game.Post(agent, call)
}

// Wait implement IAgent
func (agent *TunnelAgent) Wait(service IService, call *network.Call, timeout time.Duration) (Future, error) {
	return agent.game.Wait(agent, call, timeout)
}

// Write implement IAgent
func (agent *TunnelAgent) Write(msg *network.Message) error {
	tunnelMsg := &TunnelMsg{
		UserID: agent.userID,
		Type:   msg.Type,
		Data:   msg.Data,
	}
	return agent.gateServer.Tunnel(tunnelMsg)
}

// DriverType implement network.ISession
func (agent *TunnelAgent) DriverType() network.DriverType {
	return network.DriverTypeHost
}

// Status implement network.ISession
func (agent *TunnelAgent) Status() network.SessionStatus {
	return network.SessionStatusInConnected
}

// Handler implement network.ISession
func (agent *TunnelAgent) Handler() network.ISessionHandler {
	return nil
}

// Close implement network.ISession
func (agent *TunnelAgent) Close() {
}

// Name implement network.ISession
func (agent *TunnelAgent) Name() string {
	return "tunnel remote"
}

// RemoteAddr implement network.ISession
func (agent *TunnelAgent) RemoteAddr() string {
	return ""
}
