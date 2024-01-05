// -------------------------------------------
// @file      : user.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/5 上午11:01
// -------------------------------------------

package game

import (
	"gogs/base/gscluster"
	"gogs/game/model"
	"gogs/idl"
)

type User struct {
	*model.DBUser
	agent *gscluster.ClientAgent
}

func NewUser(clusterUser *gscluster.ClientAgent) (*User, error) {
	return &User{
		agent: clusterUser,
	}, nil
}

func (user *User) GetServerTime() (int64, idl.Err, error) {
	return 0, idl.ErrOK, nil
}

func (user *User) GetUserInfo() (*idl.UserInfo, idl.Err, error) {
	return nil, idl.ErrOK, nil
}
