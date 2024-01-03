// -------------------------------------------
// @file      : gsdocker.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/29 下午6:28
// -------------------------------------------

package gsdock

import (
	"gogs/base/config"
	"gogs/base/gserrors"
	"gogs/base/gsnet"
	log "gogs/base/logger"
	"sync"
	"sync/atomic"
	"time"
)

// GSDock 集群服务节点
type GSDock struct {
	*RPC                                                    // 远程调用管理器
	*ServiceStatusPublisher                                 // 服务状态发布器
	wg                      sync.WaitGroup                  // WaitGroup
	idgen                   uint32                          // 服务ID生成器
	Node                    *gsnet.Node                     // 网络节点
	rMutex                  sync.RWMutex                    // 远程服务 读写锁
	lMutex                  sync.RWMutex                    // 本地服务 读写锁
	tMutex                  sync.RWMutex                    // 建造者集合 读写锁
	lServices               map[ID]IService                 // 本地服务集合,通过ID索引
	rServices               map[string]*Neighbor            // 集群的邻居集合,key为其对应Session的name
	builders                map[string]ITypeBuilder         // 服务建造者集合
	lServiceEvents          chan *eventServiceStatusChanged // 本地服务状态变更事件通道
	cerExit                 chan struct{}                   // 退出通道
}

// NewGSDock 新建集群服务节点
func NewGSDock(localAddr string) *GSDock {
	dock := &GSDock{
		RPC:                    NewRPC(),
		ServiceStatusPublisher: NewServiceStatusPublisher(),
		Node:                   gsnet.NewNode(),
		lServices:              make(map[ID]IService),
		rServices:              make(map[string]*Neighbor),
		builders:               make(map[string]ITypeBuilder),
		lServiceEvents:         make(chan *eventServiceStatusChanged),
		cerExit:                make(chan struct{}, 1),
	}
	_ = dock.Node.NewDriver(
		gsnet.NewClusterDriver(
			localAddr,
			func(session gsnet.ISession) (gsnet.ISessionHandler, error) {
				return NewClusterRemote(dock, session), nil
			}))
	go func() {
		dock.wg.Add(1)
		ticker := time.Tick(config.ClusterCerInterval())
		for _ = range ticker {
			if !dock.tryCERs() {
				break
			}
		}
		log.Debug("GSDock cer exit")
		dock.wg.Done()
	}()
	return dock
}

// tryCERs
func (dock *GSDock) tryCERs() bool {
	events := make(map[ID]*eventServiceStatusChanged)
	var content []*gsnet.CER
	for {
		if len(events) > config.ClusterCerMax() {
			break
		}
		select {
		case event, ok := <-dock.lServiceEvents:
			if !ok {
				return false
			}
			events[event.service.ID()] = event
		case <-dock.cerExit:
			if len(events) > 0 {
				dock.cerExit <- struct{}{}
				goto NEXT
			}
			return false
		default:
			if len(events) == 0 {
				return true
			}
			goto NEXT
		}
	}
NEXT:
	for _, event := range events {
		cer := &gsnet.CER{
			Name: event.service.Name(),
			ID:   uint32(event.service.ID()),
			Type: event.service.Type(),
		}
		if event.status == gsnet.ServiceStatusOnline {
			cer.Add = true
		} else {
			cer.Add = false
		}
		content = append(content, cer)
	}
	cers := &gsnet.CERs{
		Data: content,
	}
	data := cers.Marshal()
	msg := &gsnet.Message{
		Type: gsnet.MessageTypeCER,
		Data: data,
	}
	dock.rMutex.RLock()
	defer dock.rMutex.RUnlock()
	for _, neighbor := range dock.rServices {
		err := neighbor.remote.Write(msg)
		if err != nil {
			log.Errorf("neighbor remote write msg err: %s", err)
		}
	}
	log.Info("GSDock cer sent")
	return true
}

// newID 生成ServiceID
func (dock *GSDock) newID() ID {
	for {
		val := atomic.AddUint32(&dock.idgen, 1)
		if _, ok := dock.lServices[ID(val)]; !ok {
			return ID(val)
		}
	}
}

// handleCERs 处理来自邻居节点的服务状态变更通知
func (dock *GSDock) handleCERs(remote *ClusterRemote, cers *gsnet.CERs) {
	dock.rMutex.Lock()
	defer dock.rMutex.Unlock()
	if neighbor, ok := dock.rServices[remote.Name()]; ok {
		for _, cer := range cers.Data {
			if cer.Add {
				if service, ok := neighbor.services[cer.Name]; ok {
					if service.RemoteID() != ID(cer.ID) || service.Type() != cer.Type {
						delete(neighbor.services, cer.Name)
						delete(neighbor.servicesByID, service.RemoteID())
						dock.ServiceStatusChanged(service, gsnet.ServiceStatusOffline)
					} else {
						continue
					}
				}
				dock.tMutex.RLock()
				builder, ok := dock.builders[cer.Type]
				dock.tMutex.RUnlock()
				if ok {
					remoteService := builder.NewRemoteService(remote, cer.Name, dock.newID(), ID(cer.ID), nil)
					neighbor.services[cer.Name] = remoteService
					neighbor.servicesByID[remoteService.RemoteID()] = remoteService
					dock.ServiceStatusChanged(remoteService, gsnet.ServiceStatusOnline)
					continue
				}
			} else {
				if service, ok := neighbor.services[cer.Name]; ok {
					if service.RemoteID() == ID(cer.ID) && service.Type() == cer.Type {
						delete(neighbor.services, cer.Name)
						delete(neighbor.servicesByID, service.RemoteID())
						dock.ServiceStatusChanged(service, gsnet.ServiceStatusOffline)
					}
				}
			}
		}
	}
}

// statusChanged 会话状态变更
func (dock *GSDock) statusChanged(remote *ClusterRemote, status gsnet.SessionStatus) {
	switch status {
	case gsnet.SessionStatusInConnected, gsnet.SessionStatusOutConnected:
		dock.rMutex.Lock()
		dock.rServices[remote.Name()] = NewNeighbor(remote)
		dock.rMutex.Unlock()
		var content []*gsnet.CER
		dock.lMutex.RLock()
		for _, service := range dock.lServices {
			cer := &gsnet.CER{
				Name: service.Name(),
				ID:   uint32(service.ID()),
				Type: service.Type(),
				Add:  true,
			}
			content = append(content, cer)
		}
		dock.lMutex.RUnlock()
		cers := &gsnet.CERs{
			Data: content,
		}
		data := cers.Marshal()
		msg := &gsnet.Message{
			Type: gsnet.MessageTypeCER,
			Data: data,
		}
		err := remote.Write(msg)
		if err != nil {
			log.Errorf("remote write msg err: %s", err)
		}
	case gsnet.SessionStatusClosed, gsnet.SessionStatusDisconnected:
		dock.rMutex.Lock()
		delete(dock.rServices, remote.Name())
		dock.rMutex.Unlock()
	}
}

// handleCall 处理来自对本地服务的调用
func (dock *GSDock) handleCall(call *gsnet.Call) (*gsnet.Return, error) {
	dock.lMutex.RLock()
	defer dock.lMutex.RUnlock()
	if service, ok := dock.lServices[ID(call.ServiceID)]; ok {
		return service.Call(call)
	}
	return nil, gserrors.Newf("local service not found: %d", call.ServiceID)
}

// Register 注册服务建造者
func (dock *GSDock) Register(builder ITypeBuilder) (ITypeBuilder, error) {
	dock.tMutex.Lock()
	defer dock.tMutex.Unlock()
	// TODO 把建造者的String改为ServiceType
	if _, ok := dock.builders[builder.String()]; ok {
		return nil, gserrors.Newf("duplicate builder: %s", builder)
	}
	dock.builders[builder.String()] = builder
	return builder, nil
}
