// -------------------------------------------
// @file      : client_agent.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午6:27
// -------------------------------------------

package gscluster

import (
	"gogs/base/mongodb"
)

// ClientAgent 客户端代理,最终挂载到Game上运行的User对象上,用于发送请求给客户端
type ClientAgent struct {
	clientService IRemoteService                          // 挂载客户端API的远程服务
	sessionID     int64                                   // 会话ID
	userID        int64                                   // 玩家ID
	actor         IActor                                  // 用户角色
	DBSave        func(client *mongodb.MongoClient) error // 存盘
}

// NewClientAgent 新建用户
func NewClientAgent(sessionID, userID int64) *ClientAgent {
	return &ClientAgent{
		sessionID: sessionID,
		userID:    userID,
	}
}

// ClientService 获取用户的客户端代理
func (agent *ClientAgent) ClientService() (IRemoteService, bool) {
	client := agent.clientService
	if client == nil {
		return nil, false
	}
	return client, true
}

// SetClientService 设置用户的客户端代理
func (agent *ClientAgent) SetClientService(client IRemoteService) {
	agent.clientService = client
}

// SessionID 获取会话ID
func (agent *ClientAgent) SessionID() int64 {
	return agent.sessionID
}

// UserID 获取玩家ID
func (agent *ClientAgent) UserID() int64 {
	return agent.userID
}

// Actor 获取用户角色代理
func (agent *ClientAgent) Actor() IActor {
	return agent.actor
}

// SetActor 设置用户角色代理
func (agent *ClientAgent) SetActor(actor IActor) {
	agent.actor = actor
}

// Save 存盘
func (agent *ClientAgent) Save(db *mongodb.MongoClient) error {
	if agent.DBSave != nil {
		return agent.DBSave(db)
	}
	return nil
}
