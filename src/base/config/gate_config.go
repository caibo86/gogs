// -------------------------------------------
// @file      : gate_config.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午9:38
// -------------------------------------------

package config

import "fmt"

// GateConfig 网关配置
type GateConfig struct {
	InnerAddr string `yaml:"innerAddr"`
	InnerPort string `yaml:"innerPort"`
	Addr      string `yaml:"addr"`
	Port      string `yaml:"port"`
	LogPath   string `yaml:"logPath"`
}

// NewGateConfig 创建网关配置
func NewGateConfig() *GateConfig {
	c := &GateConfig{}
	return c
}

// GetGateConfig 获取网关配置
func GetGateConfig() *GateConfig {
	c := manager.configMap[KeyGate]
	if c == nil {
		return nil
	}
	return c.(*GateConfig)
}

// GetType implements IConfig
func (c *GateConfig) GetType() string {
	return KeyGate
}

func (c *GateConfig) FullAddr() string {
	return fmt.Sprintf("%s:%s", c.Addr, c.Port)
}

func (c *GateConfig) FullInnerAddr() string {
	return fmt.Sprintf("%s:%s", c.InnerAddr, c.InnerPort)
}
