// -------------------------------------------
// @file      : login.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午2:54
// -------------------------------------------

package cluster

import (
	"fmt"
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
		UserServiceName: "client",
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
		log.Debugf("current Gate servers: %v", game.gateServers)
		return true
	}
	game.Host.ListenServiceType(GateServerTypeName, listener)
	return game, nil
}

// Shutdown 关闭服务器
func (game *Game) Shutdown() {
	log.Infof("%s shutdown start:", game.serverName)
	log.Infof("%s:Host closing...", game.serverName)
	game.Host.Close()
	log.Infof("%s:ActorSystem closing...", game.serverName)
	game.ActorSystem.Close()
	etcd.Exit()
	log.Infof("%s shutdown finished.", game.serverName)
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
func (game *Game) Login(ntf *UserLoginNtf, ci *ClientInfo) (int64, Err, error) {
	game.Lock()
	gateServer, ok := game.gateServers[ntf.Gate]
	game.Unlock()
	if !ok {
		log.Errorf("unable to find gate server: %s", ntf.Gate)
		return 0, ErrGateNotFound, nil
	}
	clientType := "client"
	builder, ok := game.builders[clientType]
	if !ok {
		log.Errorf("unable to find client type builder: %s", clientType)
		return 0, ErrSystem, nil
	}
	// 角色名字
	name := ActorName{
		SystemName: game.ActorSystem.name,
		Type:       game.UserServiceName,
		ID:         ntf.UserID,
	}
	// 客户端代理
	var clientAgent *ClientAgent
	// 在角色系统查找是否已经存在角色
	actor, ok := game.ActorSystem.GetActor(name.String())
	if !ok {
		clientAgent = NewClientAgent(ntf.SessionID, ntf.UserID)
		var err error
		actor, err = game.ActorSystem.NewActor(name, clientAgent)
		if actor == nil {
			log.Errorf("new actor: %s err: %s", name.String(), err)
			return 0, ErrActorName, nil
		}
	} else {
		clientAgent = actor.Context().(*ClientAgent)
		clientAgent.sessionID = ntf.SessionID
	}
	game.Lock()
	defer game.Unlock()
	remoteService := builder.NewRemoteService(newTunnelAgent(game, clientAgent.userID, gateServer),
		actor.Name(),
		game.newServiceID(),
		0,
		nil)
	clientAgent.SetClientService(remoteService)
	return clientAgent.UserID(), ErrOK, nil
}

// Logout 登出
func (game *Game) Logout(ntf *UserLoginNtf) error {
	actorName := fmt.Sprintf("%s:%s@%d", game.ActorSystem.name, game.UserServiceName, ntf.UserID)
	if actor, ok := game.ActorSystem.GetActor(actorName); ok {
		clientAgent := actor.Context().(*ClientAgent)
		if clientAgent.sessionID != ntf.SessionID {
			log.Infof("%s logout %s %s", clientAgent, clientAgent.sessionID, ntf.SessionID)
			return nil
		}
		log.Infof("%s logout", clientAgent)
		game.Lock()
		clientAgent.SetClientService(nil)
		game.Unlock()
	}
	return nil
}

// Tunnel 处理gate转发来的client消息
func (game *Game) Tunnel(msg *TunnelMsg) error {
	if msg.Type == network.MessageTypeReturn {
		callReturn, err := network.UnmarshalReturn(msg.Data)
		if err != nil {
			return err
		}
		game.Notify(callReturn)
		return nil
	}
	// 如果是调用,查找用户角色
	actorName := fmt.Sprintf("%s:%s@%d", game.ActorSystem.name, game.UserServiceName, msg.UserID)
	if actor, ok := game.ActorSystem.GetActor(actorName); ok {
		// 获取角色上下文
		clientAgent := actor.Context().(*ClientAgent)
		call, err := network.UnmarshalCall(msg.Data)
		if err != nil {
			return err
		}
		// 调用角色方法
		callReturn, err := actor.Service().Call(call)
		if err != nil {
			return err
		}
		if callReturn == nil {
			return nil
		}
		data := callReturn.Marshal()
		if clientService, ok := clientAgent.ClientService(); ok {
			msg := &network.Message{
				Type: network.MessageTypeReturn,
				Data: data,
			}
			err = clientService.Agent().Write(msg)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return cberrors.New("actor not found: %s", actorName)
}
