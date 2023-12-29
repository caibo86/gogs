// -------------------------------------------
// @file      : rpc_config.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/29 下午10:46
// -------------------------------------------

package config

type RPCConfig struct {
	Timeout int64 `yaml:"timeout"` // 超时时间,单位毫秒
}

// NewRPCConfig 创建RPC配置
func NewRPCConfig() *RPCConfig {
	c := &RPCConfig{
		Timeout: 5000,
	}
	return c
}

// GetRPCConfig 获取rpc配置
func GetRPCConfig() *RPCConfig {
	c := manager.configMap[KeyRPC]
	if c == nil {
		return nil
	}
	return c.(*RPCConfig)
}

// GetType implements IConfig
func (c *RPCConfig) GetType() string {
	return KeyRPC
}
