// -------------------------------------------
// @file      : user.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午6:27
// -------------------------------------------

package gscluster

import (
	"fmt"
	"gogs/base/mongodb"
)

type User struct {
	client    IRemoteService                          // 用户的客户端代理
	sessionID int64                                   // 会话ID
	id        int64                                   // 玩家ID
	actor     IActor                                  // 用户角色代理
	DBSave    func(client *mongodb.MongoClient) error // 存盘
}

// NewUser 新建用户
func NewUser(sessionID, userID int64) *User {
	return &User{
		sessionID: sessionID,
		id:        userID,
	}
}

// Client 获取用户的客户端代理
func (user *User) Client() (IRemoteService, bool) {
	client := user.client
	if client == nil {
		return nil, false
	}
	return client, true
}

// SetClient 设置用户的客户端代理
func (user *User) SetClient(client IRemoteService) {
	user.client = client
}

// SessionID 获取会话ID
func (user *User) SessionID() int64 {
	return user.sessionID
}

// ID 获取玩家ID
func (user *User) ID() int64 {
	return user.id
}

// String implement fmt.Stringer
func (user *User) String() string {
	return fmt.Sprintf("User(%d)(%d)", user.id, user.sessionID)
}

// Actor 获取用户角色代理
func (user *User) Actor() IActor {
	return user.actor
}

// SetActor 设置用户角色代理
func (user *User) SetActor(actor IActor) {
	user.actor = actor
}

// Save 存盘
func (user *User) Save(db *mongodb.MongoClient) error {
	if user.DBSave != nil {
		return user.DBSave(db)
	}
	return nil
}
