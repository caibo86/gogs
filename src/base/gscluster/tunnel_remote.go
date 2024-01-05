// -------------------------------------------
// @file      : tunnel_remote.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午6:44
// -------------------------------------------

package gscluster

import (
	"gogs/base/gsnet"
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

// Session implement IRemote
func (remote *TunnelRemote) Session() gsnet.ISession {
	return nil
}

// Post implement IRemote
func (remote *TunnelRemote) Post(service IService, call *gsnet.Call) error {
	return remote.game.Post(remote, call)
}

// Wait implement IRemote
func (remote *TunnelRemote) Wait(service IService, call *gsnet.Call, timeout time.Duration) (Future, error) {
	return remote.game.Wait(remote, call, timeout)
}

// Write implement IRemote
func (remote *TunnelRemote) Write(msg *gsnet.Message) error {
	tunnelMsg := &TunnelMsg{
		UserID: remote.userID,
		Type:   msg.Type,
		Data:   msg.Data,
	}
	return remote.gateServer.Tunnel(tunnelMsg)
}

// DriverType implement gsnet.ISession
func (remote *TunnelRemote) DriverType() gsnet.DriverType {
	return gsnet.DriverTypeCluster
}

// Status implement gsnet.ISession
func (remote *TunnelRemote) Status() gsnet.SessionStatus {
	return gsnet.SessionStatusInConnected
}

// Handler implement gsnet.ISession
func (remote *TunnelRemote) Handler() gsnet.ISessionHandler {
	return nil
}

// Close implement gsnet.ISession
func (remote *TunnelRemote) Close() {
}

// Name implement gsnet.ISession
func (remote *TunnelRemote) Name() string {
	return "tunnel remote"
}

// RemoteAddr implement gsnet.ISession
func (remote *TunnelRemote) RemoteAddr() string {
	return ""
}
