// -------------------------------------------
// @file      : user.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/5 上午12:06
// -------------------------------------------

package model

import "gogs/base/gscluster"

type User struct {
	*DBUser
	clusterUser *gscluster.User
}

func NewUser(clusterUser *gscluster.User) (*User, error) {
	return &User{
		clusterUser: clusterUser,
	}, nil
}

type DBUser struct {
}
