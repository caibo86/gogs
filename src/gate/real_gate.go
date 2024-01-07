// -------------------------------------------
// @file      : base_gate.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午10:20
// -------------------------------------------

package gate

import (
	"gogs/base/cluster"
	log "gogs/base/logger"
	"gogs/idl"
)

// RealGate 实际网关服务提供者
type RealGate struct {
	remote *cluster.GateRemote
}

// NewRealGate 新建网关服务提供者
func NewRealGate(remote *cluster.GateRemote) *RealGate {
	return &RealGate{
		remote: remote,
	}
}

// Login 登录
func (gate *RealGate) Login(req *idl.LoginReq, clientInfo *idl.ClientInfo) (*idl.LoginAck, idl.Err, error) {
	log.Debugf("login req:%+v, clientInfo:%+v", req, clientInfo)
	return &idl.LoginAck{
		AccountID: req.AccountID,
		UserID:    req.UserID,
		ServerID:  req.ServerID,
	}, idl.ErrOK, nil
}
