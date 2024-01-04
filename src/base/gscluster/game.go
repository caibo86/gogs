// -------------------------------------------
// @file      : game.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午2:54
// -------------------------------------------

package gscluster

import (
	"fmt"
	"gogs/base/gserrors"
	"gogs/base/gsnet"
	log "gogs/base/logger"
	"gogs/base/mongodb"
	"sync"
	"sync/atomic"
)

// Game 游戏服务器
type Game struct {
	*RPC                                    // RPC管理器
	sync.RWMutex                            // 读写锁
	ActorSystem  *ActorSystem               // 角色系统
	host         *Host                      // 集群服务器
	db           *mongodb.MongoClient       // mongodb
	gateServers  map[string]IGateServer     // 网关服务集合
	builders     map[string]IServiceBuilder // 服务构造器集合
	idgen        uint32                     // service id generator
	tUser        string                     // 玩家类型
	serverID     int64                      // 服务器ID
	serverName   string                     // 服务器名字
}

// NewGame 新建游戏服务器
func NewGame(id int64, name string, builders map[string]IServiceBuilder, localAddr string, db *mongodb.MongoClient) (
	*Game, error) {
	actorSystem, err := NewActorSystem(name, builders, localAddr, db)
	if err != nil {
		return nil, err
	}
	game := &Game{
		RPC:         NewRPC(),
		ActorSystem: actorSystem,
		host:        actorSystem.host,
		db:          db,
		gateServers: make(map[string]IGateServer),
		builders:    builders,
		tUser:       "User",
		serverID:    id,
		serverName:  name,
	}
	// 注册GameServer服务构造器
	_, err = game.host.RegisterBuilder(NewGameServerBuilder(func(service IService) (IGameServer, error) {
		return game, nil
	}))
	if err != nil {
		return nil, err
	}
	// 创建本地GameServer服务时,返回自身
	_, err = game.host.NewLocalService(GameServerTypeName, fmt.Sprintf("%s@%d", name, id))
	if err != nil {
		return nil, err
	}
	// 注册GateServer服务构造器,只能构造GateServerRemoteService
	_, err = game.host.RegisterBuilder(NewGateServerBuilder(nil))
	// 监听GateServer类型服务
	listener := func(service IService, status gsnet.ServiceStatus) bool {
		game.Lock()
		defer game.Unlock()
		if status == gsnet.ServiceStatusOnline {
			game.gateServers[service.Name()] = service.(IGateServer)
			log.Infof("gate server online name: %s type: %s id: %d",
				service.Name(), service.Type(), service.ID())
		} else {
			delete(game.gateServers, service.Name())
			log.Infof("gate server offline name: %s type: %s id: %d",
				service.Name(), service.Type(), service.ID())
		}
		log.Debugf("current gate servers: %v", game.gateServers)
		return true
	}
	game.host.ListenServiceType(GateServerTypeName, listener)
	return game, nil
}

// Shutdown 关闭服务器
func (game *Game) Shutdown() {
	log.Info("Game shutdown start:")
	log.Info("Game:host closing...")
	game.host.Close()
	log.Info("Game:ActorSystem closing...")
	game.ActorSystem.Close()
	log.Info("Game:DB closing...")
	err := game.db.Disconnect()
	if err != nil {
		log.Errorf("mongodb disconnect err: %s", err)
	}
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
		return 0, ErrUnknownService, gserrors.Newf("unregister client type builder: %s", clientType)
	}
	// 角色名字
	name := ActorName{
		SystemName: game.ActorSystem.name,
		Type:       game.tUser,
		ID:         msg.UserID,
	}
	// 用户对象
	var user *User
	// 在角色系统查找是否已经存在角色
	actor, ok := game.ActorSystem.GetActor(name.String())
	if !ok {
		user = NewUser(msg.SessionID, msg.UserID)
		var err error
		actor, err = game.ActorSystem.NewActor(name, user)
		if actor == nil {
			log.Errorf("new actor: %s err: %s", name, err)
			return 0, ErrActorName, err
		}
	} else {
		user = actor.Context().(*User)
		user.sessionID = msg.SessionID
	}
	game.Lock()
	defer game.Unlock()
	remoteService := builder.NewRemoteService(newTunnelRemote(game, user.id, gateServer),
		actor.Name(),
		game.newServiceID(),
		0,
		nil)
	user.SetClient(remoteService)
	return user.ID(), ErrOK, nil
}

// Logout 登出
func (game *Game) Logout(msg *RProxyMsg) error {
	return nil
}

func (game *Game) Tunnel(msg *TunnelMsg) error {
	return nil
}
