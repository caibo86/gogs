// -------------------------------------------
// @file      : option.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/29 下午10:16
// -------------------------------------------

package config

// yaml配置文件中的key
const (
	KeyEtcd = "etcd" // etcd配置
	KeyLog  = "log"  // 日志配置
	KeyRPC  = "rpc"  // dock配置
)

type Option func(*Manager)
