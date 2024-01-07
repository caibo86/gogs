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
	HostSessionCache        int   `yaml:"hostSessionCache"`        // 集群会话消息发送缓存大小
	HostSessionHeartbeat    int   `yaml:"hostSessionHeartbeat"`    // 集群会话心跳间隔,单位秒
	ClientSessionCache      int   `yaml:"clientSessionCache"`      // 客户端会话消息发送缓存大小
	ClusterRegistryInterval int   `yaml:"clusterRegistryInterval"` // 集群服务注册时间间隔,单位秒
	ClusterRegistryMax      int   `yaml:"clusterRegistryMax"`      // 集群单次注册服务的最大数量
	ActorGroups             int   `yaml:"actorGroups"`             // 用户散列分组数量
}

// NewRPCConfig 创建RPC配置
func NewRPCConfig() *RPCConfig {
	c := &RPCConfig{
		Timeout:                 5,
		ListenRetryInterval:     5,
		GateSessionCache:        4096,
		HostSessionCache:        4096,
		HostSessionHeartbeat:    30,
		ClientSessionCache:      64,
		ClusterRegistryInterval: 2,
		ClusterRegistryMax:      128,
		ActorGroups:             128,
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

func HostSessionCache() int {
	return GetRPCConfig().HostSessionCache
}

func HostSessionHeartbeat() time.Duration {
	return time.Duration(GetRPCConfig().HostSessionHeartbeat) * time.Second
}

func ClientSessionCache() int {
	return GetRPCConfig().ClientSessionCache
}

func ClusterRegistryInterval() time.Duration {
	return time.Duration(GetRPCConfig().ClusterRegistryInterval) * time.Second
}

func ClusterRegistryMax() int {
	return GetRPCConfig().ClusterRegistryMax
}

func ActorGroups() int {
	return GetRPCConfig().ActorGroups
}
