// -------------------------------------------
// @file      : simulator_agent.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/7 下午5:03
// -------------------------------------------

package cluster

import (
	"gogs/base/cluster/network"
	log "gogs/base/logger"
	"time"
)

// SimulatorAgent 模拟器代理
// Simulator 和 network.ClientSession 的中间层
type SimulatorAgent struct {
	simulator *Simulator
	session   network.ISession
	client    *Client
}

// NewSimulatorAgent 新建模拟器代理
func NewSimulatorAgent(simulator *Simulator, session network.ISession) *SimulatorAgent {
	return &SimulatorAgent{
		simulator: simulator,
		session:   session,
	}
}

// Name 唯一标识,其代理会话的名字
func (agent *SimulatorAgent) Name() string {
	return agent.session.Name()
}

// SetClient 设置client
func (agent *SimulatorAgent) SetClient(client *Client) {
	agent.client = client
}

// Client 获取client
func (agent *SimulatorAgent) Client() *Client {
	return agent.client
}

// Post implements IAgent
func (agent *SimulatorAgent) Post(service IService, call *network.Call) error {
	return agent.simulator.Post(agent.session, call)
}

// Wait implements IAgent
func (agent *SimulatorAgent) Wait(service IService, call *network.Call, timeout time.Duration) (Future, error) {
	return agent.simulator.Wait(agent.session, call, timeout)
}

// Write implements IAgent
func (agent *SimulatorAgent) Write(msg *network.Message) error {
	return agent.session.Write(msg)
}

// Session implements IAgent
func (agent *SimulatorAgent) Session() network.ISession {
	return agent.session
}

// Close implements IAgent
func (agent *SimulatorAgent) Close() {
}

// SessionStatusChanged implements network.ISessionHandler
func (agent *SimulatorAgent) SessionStatusChanged(status network.SessionStatus) {
	agent.simulator.sessionStatusChanged(agent, status)
}

// Read implements network.ISessionHandler
func (agent *SimulatorAgent) Read(session network.ISession, msg *network.Message) {
	switch msg.Type {
	case network.MessageTypeCall:
		go agent.handleCall(msg.Data)
	case network.MessageTypeReturn:
		go agent.handleReturn(msg.Data)
	}
}

// handleCall 处理调用
func (agent *SimulatorAgent) handleCall(data []byte) {
	call, err := network.UnmarshalCall(data)
	if err != nil {
		log.Error("%s", err)
		return
	}
	callReturn, err := agent.client.ClientService.Call(call)
	if err != nil {
		log.Error("handle call service: %d, method: %d, err: %s", call.ServiceID, call.MethodID, err)
		return
	}
	if callReturn == nil {
		return
	}
	data = callReturn.Marshal()
	msg := &network.Message{
		Type: network.MessageTypeReturn,
		Data: data,
	}
	err = agent.session.Write(msg)
	if err != nil {
		log.Error("client session write msg err: %s", err)
		return
	}
	return
}

// handleReturn 处理调用返回
func (agent *SimulatorAgent) handleReturn(data []byte) {
	callReturn, err := network.UnmarshalReturn(data)
	if err != nil {
		log.Error("unmarshal return err: %s", err)
		return
	}
	agent.simulator.Notify(callReturn)
}
