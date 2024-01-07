// -------------------------------------------
// @file      : actor_system.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/4 下午2:56
// -------------------------------------------

package cluster

import (
	"fmt"
	"gogs/base/cberrors"
	"gogs/base/cluster/network"
	"gogs/base/config"
	log "gogs/base/logger"
	"gogs/base/mongodb"
	"sync"
	"sync/atomic"
)

// ActorSystem 角色系统
type ActorSystem struct {
	*RPC                                    // RPC管理器
	name         string                     // 名字
	host         *Host                      // 集群服务器
	builders     map[string]IServiceBuilder // 服务构造器集合
	actors       map[string]IActor          // 角色集合
	actorLock    sync.RWMutex               // 角色集合读写锁
	neighbors    map[string]IActorSystem    // 邻居角色系统集合
	neighborLock sync.RWMutex               // 邻居角色系统集合读写锁
	idgen        uint32                     // userID generator
	groupLocks   []sync.Mutex               // 分组锁
	db           *mongodb.MongoClient       // 数据库
}

// NewActorSystem 新建角色系统
func NewActorSystem(name string, builders map[string]IServiceBuilder, localAddr string) (*ActorSystem, error) {
	systemName := fmt.Sprintf("%s:ActorSystem", name)
	system := &ActorSystem{
		RPC:        NewRPC(),
		name:       systemName,
		host:       NewHost(localAddr),
		builders:   builders,
		actors:     make(map[string]IActor),
		neighbors:  make(map[string]IActorSystem),
		groupLocks: make([]sync.Mutex, config.ActorGroups()),
	}
	// 注册服务构造器,本地IActorSystem服务即返回自身
	_, err := system.host.RegisterBuilder(NewActorSystemBuilder(func(service IService) (IActorSystem, error) {
		return system, nil
	}))
	if err != nil {
		return nil, err
	}
	_, err = system.host.NewLocalService(ActorSystemTypeName, systemName)
	if err != nil {
		return nil, err
	}
	listener := func(service IService, status network.ServiceStatus) bool {
		system.neighborLock.Lock()
		defer system.neighborLock.Unlock()
		if status == network.ServiceStatusOnline {
			system.neighbors[service.Name()] = service.(IActorSystem)
			log.Infof("neighbor actor system online name: %s type: %s userID: %d",
				service.Name(), service.Type(), service.ID())
		} else {
			delete(system.neighbors, service.Name())
			log.Infof("neighbor actor system offline name: %s type: %s userID: %d",
				service.Name(), service.Type(), service.ID())
		}
		return true
	}
	system.host.ListenServiceType(ActorSystemTypeName, listener)
	return system, nil
}

// Close 关闭角色系统
func (system *ActorSystem) Close() {
	system.actorLock.Lock()
	defer system.actorLock.Unlock()
	for _, actor := range system.actors {
		if context := actor.Context(); context != nil {
			log.Infof("actor: %s save context", actor.Name())
			err := context.Save(system.db)
			if err != nil {
				log.Errorf("actor: %s save context failed: %s", actor.Name(), err)
			}
		}
	}
}

// CheckNeighborExist 检查指定的邻居是否上线
func (system *ActorSystem) CheckNeighborExist(name string) bool {
	system.neighborLock.RLock()
	defer system.neighborLock.RUnlock()
	_, ok := system.neighbors[name]
	return ok
}

// newServiceID 获取一个唯一ID,用于ServiceID,非角色ID
func (system *ActorSystem) newServiceID() ID {
	return ID(atomic.AddUint32(&system.idgen, 1))
}

// DelActor 在角色集合删除指定的角色
func (system *ActorSystem) DelActor(actor IActor) {
	system.actorLock.Lock()
	defer system.actorLock.Unlock()
	if u := system.actors[actor.Name()]; u == actor {
		delete(system.actors, actor.Name())
	}
}

// GetActor 获取指定名字的角色
func (system *ActorSystem) GetActor(name string) (IActor, bool) {
	system.actorLock.RLock()
	defer system.actorLock.RUnlock()
	actor, ok := system.actors[name]
	return actor, ok
}

// NewActor 创建角色,可能有本地角色和远程角色
func (system *ActorSystem) NewActor(name ActorName, context IActorContext) (IActor, error) {
	builder, ok := system.builders[name.Type]
	if !ok {
		return nil, cberrors.New("actor builder not found: %s", name.Type)
	}
	if context == nil {
		context = NewNilActorContext()
	}
	if name.SystemName == "" || name.SystemName == system.name {
		name.SystemName = system.name
		nameStr := name.String()
		serviceID := system.newServiceID()
		locker := &system.groupLocks[serviceID%ID(len(system.groupLocks))]
		actor := newBaseActor(system, &name, locker, context)
		context.SetActor(actor)
		service, err := builder.NewService(nameStr, serviceID, context)
		if err != nil {
			return nil, err
		}
		actor.service = service
		system.actorLock.Lock()
		if old, ok := system.actors[nameStr]; ok {
			system.actorLock.Unlock()
			return old, cberrors.New("actor: %s already exist", actor.Name())
		}
		system.actors[actor.Name()] = actor
		system.actorLock.Unlock()
		return actor, nil
	}
	system.neighborLock.RLock()
	neighbor, ok := system.neighbors[name.SystemName]
	system.neighborLock.RUnlock()
	if ok {
		serviceID := system.newServiceID()
		locker := &system.groupLocks[serviceID%ID(len(system.groupLocks))]
		nameStr := name.String()
		remote := newActorRemote(system, nameStr, neighbor)
		actor := newBaseActor(system, &name, locker, context)
		context.SetActor(actor)
		actor.service = builder.NewRemoteService(remote, nameStr, serviceID, 0, actor)
		system.actorLock.Lock()
		if old, ok := system.actors[nameStr]; ok {
			system.actorLock.Unlock()
			return old, cberrors.New("actor: %s already exist", actor.Name())
		}
		system.actors[nameStr] = actor
		system.actorLock.Unlock()
		return actor, nil
	}
	return nil, cberrors.New("neighbor actor system not found: %s", name.SystemName)
}

// ActorInvoke 角色调用
func (system *ActorSystem) ActorInvoke(msg *ActorMsg) (*network.Return, Err, error) {
	system.actorLock.RLock()
	actor, ok := system.actors[msg.ActorName]
	system.actorLock.RUnlock()
	if !ok {
		return nil, ErrActorNotFound, nil
	}
	call, err := network.UnmarshalCall(msg.Data)
	if err != nil {
		return nil, ErrUnmarshal, err
	}
	callReturn, err := actor.Service().Call(call)
	if err != nil {
		return nil, ErrSystem, err
	}
	if callReturn == nil {
		callReturn = &network.Return{
			ID:        call.ID,
			ServiceID: call.ServiceID,
		}
	}
	return callReturn, ErrOK, nil
}
