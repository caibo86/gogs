// -------------------------------------------
// @file      : simulator_config.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/7 下午6:15
// -------------------------------------------

package config

// SimulatorConfig 游戏服配置
type SimulatorConfig struct {
	LogPath string `yaml:"logPath"`
}

// NewSimulatorConfig 创建模拟器
func NewSimulatorConfig() *SimulatorConfig {
	c := &SimulatorConfig{}
	return c
}

// GetSimulatorConfig 获取游戏服配置
func GetSimulatorConfig() *SimulatorConfig {
	c := manager.configMap[KeySimulator]
	if c == nil {
		return nil
	}
	return c.(*SimulatorConfig)
}

// GetType implements IConfig
func (c *SimulatorConfig) GetType() string {
	return KeySimulator
}
