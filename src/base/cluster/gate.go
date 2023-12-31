// -------------------------------------------
// @file      : Gate.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 上午11:15
// -------------------------------------------

package cluster

import (
	"gogs/base/cberrors"
	"gogs/base/cluster/network"
	log "gogs/base/logger"
	"sync"
	"sync/atomic"
)

// Gate 网关服务器
type Gate struct {
	*RPC                              // RPC管理器
	sync.RWMutex                      // 读写锁
	name         string               // 网关名字
	host         *Host                // 集群服务器
	gameServers  map[ID]IService      // Game以GameServer形式,保存在Gate
	agents       map[int64]*GateAgent // GateAgent列表,通过UserID索引
	builder      IServiceBuilder
	idgen        int64 // session userID generator
}

// NewGate 新建网关 localAddr本地对客户端监听地址,hostAddr集群节点地址
func NewGate(name, localAddr, hostAddr string, builder IServiceBuilder, protocol network.ProtocolType) (*Gate, error) {
	gate := &Gate{
		name:        name,
		gameServers: make(map[ID]IService),
		agents:      make(map[int64]*GateAgent),
		builder:     builder,
	}
	// 为网关创建集群节点服务器
	gate.host = NewHost(hostAddr)
	// 注册GateDriver
	sessionHandlerBuilder := func(session network.ISession) (network.ISessionHandler, error) {
		return newGateAgent(gate, session, gate.GenSessionID())
	}
	err := gate.host.Node.NewDriver(network.NewGateDriver(localAddr, sessionHandlerBuilder, protocol))
	if err != nil {
		return nil, err
	}
	// 参数为nil,表示只能构造远程服务
	gameServerBuilder := NewGameServerBuilder(nil)
	_, err = gate.host.RegisterBuilder(gameServerBuilder)
	if err != nil {
		return nil, err
	}
	// 在集群内通过服务类型监听GameServer服务
	listener := func(service IService, status network.ServiceStatus) bool {
		gate.Lock()
		defer gate.Unlock()
		if status == network.ServiceStatusOnline {
			gate.gameServers[service.ID()] = service
			log.Infof("service GameServerRemoteService online name: %s, type: %s, id: %d, remote: %d",
				service.Name(), service.Type(), service.ID(), service.(IRemoteService).RemoteID())
			log.Infof("service GameServerRemoteService list: %v", gate.gameServers)
		} else {
			delete(gate.gameServers, service.ID())
			log.Infof("service GameServerRemoteService offline name: %s, type: %s, id: %d, remote: %d",
				service.Name(), service.Type(), service.ID(), service.(IRemoteService).RemoteID())
		}
		return true
	}
	gate.host.ListenServiceType(GameServerTypeName, listener)
	// 新建GateServerBuilder,创建本地GateServer服务时,把Gate作为服务实际提供者
	gateServerBuilder := NewGateServerBuilder(func(service IService) (IGateServer, error) {
		return gate, nil
	})
	_, err = gate.host.RegisterBuilder(gateServerBuilder)
	if err != nil {
		return nil, err
	}
	// 创建一个本地GateServer服务
	_, err = gate.host.NewLocalService(GateServerTypeName, name)
	if err != nil {
		return nil, err
	}
	return gate, nil
}

// Close 关闭网关
func (gate *Gate) Close() {
	gate.host.Close()
}

// String implement fmt.Stringer
func (gate *Gate) String() string {
	return gate.name
}

// Name implement IService
func (gate *Gate) Name() string {
	return gate.name
}

// GenSessionID 生成唯一会话ID
func (gate *Gate) GenSessionID() int64 {
	return atomic.AddInt64(&gate.idgen, 1)
}

// sessionStatusChanged 会话状态变化
func (gate *Gate) sessionStatusChanged(agent *GateAgent, status network.SessionStatus) {
	gate.Lock()
	defer gate.Unlock()
	if status == network.SessionStatusInConnected {
		gate.agents[agent.userID] = agent
	} else {
		delete(gate.agents, agent.userID)
	}
}

// Login 用户登录
func (gate *Gate) Login(agent *GateAgent, ntf *UserLoginNtf, ci *ClientInfo) (Err, error) {
	ntf.SessionID = agent.sessionID
	ntf.Gate = gate.name

	gate.Lock()
	defer gate.Unlock()
	// TODO 这里
	for _, service := range gate.gameServers {
		if ntf.ServerID == 0 || GetIDByName(service.Name()) == ntf.ServerID {
			gameServer := service.(IGameServer)
			userID, code, err := gameServer.Login(ntf, ci)
			if err != nil || code != ErrOK {
				return code, cberrors.New("call GameServer#Login(%s) code: %s, err: %v", gameServer, code, err)
			}
			agent.gameServer = gameServer
			agent.userID = userID
			go gate.sessionStatusChanged(agent, network.SessionStatusInConnected)
			return ErrOK, nil
		}
	}
	return ErrOK, cberrors.New("unable to find game server: %d", ntf.ServerID)
}

// Tunnel 转发从Game->Client的消息,通过UserID找到对应的Session
func (gate *Gate) Tunnel(msg *TunnelMsg) error {
	gate.RLock()
	remote, ok := gate.agents[msg.UserID]
	gate.RUnlock()
	if ok {
		return remote.session.Write(&network.Message{
			Type: msg.Type,
			Data: msg.Data,
		})
	}
	return nil
}
