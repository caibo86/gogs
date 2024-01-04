// -------------------------------------------
// @file      : login_config.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午9:45
// -------------------------------------------

package config

// LoginConfig 登录服务器配置
type LoginConfig struct {
}

// NewLoginConfig 创建登录服务器配置
func NewLoginConfig() *LoginConfig {
	c := &LoginConfig{}
	return c
}

// GetLoginConfig 获取登录服务器配置
func GetLoginConfig() *LoginConfig {
	c := manager.configMap[KeyLogin]
	if c == nil {
		return nil
	}
	return c.(*LoginConfig)
}

// GetType implements IConfig
func (c *LoginConfig) GetType() string {
	return KeyLogin
}
