// -------------------------------------------
// @file      : rpc_config.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/29 下午10:46
// -------------------------------------------

package config

import "time"

type RPCConfig struct {
	Timeout                 int64 `yaml:"timeout"`                 // 超时时间,单位秒
	ListenRetryInterval     int64 `yaml:"listenRetryInterval"`     // 监听重试时间,单位秒
	GateSessionCache        int   `yaml:"gateSessionCache"`        // 网关会话消息发送缓存大小
	ClusterSessionCache     int   `yaml:"clusterSessionCache"`     // 集群会话消息发送缓存大小
	ClusterSessionHeartbeat int   `yaml:"clusterSessionHeartbeat"` // 集群会话心跳间隔,单位秒
	ClientSessionCache      int   `yaml:"clientSessionCache"`      // 客户端会话消息发送缓存大小
}

// NewRPCConfig 创建RPC配置
func NewRPCConfig() *RPCConfig {
	c := &RPCConfig{
		Timeout:                 5,
		ListenRetryInterval:     5,
		GateSessionCache:        4096,
		ClusterSessionCache:     4096,
		ClusterSessionHeartbeat: 30,
		ClientSessionCache:      64,
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

func RPCTimeout() time.Duration {
	return time.Duration(GetRPCConfig().Timeout) * time.Second
}

func ListenRetryInterval() time.Duration {
	return time.Duration(GetRPCConfig().ListenRetryInterval) * time.Second
}

func GateSessionCache() int {
	return GetRPCConfig().GateSessionCache
}

func ClusterSessionCache() int {
	return GetRPCConfig().ClusterSessionCache
}

func ClusterSessionHeartbeat() time.Duration {
	return time.Duration(GetRPCConfig().ClusterSessionHeartbeat) * time.Second
}

func ClientSessionCache() int {
	return GetRPCConfig().ClientSessionCache
}
