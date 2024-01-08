// -------------------------------------------
// @file      : user.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/5 上午11:01
// -------------------------------------------

package game

import (
	"gogs/base/cluster"
	"gogs/cb"
	"gogs/game/model"
)

type User struct {
	*model.DBUser
	agent *cluster.ClientAgent
}

func NewUser(clusterUser *cluster.ClientAgent) (*User, error) {
	return &User{
		agent: clusterUser,
	}, nil
}

func (user *User) GetServerTime() (int64, cb.Code, error) {
	return 0, cb.CodeOK, nil
}

func (user *User) GetUserInfo() (*cb.UserInfo, cb.Code, error) {
	return nil, cb.CodeOK, nil
}
