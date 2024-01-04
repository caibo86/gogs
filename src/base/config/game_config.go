// -------------------------------------------
// @file      : game_config.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午9:43
// -------------------------------------------

package config

// GameConfig 游戏服配置
type GameConfig struct {
	LogPath string `yaml:"logPath"`
	DBName  string `yaml:"dbName"`
}

// NewGameConfig 创建游戏服配置
func NewGameConfig() *GameConfig {
	c := &GameConfig{}
	return c
}

// GetGameConfig 获取游戏服配置
func GetGameConfig() *GameConfig {
	c := manager.configMap[KeyGame]
	if c == nil {
		return nil
	}
	return c.(*GameConfig)
}

// GetType implements IConfig
func (c *GameConfig) GetType() string {
	return KeyGame
}
