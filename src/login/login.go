// -------------------------------------------
// @file      : login.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 上午11:13
// -------------------------------------------

package login

import (
	"fmt"
	"go.uber.org/zap"
	"gogs/base/cberrors"
	"gogs/base/cluster"
	"gogs/base/cluster/network"
	"gogs/base/config"
	"gogs/base/etcd"
	log "gogs/base/logger"
	"gogs/cb"
	"gogs/game/model"
	"runtime"
)

var server *cluster.Normal

func Main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// 日志配置
	logConfig := config.GetLogConfig()
	// 登录服配置
	loginConfig := config.GetLoginConfig()
	if logConfig == nil {
		cberrors.Panic("unable to find log config")
	}
	if loginConfig == nil {
		cberrors.Panic("unable to find login config")
	}
	// 初始化日志
	log.Init(
		log.SetFilename(loginConfig.LogPath),
		log.SetIsOpenFile(logConfig.IsOpenFile),
		log.SetIsOpenErrorFile(logConfig.IsOpenErrorFile),
		log.SetIsOpenConsole(logConfig.IsOpenConsole),
		log.SetIsAsync(logConfig.IsAsync),
		log.SetMaxFileSize(int(logConfig.Maxsize)),
		log.SetStacktrace(zap.PanicLevel),
	)
	etcdConfig := config.GetEtcdConfig()
	if etcdConfig == nil {
		cberrors.Panic("unable to find etcd config")
	}
	if err := etcd.Init(etcdConfig, nil); err != nil {
		cberrors.Panic("etcd init err:%s", err)
	}
	defer func() {
		// 等待异步日志写入完成
		_ = log.Close()
	}()
	model.InitMongoDB(config.ServerID)
	RegisterBuilders()
	name := fmt.Sprintf("%s:%d", config.ServerType, config.ServerID)
	server = cluster.NewNormal(name, builders, "")
	server.Host.NewService(cb.LoginTypeName, name, nil)
	etcd.SetServiceCallback(EtcdNodeEventListener)
	CheckGateConn()
	// 处理系统信号
	ProcessSignal()
	<-exitChan
	server.Shutdown()
	model.CloseMongoDB()
}

// CheckGateConn 检查网关连接
func CheckGateConn() {
	nodes, err := etcd.GetDepListByType(etcd.ServerTypeGate)
	if err != nil {
		log.Errorf("get etcd gate list err:%s", err)
		return
	}
	for _, node := range nodes {
		if _, ok := server.Host.Node.GetSession(network.DriverTypeHost, node.GetConnectURL()); !ok {
			_, err := server.Host.Connect(node.GetConnectURL())
			if err != nil {
				log.Errorf("connect to gate err:%s", err)
			}
		}
	}
}

// EtcdNodeEventListener 注册到etcd组件的节点状态变更事件处理器
func EtcdNodeEventListener(nodeEvent *etcd.NodeEvent) {
	switch nodeEvent.Event {
	case etcd.EventAdd:
		if nodeEvent.Node.GetType() == etcd.ServerTypeGate {
			if _, ok := server.Host.Node.GetSession(network.DriverTypeHost, nodeEvent.Node.GetConnectURL()); !ok {
				_, err := server.Host.Connect(nodeEvent.Node.GetConnectURL())
				if err != nil {
					log.Errorf("connect to gate err:%s", err)
				}
			}
		}
	case etcd.EventDelete:
		if nodeEvent.Node.GetType() == etcd.ServerTypeGate {
			session, ok := server.Host.Node.GetSession(network.DriverTypeHost, nodeEvent.Node.GetConnectURL())
			if ok {
				session.Close()
			}
		}
	default:
		log.Errorf("unknown etcd event:%+v", nodeEvent)
	}

}
