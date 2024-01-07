// -------------------------------------------
// @file      : gate.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午9:35
// -------------------------------------------

package gate

import (
	"fmt"
	"gogs/base/cberrors"
	"gogs/base/cluster"
	"gogs/base/cluster/network"
	"gogs/base/config"
	"gogs/base/etcd"
	log "gogs/base/logger"
	"gogs/idl"
	"runtime"
)

var (
	// exitChan 退出信号
	exitChan = make(chan struct{}, 1)
)

func Main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// 日志配置
	logConfig := config.GetLogConfig()
	// 网关配置
	gateConfig := config.GetGateConfig()
	if logConfig == nil {
		cberrors.Panic("unable to find log config")
	}
	if gateConfig == nil {
		cberrors.Panic("unable to find gate config")
	}
	// 初始化日志
	log.Init(
		log.SetFilename(gateConfig.LogPath),
		log.SetIsOpenFile(logConfig.IsOpenFile),
		log.SetIsOpenErrorFile(logConfig.IsOpenErrorFile),
		log.SetIsOpenConsole(logConfig.IsOpenConsole),
		log.SetIsAsync(logConfig.IsAsync),
		log.SetMaxFileSize(int(logConfig.Maxsize)),
	)
	defer func() {
		// 等待异步日志写入完成
		_ = log.Close()
	}()
	// 处理系统信号
	ProcessSignal()
	// 外部地址
	addr := gateConfig.FullAddr()
	// 内部地址
	hostAddr := gateConfig.FullInnerAddr()
	// 网关名字
	name := fmt.Sprintf("%s:%d", config.ServerType, config.ServerID)
	// 网关服务构造器
	builder := idl.NewGateBuilder(func(service cluster.IService) (idl.IGate, error) {
		return NewRealGate(service.Context().(*cluster.GateRemote)), nil
	})
	log.Infof("gate: %s addr: %s inner addr: %s", name, addr, hostAddr)
	server, err := cluster.NewGate(name, addr, hostAddr, builder, network.ProtocolTCP)
	if err != nil {
		cberrors.Panic(err.Error())
	}
	// 启动监听后再启动etcd组件
	etcdConfig := config.GetEtcdConfig()
	if etcdConfig == nil {
		cberrors.Panic("unable to find etcd config")
	}
	config.Adjust(
		config.SetEtcdServiceType(config.ServerType),
		config.SetEtcdServiceID(config.ServerID),
		config.SetEtcdServiceAddr(gateConfig.InnerAddr),
		config.SetEtcdServicePort(gateConfig.InnerPort),
	)
	if err := etcd.Init(etcdConfig, nil); err != nil {
		cberrors.Panic("etcd init err:%s", err)
	}
	<-exitChan
	server.Close()
	etcd.Exit()
}
