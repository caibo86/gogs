// -------------------------------------------
// @file      : log_config.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/29 下午10:46
// -------------------------------------------

package config

// LogConfig 日志配置
type LogConfig struct {
	Level           int32 `yaml:"level"`
	Maxsize         int32 `yaml:"maxSize"`
	Flag            uint8 `yaml:"flag"`
	Profile         bool  `yaml:"profile"`
	IsOpenConsole   bool  `yaml:"isOpenConsole"`
	IsOpenFile      bool  `yaml:"isOpenFile"`
	IsOpenErrorFile bool  `yaml:"isOpenErrorFile"`
	IsAsync         bool  `yaml:"isAsync"`
}

// NewLogConfig 创建日志配置
func NewLogConfig() *LogConfig {
	c := &LogConfig{
		Level:           0,
		Maxsize:         128,
		IsAsync:         false,
		IsOpenFile:      true,
		IsOpenErrorFile: true,
		IsOpenConsole:   true,
	}
	return c
}

// GetLogConfig 获取日志配置
func GetLogConfig() *LogConfig {
	c := manager.configMap[KeyLog]
	if c == nil {
		return nil
	}
	return c.(*LogConfig)
}

// GetType implements IConfig
func (c *LogConfig) GetType() string {
	return KeyLog
}
