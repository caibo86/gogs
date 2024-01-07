// -------------------------------------------
// @file      : tunnel_remote.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午6:44
// -------------------------------------------

package cluster

import (
	"gogs/base/cluster/network"
	"time"
)

// TunnelRemote game和GateServerRemoteService的中间层
// 用于game调用client的接口
// 通过GateServer的Tunnel接口将消息发送给client
type TunnelRemote struct {
	game       *Game
	userID     int64
	gateServer IGateServer
}

func newTunnelRemote(game *Game, userID int64, gateServer IGateServer) *TunnelRemote {
	return &TunnelRemote{
		game:       game,
		userID:     userID,
		gateServer: gateServer,
	}
}

// Session implement IAgent
func (remote *TunnelRemote) Session() network.ISession {
	return nil
}

// Post implement IAgent
func (remote *TunnelRemote) Post(service IService, call *network.Call) error {
	return remote.game.Post(remote, call)
}

// Wait implement IAgent
func (remote *TunnelRemote) Wait(service IService, call *network.Call, timeout time.Duration) (Future, error) {
	return remote.game.Wait(remote, call, timeout)
}

// Write implement IAgent
func (remote *TunnelRemote) Write(msg *network.Message) error {
	tunnelMsg := &TunnelMsg{
		UserID: remote.userID,
		Type:   msg.Type,
		Data:   msg.Data,
	}
	return remote.gateServer.Tunnel(tunnelMsg)
}

// DriverType implement network.ISession
func (remote *TunnelRemote) DriverType() network.DriverType {
	return network.DriverTypeHost
}

// Status implement network.ISession
func (remote *TunnelRemote) Status() network.SessionStatus {
	return network.SessionStatusInConnected
}

// Handler implement network.ISession
func (remote *TunnelRemote) Handler() network.ISessionHandler {
	return nil
}

// Close implement network.ISession
func (remote *TunnelRemote) Close() {
}

// Name implement network.ISession
func (remote *TunnelRemote) Name() string {
	return "tunnel remote"
}

// RemoteAddr implement network.ISession
func (remote *TunnelRemote) RemoteAddr() string {
	return ""
}
