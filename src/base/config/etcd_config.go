// -------------------------------------------
// @file      : etcd_config.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/29 下午10:46
// -------------------------------------------

package config

// EtcdConfig etcd服务配置
type EtcdConfig struct {
	Root        string   `yaml:"root"`     // 前缀
	Addr        []string `yaml:"addr"`     // etcd地址列表
	User        string   `yaml:"user"`     // etcd用户名
	Password    string   `yaml:"password"` // etcd密码
	DepList     []string `yaml:"depList"`  // 依赖的服务列表
	TTL         int64    `yaml:"ttl"`      // 租约时长second
	ServiceType string   // 本app服务类型
	ServiceID   int64    // 本app实例id
	ServiceAddr string   // 本app实例地址
	ServicePort string   // 本app实例端口
}

// NewEtcdConfig 创建etcd配置
func NewEtcdConfig() *EtcdConfig {
	c := &EtcdConfig{}
	return c
}

// GetEtcdConfig 获取etcd配置
func GetEtcdConfig() *EtcdConfig {
	c := manager.configMap[KeyEtcd]
	if c == nil {
		return nil
	}
	return c.(*EtcdConfig)
}

// GetType implement IConfig
func (c *EtcdConfig) GetType() string {
	return KeyEtcd
}
