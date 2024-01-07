// -------------------------------------------
// @file      : game.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午2:54
// -------------------------------------------

package cluster

import (
	"gogs/base/cberrors"
	"gogs/base/cluster/network"
	"gogs/base/etcd"
	log "gogs/base/logger"
	"sync"
	"sync/atomic"
)

// Game 游戏服务器
type Game struct {
	*RPC                                       // RPC管理器
	sync.RWMutex                               // 读写锁
	ActorSystem     *ActorSystem               // 角色系统
	Host            *Host                      // 集群服务器
	gateServers     map[string]IGateServer     // 网关服务集合
	builders        map[string]IServiceBuilder // 服务构造器集合
	idgen           uint32                     // service userID generator
	UserServiceName string                     // 实际提供API服务的本地服务名字
	serverID        int64                      // 服务器ID
	serverName      string                     // 服务器名字
}

// NewGame 新建游戏服务器
func NewGame(name string, builders map[string]IServiceBuilder, localAddr string) (
	*Game, error) {
	actorSystem, err := NewActorSystem(name, builders, localAddr)
	if err != nil {
		return nil, err
	}
	game := &Game{
		RPC:             NewRPC(),
		ActorSystem:     actorSystem,
		Host:            actorSystem.host,
		gateServers:     make(map[string]IGateServer),
		builders:        builders,
		UserServiceName: "GameAPI",
		serverName:      name,
	}
	// 注册GameServer服务构造器
	_, err = game.Host.RegisterBuilder(NewGameServerBuilder(func(service IService) (IGameServer, error) {
		return game, nil
	}))
	if err != nil {
		return nil, err
	}
	// 创建本地GameServer服务时,返回自身
	_, err = game.Host.NewLocalService(GameServerTypeName, name)
	if err != nil {
		return nil, err
	}
	// 注册GateServer服务构造器,只能构造GateServerRemoteService
	_, err = game.Host.RegisterBuilder(NewGateServerBuilder(nil))
	// 监听GateServer类型服务
	listener := func(service IService, status network.ServiceStatus) bool {
		game.Lock()
		defer game.Unlock()
		if status == network.ServiceStatusOnline {
			game.gateServers[service.Name()] = service.(IGateServer)
			log.Infof("service GateServerRemoteService online name: %s, type: %s, id: %d, remote id: %d",
				service.Name(), service.Type(), service.ID(), service.(IRemoteService).RemoteID())
		} else {
			delete(game.gateServers, service.Name())
			log.Infof("service GateServerRemoteService  offline name: %s, type: %s, id: %d, remote id: %d",
				service.Name(), service.Type(), service.ID(), service.(IRemoteService).RemoteID())
		}
		log.Debugf("current gate servers: %v", game.gateServers)
		return true
	}
	game.Host.ListenServiceType(GateServerTypeName, listener)
	return game, nil
}

// Shutdown 关闭服务器
func (game *Game) Shutdown() {
	log.Info("Game shutdown start:")
	log.Info("Game:Host closing...")
	game.Host.Close()
	log.Info("Game:ActorSystem closing...")
	game.ActorSystem.Close()
	etcd.Exit()
	log.Info("Game:DB closing...")

	log.Info("Game shutdown finished.")
}

// newServiceID 生成一个唯一ID
func (game *Game) newServiceID() ID {
	return ID(atomic.AddUint32(&game.idgen, 1))
}

// Name 服务器名字
func (game *Game) Name() string {
	return game.serverName
}

// Login 登录
func (game *Game) Login(msg *RProxyMsg) (int64, Err, error) {
	game.Lock()
	gateServer, ok := game.gateServers[msg.Gate]
	game.Unlock()
	if !ok {
		return 0, ErrRProxy, nil
	}
	clientType := "client"
	builder, ok := game.builders[clientType]
	if !ok {
		return 0, ErrUnknownService, cberrors.New("unable to find client type builder: %s", clientType)
	}
	// 角色名字
	name := ActorName{
		SystemName: game.ActorSystem.name,
		Type:       game.UserServiceName,
		ID:         msg.UserID,
	}
	// 客户端代理
	var clientAgent *ClientAgent
	// 在角色系统查找是否已经存在角色
	actor, ok := game.ActorSystem.GetActor(name.String())
	if !ok {
		clientAgent = NewClientAgent(msg.SessionID, msg.UserID)
		var err error
		actor, err = game.ActorSystem.NewActor(name, clientAgent)
		if actor == nil {
			log.Errorf("new actor: %s err: %s", name, err)
			return 0, ErrActorName, err
		}
	} else {
		clientAgent = actor.Context().(*ClientAgent)
		clientAgent.sessionID = msg.SessionID
	}
	game.Lock()
	defer game.Unlock()
	remoteService := builder.NewRemoteService(newTunnelRemote(game, clientAgent.userID, gateServer),
		actor.Name(),
		game.newServiceID(),
		0,
		nil)
	clientAgent.SetClientService(remoteService)
	return clientAgent.UserID(), ErrOK, nil
}

// Logout 登出
func (game *Game) Logout(msg *RProxyMsg) error {
	return nil
}

func (game *Game) Tunnel(msg *TunnelMsg) error {
	return nil
}
