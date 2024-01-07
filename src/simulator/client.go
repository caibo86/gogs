// -------------------------------------------
// @file      : client.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/7 下午6:21
// -------------------------------------------

package simulator

import "gogs/cb"

type ClientAPI struct {
}

func NewClientAPI() *ClientAPI {
	return &ClientAPI{}
}

func (c *ClientAPI) GetClientInfo() (ret0 *cb.ClientInfo, err error) {
	return
}
