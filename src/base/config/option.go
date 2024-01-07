// -------------------------------------------
// @file      : option.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/29 下午10:16
// -------------------------------------------

package config

// yaml配置文件中的key
const (
	KeyEtcd      = "etcd"      // etcd配置
	KeyLog       = "log"       // 日志配置
	KeyRPC       = "rpc"       // rpc host配置
	KeyGate      = "gate"      // 网关服务器配置
	KeyGame      = "game"      // 游戏服务器配置
	KeyLogin     = "login"     // 登录服务器配置
	KeyMap       = "map"       // 地图服务器配置
	KeyChat      = "chat"      // 聊天服务器配置
	KeyActivity  = "activity"  // 活动服务器配置
	KeyAlliance  = "alliance"  // 联盟服务器配置
	KeyFriend    = "friend"    // 好友服务器配置
	KeyGM        = "gm"        // GM服务器配置
	KeyMail      = "mail"      // 邮件服务器配置
	KeySimulator = "simulator" // 客户端模拟器配置
)

type Option func(*Manager)

// SetEtcdServiceType 把自身注册到etcd的服务类型
func SetEtcdServiceType(serviceType string) Option {
	return func(m *Manager) {
		if m == nil {
			return
		}
		config := m.GetConfigByType(KeyEtcd)
		if config != nil {
			config.(*EtcdConfig).ServiceType = serviceType
		}
	}
}

// SetEtcdServiceID 自身注册到etcd的服务实例id
func SetEtcdServiceID(serviceID int64) Option {
	return func(m *Manager) {
		if m == nil {
			return
		}
		config := m.GetConfigByType(KeyEtcd)
		if config != nil {
			config.(*EtcdConfig).ServiceID = serviceID
		}
	}
}

// SetEtcdServiceAddr 自身注册到etcd的服务地址
func SetEtcdServiceAddr(addr string) Option {
	return func(m *Manager) {
		if m == nil {
			return
		}
		config := m.GetConfigByType(KeyEtcd)
		if config != nil {
			config.(*EtcdConfig).ServiceAddr = addr
		}
	}
}

// SetEtcdServicePort 自身注册到etcd的服务端口
func SetEtcdServicePort(port string) Option {
	return func(m *Manager) {
		if m == nil {
			return
		}
		config := m.GetConfigByType(KeyEtcd)
		if config != nil {
			config.(*EtcdConfig).ServicePort = port
		}
	}
}
