// -------------------------------------------
// @file      : game.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 上午11:13
// -------------------------------------------

package game

import (
	"go.uber.org/zap"
	"gogs/base/config"
	"gogs/base/etcd"
	"gogs/base/gscluster"
	"gogs/base/gserrors"
	"gogs/base/gsnet"
	log "gogs/base/logger"
	"gogs/game/model"
	"runtime"
)

var server *gscluster.Game

func Main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// 日志配置
	logConfig := config.GetLogConfig()
	// 游戏服配置
	gameConfig := config.GetGameConfig()
	if logConfig == nil {
		gserrors.Panic("unable to find log config")
	}
	if gameConfig == nil {
		gserrors.Panic("unable to find game config")
	}
	// 初始化日志
	log.Init(
		log.SetFilename(gameConfig.LogPath),
		log.SetIsOpenFile(logConfig.IsOpenFile),
		log.SetIsOpenErrorFile(logConfig.IsOpenErrorFile),
		log.SetIsOpenConsole(logConfig.IsOpenConsole),
		log.SetIsAsync(logConfig.IsAsync),
		log.SetMaxFileSize(int(logConfig.Maxsize)),
		log.SetStacktrace(zap.PanicLevel),
	)
	etcdConfig := config.GetEtcdConfig()
	if etcdConfig == nil {
		gserrors.Panic("unable to find etcd config")
	}
	if err := etcd.Init(etcdConfig, nil); err != nil {
		gserrors.Panicf("etcd init err:%s", err)
	}
	defer func() {
		// 等待异步日志写入完成
		_ = log.Close()
	}()
	model.InitMongoDB(config.ServerID)
	RegisterBuilders()
	var err error
	server, err = gscluster.NewGame(
		config.ServerID,
		config.ServerType,
		builders,
		"localhost:9102",
	)
	if err != nil {
		gserrors.Panicf("new game err:%s", err)
	}
	etcd.SetServiceCallback(EtcdNodeEventListener)
	// 处理系统信号
	ProcessSignal()
	<-exitChan
	server.Shutdown()
	model.CloseMongoDB()
}

// EtcdNodeEventListener 注册到etcd组件的节点状态变更事件处理器
func EtcdNodeEventListener(nodeEvent *etcd.NodeEvent) {
	switch nodeEvent.Event {
	case etcd.EventAdd:
		if nodeEvent.Node.GetType() == etcd.ServerTypeGate {
			if _, ok := server.Host.Node.GetSession(gsnet.DriverTypeCluster, nodeEvent.Node.GetConnectURL()); !ok {
				_, err := server.Host.Connect(nodeEvent.Node.GetConnectURL())
				if err != nil {
					log.Errorf("connect to gate err:%s", err)
				}
			}
		}
	case etcd.EventDelete:
		if nodeEvent.Node.GetType() == etcd.ServerTypeGate {
			session, ok := server.Host.Node.GetSession(gsnet.DriverTypeCluster, nodeEvent.Node.GetConnectURL())
			if ok {
				session.Close()
			}
		}
	default:
		log.Errorf("unknown etcd event:%+v", nodeEvent)
	}

}
