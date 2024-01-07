// -------------------------------------------
// @file      : service_status.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/3 下午11:40
// -------------------------------------------

package cluster

import "gogs/base/cluster/network"

// ServiceListener 服务状态监听器,返回false则停止监听
type ServiceListener func(service IService, status network.ServiceStatus) bool

// eventListenService 监听服务事件(通过服务名)
type eventListenService struct {
	name     string          // 监听的服务名称
	listener ServiceListener // 监听器
}

// eventListenServiceType 监听服务事件(通过服务类型)
type eventListenServiceType struct {
	typename string          // 监听的服务类型
	listener ServiceListener // 监听器
}

// eventServiceStatusChanged 监听服务事件(通过服务状态)
type eventServiceStatusChanged struct {
	service IService              // 目标服务
	status  network.ServiceStatus // 需要监听的状态
}

// ServiceStatusPublisher 服务状态发布器
type ServiceStatusPublisher struct {
	serviceListeners     map[string][]ServiceListener // 按服务名字分类的监听器列表
	serviceTypeListeners map[string][]ServiceListener // 按服务类型分类的监听器列表
	services             map[string]IService          // 已上线的服务
	events               chan interface{}             // 事件通道
}

// NewServiceStatusPublisher 新建服务状态发布器
func NewServiceStatusPublisher() *ServiceStatusPublisher {
	publisher := &ServiceStatusPublisher{
		serviceListeners:     make(map[string][]ServiceListener),
		serviceTypeListeners: make(map[string][]ServiceListener),
		services:             make(map[string]IService),
		events:               make(chan interface{}, 1024),
	}
	// 启动发布器
	go publisher.handleEvent()
	return publisher
}

// handleEvent 事件处理
func (publisher *ServiceStatusPublisher) handleEvent() {
	var mark bool
	// 从事件接收通道接收事件
	for event := range publisher.events {
		switch event.(type) {
		case *eventListenService:
			mark = true
			// 注册服务监听事件
			e := event.(*eventListenService)
			for name, service := range publisher.services {
				if e.name == name {
					if !e.listener(service, network.ServiceStatusOnline) {
						mark = false
					}
				}
			}
			if mark {
				publisher.serviceListeners[e.name] = append(publisher.serviceListeners[e.name], e.listener)
			}
		case *eventListenServiceType:
			mark = true
			e := event.(*eventListenServiceType)
			for _, service := range publisher.services {
				if e.typename == service.Type() {
					if !e.listener(service, network.ServiceStatusOnline) {
						mark = false
					}
				}
			}
			if mark {
				publisher.serviceTypeListeners[e.typename] = append(publisher.serviceTypeListeners[e.typename], e.listener)
			}
		case *eventServiceStatusChanged:
			e := event.(*eventServiceStatusChanged)
			if e.status == network.ServiceStatusOnline {
				publisher.services[e.service.Name()] = e.service
			} else {
				delete(publisher.services, e.service.Name())
			}
			listeners := publisher.serviceListeners[e.service.Name()]
			var tmp []ServiceListener
			for _, listener := range listeners {
				if listener(e.service, e.status) {
					tmp = append(tmp, listener)
				}
			}
			publisher.serviceListeners[e.service.Name()] = tmp
			listeners = publisher.serviceTypeListeners[e.service.Type()]
			tmp = nil
			for _, listener := range listeners {
				if listener(e.service, e.status) {
					tmp = append(tmp, listener)
				}
			}
			publisher.serviceTypeListeners[e.service.Type()] = tmp
		}
	}
}

// ServiceStatusChanged 服务状态变更
func (publisher *ServiceStatusPublisher) ServiceStatusChanged(service IService, status network.ServiceStatus) {
	publisher.events <- &eventServiceStatusChanged{
		service: service,
		status:  status,
	}
}

// ListenService 监听服务事件(通过服务名)
func (publisher *ServiceStatusPublisher) ListenService(name string, listener ServiceListener) {
	publisher.events <- &eventListenService{
		name:     name,
		listener: listener,
	}
}

// ListenServiceType 监听服务事件(通过服务类型)
func (publisher *ServiceStatusPublisher) ListenServiceType(typename string, listener ServiceListener) {
	publisher.events <- &eventListenServiceType{
		typename: typename,
		listener: listener,
	}
}
