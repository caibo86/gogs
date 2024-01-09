// -------------------------------------------
// @file      : simulator.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/7 下午6:17
// -------------------------------------------

package simulator

import (
	"go.uber.org/zap"
	"gogs/base/cberrors"
	"gogs/base/cluster"
	"gogs/base/cluster/network"
	"gogs/base/config"
	log "gogs/base/logger"
	"gogs/cb"
	"strconv"
)

func Main() {
	// 日志配置
	logConfig := config.GetLogConfig()
	// 模拟器配置
	simulatorConfig := config.GetSimulatorConfig()
	if logConfig == nil {
		cberrors.Panic("unable to find log config")
	}
	if simulatorConfig == nil {
		cberrors.Panic("unable to find game config")
	}
	// 初始化日志
	log.Init(
		log.SetFilename(simulatorConfig.LogPath),
		log.SetIsOpenFile(logConfig.IsOpenFile),
		log.SetIsOpenErrorFile(logConfig.IsOpenErrorFile),
		log.SetIsOpenConsole(logConfig.IsOpenConsole),
		log.SetIsAsync(logConfig.IsAsync),
		log.SetMaxFileSize(int(logConfig.Maxsize)),
		log.SetStacktrace(zap.PanicLevel),
	)
	builders := map[string]cluster.IServiceBuilder{
		"gate": cb.NewGateBuilder(nil),
		"game": cb.NewGameBuilder(nil),
		"client": cb.NewClientAPIBuilder(func(service cluster.IService) (cb.IClientAPI, error) {
			return NewClientAPI(), nil
		}),
	}
	simulator, err := cluster.NewSimulator(
		"localhost:9100",
		builders,
		network.ProtocolTCP,
	)
	if err != nil {
		log.Errorf("new simulator err:%s", err)
		return
	}
	userID := 1
	client, err := simulator.Connect(strconv.Itoa(userID))
	if err != nil {
		log.Errorf("connect err:%s", err)
		return
	}
	ack, code, err := client.GateServer.(*cb.GateRemoteService).Login(&cb.LoginReq{
		AccountID:   100,
		Token:       "abc",
		UserID:      999,
		ServerID:    1,
		AccountType: cb.AccountTypeTest,
	}, &cb.ClientInfo{})
	if err != nil {
		log.Errorf("login err:%s", err)
		return
	}
	log.Debugf("login ack:%+v code:%d", ack, code)
	<-make(chan struct{})
}
