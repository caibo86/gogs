// -------------------------------------------
// @file      : simulator.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/7 下午5:03
// -------------------------------------------

package cluster

import (
	"gogs/base/cberrors"
	"gogs/base/cluster/network"
	log "gogs/base/logger"
	"strconv"
	"sync/atomic"
	"time"
)

type Client struct {
	Name          string
	ClientService IService
	GateServer    IRemoteService
	GameServer    IRemoteService
	agent         *SimulatorAgent
}

// Agent 获取模拟器代理
func (client *Client) Agent() *SimulatorAgent {
	return client.agent
}

// Simulator 客户端模拟器
type Simulator struct {
	*RPC                                               // RPC集中管理器
	*ServiceStatusPublisher                            // 服务状态发布器
	driver                  network.IDriver            // 驱动
	builders                map[string]IServiceBuilder // 服务构造者集合
	idgen                   uint32                     // serviceID生成器
}

// NewSimulator 生成客户端模拟器
func NewSimulator(remoteAddr string, builders map[string]IServiceBuilder, protocol network.ProtocolType) (*Simulator, error) {
	simulator := &Simulator{
		RPC:                    NewRPC(),
		ServiceStatusPublisher: NewServiceStatusPublisher(),
		builders:               builders, // 指定建造者集合
	}
	// 客户端驱动
	simulator.driver = network.NewClientDriver(
		remoteAddr,
		func(session network.ISession) (network.ISessionHandler, error) {
			return NewSimulatorAgent(simulator, session), nil
		},
		protocol,
	)
	return simulator, nil
}

// Run 启动客户端模拟器
func (simulator *Simulator) Run(nums int) error {
	// 开启指定数量的连接
	for i := 0; i < nums; i++ {
		_, err := simulator.driver.NewSession(strconv.Itoa(i), network.ConnectionTypeOut)
		if err != nil {
			return err
		}
	}
	<-make(chan bool)
	return nil
}

// Connect 连接
func (simulator *Simulator) Connect(name string) (*Client, error) {
	_, err := simulator.driver.NewSession(name, network.ConnectionTypeOut)
	if err != nil {
		return nil, err
	}
	result := make(chan IService)
	listener := func(service IService, status network.ServiceStatus) bool {
		if status == network.ServiceStatusOnline {
			result <- service
		}
		return false
	}
	simulator.ListenService(name, listener)
	select {
	case service := <-result:
		return service.Context().(*Client), nil
	case <-time.After(time.Second * 5):
		return nil, cberrors.New("simulator connect timeout")
	}
}

// newServiceID 生成一个新的服务ID
func (simulator *Simulator) newServiceID() ID {
	return ID(atomic.AddUint32(&simulator.idgen, 1))
}

// sessionStatusChanged 会话状态改变
func (simulator *Simulator) sessionStatusChanged(agent *SimulatorAgent, status network.SessionStatus) {
	switch status {
	case network.SessionStatusOutConnected:
		// 连接成功生成一个客户端
		csBuilder, ok := simulator.builders["client"]
		if !ok {
			log.Debug("unable to find a builder for service type: client")
			return
		}
		gateServerBuilder, ok := simulator.builders["Gate"]
		if !ok {
			log.Debug("unable to find a builder for service type: Gate")
			return
		}
		gameServerBuilder, ok := simulator.builders["game"]
		if !ok {
			log.Debug("unable to find a builder for service type: game")
			return
		}
		client := &Client{
			Name:       agent.Name(),
			agent:      agent,
			GateServer: gateServerBuilder.NewRemoteService(agent, "Gate", simulator.newServiceID(), gateID, nil),
			GameServer: gameServerBuilder.NewRemoteService(agent, "game", simulator.newServiceID(), gameID, nil),
		}
		clientService, err := csBuilder.NewService(agent.Name(), simulator.newServiceID(), client)
		if err != nil {
			log.Errorf("create local client service err: %s", err)
			return
		}
		client.ClientService = clientService
		agent.SetClient(client)
		simulator.ServiceStatusChanged(clientService, network.ServiceStatusOnline)
	case network.SessionStatusDisconnected:
		// 通知连接断开
		if client := agent.Client(); client != nil {
			simulator.ServiceStatusChanged(client.ClientService, network.ServiceStatusOffline)
		}
	}
}
