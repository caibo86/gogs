// -------------------------------------------
// @file      : map_config.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午9:45
// -------------------------------------------

package config

// MapConfig 地图服务器配置
type MapConfig struct {
}

// NewMapConfig 创建地图服务器配置
func NewMapConfig() *MapConfig {
	c := &MapConfig{}
	return c
}

// GetMapConfig 获取地图服务器配置
func GetMapConfig() *MapConfig {
	c := manager.configMap[KeyMap]
	if c == nil {
		return nil
	}
	return c.(*MapConfig)
}

// GetType implements IConfig
func (c *MapConfig) GetType() string {
	return KeyMap
}
