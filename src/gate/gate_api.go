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
	if api.GameServer() != nil {
		return nil, cb.CodeDuplicateLogin, nil
	}
	if api.authed {
		return nil, cb.CodeDuplicateLogin, nil
	}
	// TODO Login服务器处理登录
	_, err := api.GateAgent.Gate.Login(api.GateAgent, req, clientInfo)
	// TODO
}
