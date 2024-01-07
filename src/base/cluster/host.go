// -------------------------------------------
// @file      : cluster.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/29 下午6:28
// -------------------------------------------

package cluster

import (
	"gogs/base/cberrors"
	"gogs/base/cluster/network"
	"gogs/base/config"
	log "gogs/base/logger"
	"sync"
	"sync/atomic"
	"time"
)

// Host 集群中的服务器
type Host struct {
	*RPC                                                    // 远程调用管理器
	*ServiceStatusPublisher                                 // 服务状态发布器
	wg                      sync.WaitGroup                  // WaitGroup
	idgen                   uint32                          // 服务ID生成器
	Node                    *network.Node                   // 网络层节点
	neighborMutex           sync.RWMutex                    // 集群中的邻居集合 读写锁
	neighbors               map[string]*Neighbor            // 集群中的邻居集合,通过Session.Name索引
	localServiceMutex       sync.RWMutex                    // 本地服务集合 读写锁
	localServices           map[ID]IService                 // 本地服务集合,通过ID索引
	builderMutex            sync.RWMutex                    // 服务构造器集合 读写锁
	builders                map[string]IServiceBuilder      // 服务构造器集合,通过ServiceType索引
	localServiceEvents      chan *eventServiceStatusChanged // 本地服务状态变更事件通道
	registryExit            chan struct{}                   // 关闭服务注册的信号
}

// NewHost 新建集群服务器
func NewHost(localAddr string) *Host {
	host := &Host{
		RPC:                    NewRPC(),
		ServiceStatusPublisher: NewServiceStatusPublisher(),
		Node:                   network.NewNode(),
		localServices:          make(map[ID]IService),
		neighbors:              make(map[string]*Neighbor),
		builders:               make(map[string]IServiceBuilder),
		localServiceEvents:     make(chan *eventServiceStatusChanged),
		registryExit:           make(chan struct{}, 1),
	}
	_ = host.Node.NewDriver(
		network.NewClusterDriver(
			localAddr,
			func(session network.ISession) (network.ISessionHandler, error) {
				return NewClusterRemote(host, session), nil
			}))
	go func() {
		// 启动定时器,定时向邻居节点注册服务
		host.wg.Add(1)
		ticker := time.Tick(config.ClusterRegistryInterval())
		for range ticker {
			if !host.tryServiceRegistry() {
				break
			}
		}
		log.Debug("Host service registry ticker exit")
		host.wg.Done()
	}()
	return host
}

// tryServiceRegistry 尝试向邻居节点注册服务
func (host *Host) tryServiceRegistry() bool {
	events := make(map[ID]*eventServiceStatusChanged)
	var content []*network.ServiceRegistry
	for {
		if len(events) > config.ClusterRegistryMax() {
			break
		}
		select {
		case event, ok := <-host.localServiceEvents:
			if !ok {
				return false
			}
			events[event.service.ID()] = event
		case <-host.registryExit:
			// 收到关闭服务注册的信号
			if len(events) > 0 {
				// 还有未处理的事件,继续处理
				host.registryExit <- struct{}{}
				goto DoRegistry
			}
			return false
		default:
			if len(events) == 0 {
				return true
			}
			goto DoRegistry
		}
	}
DoRegistry:
	for _, event := range events {
		registry := &network.ServiceRegistry{
			ServiceID:   uint32(event.service.ID()),
			ServiceType: event.service.Type(),
			ServiceName: event.service.Name(),
		}
		if event.status == network.ServiceStatusOnline {
			registry.Add = true
		} else {
			registry.Add = false
		}
		content = append(content, registry)
	}
	srd := &network.ServiceRegistryData{
		Data: content,
	}
	data := srd.Marshal()
	msg := &network.Message{
		Type: network.MessageTypeRegistry,
		Data: data,
	}
	host.neighborMutex.RLock()
	defer host.neighborMutex.RUnlock()
	for _, neighbor := range host.neighbors {
		err := neighbor.clusterRemote.Write(msg)
		if err != nil {
			log.Errorf("neighbors clusterRemote write msg err: %s", err)
		}
	}
	log.Info("Host service registry data sent")
	return true
}

// newID 生成唯一ServiceID
func (host *Host) newID() ID {
	for {
		val := atomic.AddUint32(&host.idgen, 1)
		if _, ok := host.localServices[ID(val)]; !ok {
			return ID(val)
		}
	}
}

// handleServiceRegistry 处理来自邻居节点的服务注册消息
func (host *Host) handleServiceRegistry(remote *HostAgent, srd *network.ServiceRegistryData) {
	host.neighborMutex.Lock()
	defer host.neighborMutex.Unlock()
	if neighbor, ok := host.neighbors[remote.Name()]; ok {
		for _, registry := range srd.Data {
			if registry.Add {
				// 服务注册
				if service, ok := neighbor.services[registry.ServiceName]; ok {
					if service.RemoteID() != ID(registry.ServiceID) || service.Type() != registry.ServiceType {
						delete(neighbor.services, registry.ServiceName)
						delete(neighbor.servicesByID, service.RemoteID())
						host.ServiceStatusChanged(service, network.ServiceStatusOffline)
					} else {
						continue
					}
				}
				host.builderMutex.RLock()
				builder, ok := host.builders[registry.ServiceType]
				host.builderMutex.RUnlock()
				if ok {
					remoteService := builder.NewRemoteService(
						remote, registry.ServiceName, host.newID(), ID(registry.ServiceID), nil,
					)
					neighbor.services[registry.ServiceName] = remoteService
					neighbor.servicesByID[remoteService.RemoteID()] = remoteService
					host.ServiceStatusChanged(remoteService, network.ServiceStatusOnline)
					continue
				}
			} else {
				// 服务注销
				if service, ok := neighbor.services[registry.ServiceName]; ok {
					if service.RemoteID() == ID(registry.ServiceID) && service.Type() == registry.ServiceType {
						delete(neighbor.services, registry.ServiceName)
						delete(neighbor.servicesByID, service.RemoteID())
						host.ServiceStatusChanged(service, network.ServiceStatusOffline)
					}
				}
			}
		}
	}
}

// sessionStatusChanged 会话状态变更
func (host *Host) sessionStatusChanged(remote *HostAgent, status network.SessionStatus) {
	switch status {
	case network.SessionStatusInConnected, network.SessionStatusOutConnected:
		// 邻居节点连接成功
		host.neighborMutex.Lock()
		host.neighbors[remote.Name()] = NewNeighbor(remote)
		host.neighborMutex.Unlock()
		// 向邻居节点注册本地服务
		var content []*network.ServiceRegistry
		host.localServiceMutex.RLock()
		for _, service := range host.localServices {
			registry := &network.ServiceRegistry{
				Add:         true,
				ServiceID:   uint32(service.ID()),
				ServiceType: service.Type(),
				ServiceName: service.Name(),
			}
			content = append(content, registry)
		}
		host.localServiceMutex.RUnlock()
		srd := &network.ServiceRegistryData{
			Data: content,
		}
		data := srd.Marshal()
		msg := &network.Message{
			Type: network.MessageTypeRegistry,
			Data: data,
		}
		err := remote.Write(msg)
		if err != nil {
			log.Errorf("clusterRemote:%s  write msg err: %s", remote, err)
		}
	case network.SessionStatusClosed, network.SessionStatusDisconnected:
		// 邻居节点断开连接
		host.neighborMutex.Lock()
		delete(host.neighbors, remote.Name())
		host.neighborMutex.Unlock()
	}
}

// handleCall 处理来自对本地服务的调用
func (host *Host) handleCall(call *network.Call) (*network.Return, error) {
	host.localServiceMutex.RLock()
	defer host.localServiceMutex.RUnlock()
	if service, ok := host.localServices[ID(call.ServiceID)]; ok {
		return service.Call(call)
	}
	return nil, cberrors.New("local service not found: %d", call.ServiceID)
}

// RegisterBuilder 注册服务构造器
func (host *Host) RegisterBuilder(builder IServiceBuilder) (IServiceBuilder, error) {
	host.builderMutex.Lock()
	defer host.builderMutex.Unlock()
	if _, ok := host.builders[builder.ServiceType()]; ok {
		return nil, cberrors.New("duplicate service builder: %s", builder)
	}
	host.builders[builder.ServiceType()] = builder
	return builder, nil
}

// UnregisterBuilder 注销服务构造器
func (host *Host) UnregisterBuilder(serviceType string) IServiceBuilder {
	host.builderMutex.Lock()
	defer host.builderMutex.Unlock()
	builder := host.builders[serviceType]
	delete(host.builders, serviceType)
	return builder
}

// Connect 对集群中指定地址的节点发起连接
func (host *Host) Connect(remoteAddr string) (IAgent, error) {
	session, err := host.Node.NewSession(network.DriverTypeCluster, remoteAddr, network.ConnectionTypeOut)
	if session == nil {
		return nil, err
	}
	return session.Handler().(IAgent), err
}

// NewService 新建本地服务
func (host *Host) NewService(serviceType string, name string, context interface{}) (IService, error) {
	host.localServiceMutex.RLock()
	// 检查是否有重名的服务
	for _, service := range host.localServices {
		if service.Name() == name {
			host.localServiceMutex.RUnlock()
			if service.Type() == serviceType {
				return service, nil
			}
			return service, cberrors.New("duplicate service name: %s with different type, expect: %s, found: %s",
				name, serviceType, service.Type())
		}
	}
	host.localServiceMutex.RUnlock()

	host.builderMutex.RLock()
	builder, ok := host.builders[serviceType]
	host.builderMutex.RUnlock()
	if !ok {
		return nil, cberrors.New("service builder not found for type: %s", serviceType)
	}
	service, err := builder.NewService(name, host.newID(), context)
	if err != nil {
		return nil, err
	}
	host.localServiceMutex.Lock()
	host.localServices[service.ID()] = service
	host.localServiceMutex.Unlock()
	host.localServiceEvents <- &eventServiceStatusChanged{
		service: service,
		status:  network.ServiceStatusOnline,
	}
	host.ServiceStatusChanged(service, network.ServiceStatusOnline)
	log.Infof("local service online name: %s, type: %s, id: %d", service.Name(), service.Type(), service.ID())
	return service, nil
}

// NewLocalService 新建本地服务
func (host *Host) NewLocalService(serviceType string, name string) (IService, error) {
	return host.NewService(serviceType, name, nil)
}

// Close 关闭集群服务器
func (host *Host) Close() {
	host.AllLocalServiceOffline()
	host.registryExit <- struct{}{}
	host.wg.Wait()
}

// AllLocalServiceOffline 关闭所有本地服务
func (host *Host) AllLocalServiceOffline() {
	for _, service := range host.localServices {
		host.localServiceEvents <- &eventServiceStatusChanged{
			service: service,
			status:  network.ServiceStatusOffline,
		}
	}
}
