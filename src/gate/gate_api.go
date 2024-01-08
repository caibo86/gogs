// -------------------------------------------
// @file      : gate_api.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午10:20
// -------------------------------------------

package gate

import (
	"gogs/base/cluster"
	log "gogs/base/logger"
	"gogs/cb"
)

// API 网关对Client提供的API
type API struct {
	GateAgent *cluster.GateAgent
	authed    bool
}

// NewAPI 新建网关服务提供者
func NewAPI(agent *cluster.GateAgent) *API {
	return &API{
		GateAgent: agent,
	}
}

// Login 登录
func (api *API) Login(req *cb.LoginReq, clientInfo *cb.ClientInfo) (*cb.LoginAck, cb.Code, error) {
	log.Debugf("login req:%+v, clientInfo:%+v", req, clientInfo)
	if api.GateAgent.GameServer() != nil {
		return nil, cb.CodeDuplicateLogin, nil
	}
	if api.authed {
		return nil, cb.CodeDuplicateLogin, nil
	}
	// TODO Login服务器处理登录
	ntf := &cluster.UserLoginNtf{
		UserID:    req.UserID,
		ServerID:  req.ServerID,
		AccountID: req.AccountID,
	}
	ci := &cluster.ClientInfo{
		OpenUDID:             clientInfo.OpenUDID,
		Language:             clientInfo.Language,
		OS:                   clientInfo.OS,
		ClientVersion:        clientInfo.ClientVersion,
		ClientMemoryCapacity: clientInfo.ClientMemoryCapacity,
		ClientDeviceLevel:    clientInfo.ClientDeviceLevel,
		ClientChannel:        clientInfo.ClientChannel,
		IsAndroidEmulator:    clientInfo.IsAndroidEmulator,
	}
	internalErr, err := api.GateAgent.Gate.Login(api.GateAgent, ntf, ci)
	if err != nil || internalErr != cluster.ErrOK {
		log.Errorf("login error, internalErr: %s, err: %s", internalErr, err)
		return nil, cb.CodeSystemErr, nil
	}
	ack := &cb.LoginAck{
		UserID:    req.UserID,
		AccountID: req.AccountID,
		ServerID:  req.ServerID,
	}
	return ack, cb.CodeOK, nil
}
