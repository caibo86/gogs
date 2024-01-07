// -------------------------------------------
// @file      : config.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/29 下午10:04
// -------------------------------------------

package config

import (
	"encoding/json"
	"fmt"
	"gogs/base/cberrors"
	"gogs/base/misc"
	"gopkg.in/yaml.v2"
	"os"
	"path"
)

var (
	manager  = &Manager{} // 全局配置管理器
	basePath string       // 配置文件路径
)

func init() {
	basePath = misc.GetPathInRootDir("config")
}

// IConfig 配置接口
type IConfig interface {
	GetType() string
}

// Manager 全局配置管理器
type Manager struct {
	keys      []string
	configMap map[string]IConfig
}

// AddConfig 添加配置
func (m *Manager) AddConfig(c IConfig) {
	if m.configMap == nil {
		m.configMap = make(map[string]IConfig)
	}
	if _, ok := m.configMap[c.GetType()]; ok {
		panic(fmt.Errorf("dup subconfig type:%s", c.GetType()))
	}
	m.configMap[c.GetType()] = c
}

// GetConfigByType 根据类型获取配置
func (m *Manager) GetConfigByType(t string) IConfig {
	return m.configMap[t]
}

// String 实现fmt.Stringer接口
func (m *Manager) String() string {
	data, _ := json.Marshal(m.configMap)
	return string(data)
}

// CheckConfig 检查配置是否完整
func (m *Manager) CheckConfig() {
	for _, key := range m.keys {
		if _, ok := m.configMap[key]; !ok {
			cberrors.Panic("config %s not set", key)
		}
	}
}

// With 设置需要使用的配置
func With(keys ...string) {
	manager.keys = keys
	for _, key := range manager.keys {
		switch key {
		case KeyEtcd:
			manager.AddConfig(NewEtcdConfig())
		case KeyLog:
			manager.AddConfig(NewLogConfig())
		case KeyRPC:
			manager.AddConfig(NewRPCConfig())
		case KeyGate:
			manager.AddConfig(NewGateConfig())
		case KeyGame:
			manager.AddConfig(NewGameConfig())
		case KeyLogin:
			manager.AddConfig(NewLoginConfig())
		case KeyMap:
			manager.AddConfig(NewMapConfig())
		default:
			cberrors.Panic("unknown config type:%s", key)
		}
	}
}

// GlobalConfig 打印全局所有配置
func GlobalConfig() string {
	return manager.String()
}

// Adjust 根据选项调整配置
func Adjust(options ...Option) {
	for _, option := range options {
		option(manager)
	}
}

// LoadGlobalConfig 读取yaml配置文件
// 配置读取完成前不使用logger
func LoadGlobalConfig(cfgFilename string) {
	baseFilename := path.Join(basePath, "base.yml")
	baseData, err := os.ReadFile(baseFilename)
	if err != nil {
		cberrors.PanicfWith(err, "LoadGlobalConfig ReadFile(%s) err", baseFilename)
		return
	}
	filename := path.Join(basePath, cfgFilename)
	data, err := os.ReadFile(filename)
	if err != nil {
		cberrors.PanicfWith(err, "LoadGlobalConfig ReadFile(%s) err", filename)
		return
	}
	content := os.ExpandEnv(string(baseData) + string(data))
	if err = ParseGlobalConfig(content); err != nil {
		cberrors.PanicWith(err, "ParseGlobalConfig err")
	}
}

// ParseGlobalConfig 解析配置
func ParseGlobalConfig(content string) error {
	temp := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(content), &temp); err != nil {
		return err
	}
	for _, config := range manager.configMap {
		data, err := yaml.Marshal(temp[config.GetType()])
		if err != nil {
			return err
		}
		if err = yaml.Unmarshal(data, config); err != nil {
			return err
		}
	}
	manager.CheckConfig()
	return nil
}
