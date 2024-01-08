// -------------------------------------------
// @file      : user.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/5 上午11:01
// -------------------------------------------

package login

import (
	"gogs/cb"
)

type API struct {
}

func NewUser() (*API, error) {
	return &API{}, nil
}

func (user *API) GetServerTime() (int64, cb.Code, error) {
	return 0, cb.CodeOK, nil
}

func (user *API) GetUserInfo() (*cb.UserInfo, cb.Code, error) {
	return nil, cb.CodeOK, nil
}
