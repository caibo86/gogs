package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"gogs/base/config"
	log "gogs/base/logger"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var global *Service
var once sync.Once

type EventHandler func(event *NodeEvent)

// Service etcd服务配置
type Service struct {
	sync.RWMutex
	Root           string                // 前缀
	Addr           []string              // etcd地址列表
	User           string                // etcd用户名
	Password       string                // etcd密码
	DepList        []string              // 依赖的服务列表
	TTL            int64                 // 租约时长second
	serviceType    string                // 本app服务类型
	serviceID      int64                 // 本app实例id
	serviceAddr    string                // 本app监听地址
	servicePort    string                // 本app端口
	closeSig       chan chan struct{}    // 关闭信号
	nodeUpdateChan chan NodeInfo         // 节点信息变动channel
	state          int32                 // 服务状态
	dependenceMap  map[string][]NodeInfo // 依赖服务列表
	client         *clientv3.Client      // etcd v3 client
	leaseID        clientv3.LeaseID      // 租约id
	callback       EventHandler          // key changed callback
}

// Check 检查依赖
func (s *Service) Check() bool {
	set := make(map[string]struct{})
	set[s.formatType()] = struct{}{} // 不依赖自身相同服务
	for _, sub := range s.DepList {
		// 配置去重
		if _, ok := set[sub]; ok {
			return false
		}
		set[sub] = struct{}{}
	}
	return s.serviceID != ServerIDDefault && s.serviceType != ServerTypeDefault
}

// Init 服务初始化
func Init(c *config.EtcdConfig, eventHandler EventHandler) error {
	var err error
	once.Do(func() {
		global = &Service{}
		global.Root = c.Root
		global.Addr = c.Addr
		global.User = c.User
		global.Password = c.Password
		global.DepList = c.DepList
		global.TTL = c.TTL
		global.callback = DefaultEtcdCallback
		if eventHandler != nil {
			global.callback = eventHandler
		}
		global.serviceID = c.ServiceID
		global.serviceType = c.ServiceType
		global.serviceAddr = c.ServiceAddr
		global.servicePort = c.ServicePort

		log.Debugf("etcd serviceType:%s serviceID:%d", global.serviceType, global.serviceID)
		if global.serviceAddr != "" {
			log.Debugf("etcd serviceAddr:%s serverPort:%s", global.serviceAddr, global.servicePort)
		}

		if !global.Check() {
			err = fmt.Errorf("[ETCD] etcd config is invalid")
		}
		global.closeSig = make(chan chan struct{})
		global.nodeUpdateChan = make(chan NodeInfo, 128)
		global.setState(StateDisconnect)
		watchChan, leaseChan, err := global.registrationAndDiscovery()
		if err != nil {
			err = fmt.Errorf("[ETCD] registrationAndDiscovery err:%s", err)
		}
		go global.loop(watchChan, leaseChan)
	})
	return err
}

// UpdateNodeWithExtra 注册本服务信息带额外信息
func (s *Service) UpdateNodeWithExtra(extra NodeInfo) {
	nodeInfo := s.GenNodeInfo()
	for k, v := range nodeInfo {
		extra[k] = v
	}
	select {
	case s.nodeUpdateChan <- extra:
		log.Debugf("[ETCD] UpdateNodeWithExtra info:%v", extra)
	default:
		log.Errorf("[ETCD] nodeUpdateChan is full, info:%v", extra)
	}
}

// UpdateNodeWithExtra 注册本服务信息带额外信息
func UpdateNodeWithExtra(extra NodeInfo) {
	global.UpdateNodeWithExtra(extra)
}

// UpdateOtherNode 注册一个节点信息
func (s *Service) UpdateOtherNode(node NodeInfo) {
	select {
	case s.nodeUpdateChan <- node:
		log.Debugf("[ETCD] UpdateOtherNode info:%v", node)
	default:
		log.Errorf("[ETCD] nodeUpdateChan is full, info:%v", node)
	}
}

// UpdateOtherNode 注册一个节点信息
func UpdateOtherNode(node NodeInfo) {
	global.UpdateOtherNode(node)
}

// SyncUpdateNodeInfo 同步注册信息 注意线程
func (s *Service) SyncUpdateNodeInfo(node NodeInfo) error {
	return s.updateNodeInfo(node)
}

// GetDepListByType 获取指定类型的依赖服务列表
func (s *Service) GetDepListByType(t string) ([]NodeInfo, error) {
	s.RLock()
	defer s.RUnlock()

	if s.dependenceMap == nil {
		return nil, fmt.Errorf("[ETCD] service dependence nil")
	}
	deps, ok := s.dependenceMap[t]
	if !ok {
		return nil, fmt.Errorf("[ETCD] service dependence not found. type:%s", t)
	}
	ret := make([]NodeInfo, len(deps))
	for i, node := range deps {
		ret[i] = node.Clone()
	}
	return ret, nil
}

func GetDepListByType(t string) ([]NodeInfo, error) {
	return global.GetDepListByType(t)
}

// GetDepByTypeAndID 获取指定类型和ID的依赖服务
func (s *Service) GetDepByTypeAndID(t string, id int64) (NodeInfo, error) {
	s.RLock()
	defer s.RUnlock()
	if s.dependenceMap == nil {
		return nil, fmt.Errorf("[ETCD] service dependence nil")
	}
	deps, ok := s.dependenceMap[t]
	if !ok {
		return nil, fmt.Errorf("[ETCD] service dependence not found. type:%s", t)
	}
	for _, node := range deps {
		if node.GetID() == id {
			return node.Clone(), nil
		}
	}
	return nil, fmt.Errorf("[ETCD] service dependence not found. type:%s id:%d", t, id)
}

func GetDepByTypeAndID(t string, id int64) (NodeInfo, error) {
	return global.GetDepByTypeAndID(t, id)
}

// Exit 关闭服务
func (s *Service) Exit() {
	finish := make(chan struct{}, 1)
	s.closeSig <- finish
	// 此处同步等待
	<-finish
	log.Infof("[ETCD] exit done")
}

func Exit() {
	global.Exit()
}

// loop etcd监听线程
func (s *Service) loop(watchChan clientv3.WatchChan, leaseChan <-chan *clientv3.LeaseKeepAliveResponse) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("[ETCD] panic recover:%s", err, zap.Stack("stack"))
		}
	}()
	reconnectTimer := time.NewTimer(ReconnectDuration)
	setEtcdDisconnect := func() {
		watchChan = nil
		leaseChan = nil
		s.setState(StateDisconnect)
	}
	for {
		select {
		case watchRsp, ok := <-watchChan:
			if !ok {
				log.Error("[ETCD] watchChan false")
				setEtcdDisconnect()
				continue
			}
			for _, event := range watchRsp.Events {
				switch event.Type {
				case clientv3.EventTypePut:
					s.addDependence(string(event.Kv.Key), string(event.Kv.Value))
				case clientv3.EventTypeDelete:
					s.delDependence(string(event.Kv.Key))
				}
			}
		case leaseRsp, ok := <-leaseChan:
			if leaseRsp == nil || !ok {
				log.Error("[ETCD] leaseRsp false")
				setEtcdDisconnect()
			}
		case <-reconnectTimer.C:
			reconnectTimer.Reset(ReconnectDuration)
			if !s.needReconnect() {
				continue
			}
			log.Error("[ETCD] begin reconnect")
			s.close()
			newWatchChan, newLeaseChan, err := s.registrationAndDiscovery()
			if err != nil {
				log.Errorf("[ETCD] reconnect err:%s", err)
				continue
			}
			watchChan = newWatchChan
			leaseChan = newLeaseChan
		case info, ok := <-s.nodeUpdateChan:
			if !ok {
				log.Error("[ETCD] nodeUpdateChan closed")
				continue
			}
			log.Debugf("[ETCD] start update nodeInfo:%d", info.GetID())
			if err := s.updateNodeInfo(info); err != nil {
				log.Errorf("[ETCD] updateNodeInfo err:%s", err)
				continue
			}
			log.Debugf("[ETCD] finish update nodeInfo:%d", info.GetID())
		case finish := <-s.closeSig:
			if finish == nil {
				log.Warnf("[ETCD] closeSig get a nil chan")
				continue
			}
			s.close()
			s.setState(StateClosed)
			finish <- struct{}{}
			return
		}
	}
}

// formatType 当前服务的类型key 不包含实例id
func (s *Service) formatType() string {
	return fmt.Sprintf("%s/%s", s.Root, s.serviceType)
}

// setState 设置当前服务注册状态
func (s *Service) setState(state int32) {
	atomic.StoreInt32(&s.state, state)
}

func (s *Service) needReconnect() bool {
	return atomic.LoadInt32(&s.state) == StateDisconnect
}

// getDepTypeByKey 从依赖服务的key中获取其类型
func (s *Service) getDepTypeByKey(key string) (string, error) {
	arr := strings.Split(key, "/")
	if len(arr) < 3 {
		log.Errorf("[ETCD] invalid key: %s", key)
		return "", fmt.Errorf("[ETCD] invalid key: %s", key)
	}
	return arr[1], nil
}

// updateNodeInfo 向etcd更新节点信息
func (s *Service) updateNodeInfo(nodeInfo NodeInfo) error {
	if !(atomic.LoadInt32(&s.state) == StateConnected) {
		return fmt.Errorf("[ETCD] etcd is not connected")
	}
	key := s.formatKey(nodeInfo.GetType(), nodeInfo.GetID())
	infoJson, err := json.Marshal(nodeInfo)
	if err != nil {
		return fmt.Errorf("[ETCD] node json.Marshal err:%s", err)
	}
	value := string(infoJson)
	ctx, cancel := context.WithTimeout(context.TODO(), Timeout)
	defer cancel()
	if nodeInfo.GetType() == s.serviceType && nodeInfo.GetID() == s.serviceID {
		log.Debugf("[ETCD] Put(key:%s value:%s leaseID:%d)", key, value, s.leaseID)
		_, err = s.client.Put(ctx, key, value, clientv3.WithLease(s.leaseID))
	} else {
		log.Debugf("[ETCD] Put(key:%s value:%s)", key, value)
		_, err = s.client.Put(ctx, key, value)
	}
	if err != nil {
		return fmt.Errorf("[ETCD] Put err:%s", err)
	}
	return nil
}

// clearDependence 清空所有依赖服务的配置
func (s *Service) clearDependence() {
	s.Lock()
	defer s.Unlock()
	s.dependenceMap = make(map[string][]NodeInfo)
	for _, service := range s.DepList {
		s.dependenceMap[service] = make([]NodeInfo, 0)
	}
}

// delDependence 删除依赖服务的信息
func (s *Service) delDependence(key string) {
	s.Lock()
	defer s.Unlock()
	serviceAddr := strings.Split(key, "/")
	if len(serviceAddr) < 3 {
		log.Errorf("[ETCD] invalid key:%s", key)
		return
	}
	depType := serviceAddr[1]
	for index, info := range s.dependenceMap[depType] {
		if info == nil {
			continue
		}
		if s.formatKey(info.GetType(), info.GetID()) == key {
			if s.callback != nil {
				s.callback(&NodeEvent{
					Node:  info,
					Event: EventDelete,
				})
			}
			s.dependenceMap[depType] = append(s.dependenceMap[depType][:index], s.dependenceMap[depType][index+1:]...)
			break
		}
	}
}

// addDependence 添加依赖服务的信息
func (s *Service) addDependence(key, value string) {
	s.Lock()
	defer s.Unlock()
	var newNodeInfo NodeInfo
	depType, err := s.getDepTypeByKey(key)
	if err != nil {
		return
	}
	newNodeInfo = NodeDecode(value)
	if newNodeInfo == nil || newNodeInfo.IsInvalid() {
		log.Warnf("[ETCD] node is invalid: %v", newNodeInfo)
		return
	}

	// 检查是否依赖该服务
	if _, ok := s.dependenceMap[depType]; !ok {
		return
	}

	for idx, node := range s.dependenceMap[depType] {
		if node == nil {
			continue
		}
		if node.GetID() != newNodeInfo.GetID() {
			continue
		}
		if IsNodeEqual(node, newNodeInfo) {
			return
		}
		s.dependenceMap[depType][idx] = newNodeInfo
		if s.callback != nil {
			s.callback(&NodeEvent{
				Node:  newNodeInfo,
				Event: EventUpdate,
			})
		}
		return
	}
	s.dependenceMap[depType] = append(s.dependenceMap[depType], newNodeInfo)
	if s.callback != nil {
		s.callback(&NodeEvent{
			Node:  newNodeInfo,
			Event: EventAdd,
		})
	}
}

// preKey 同组服务前缀+"/"
func (s *Service) preKey() string {
	return s.Root + "/"
}

// formatKey 组装服务用于etcd的key
func (s *Service) formatKey(serviceType string, serviceID int64) string {
	return fmt.Sprintf("%s/%s/%d", s.Root, serviceType, serviceID)
}

// fetchRemoteCfg 拉取etcd中的配置信息
func (s *Service) fetchRemoteCfg() error {
	ctx, cancel := context.WithTimeout(context.TODO(), Timeout)
	defer cancel()
	// 取所有前缀一致的key
	rsp, err := s.client.Get(ctx, s.preKey(), clientv3.WithPrefix())
	if err != nil {
		return fmt.Errorf("[ETCD] fetch deps err: %s", err)
	}
	if rsp == nil {
		return fmt.Errorf("[ETCD] fetch deps ret nil")
	}
	for _, ev := range rsp.Kvs {
		key := string(ev.Key)
		value := string(ev.Value)
		// 根据DepList取需要的节点信息存储
		s.addDependence(key, value)
	}
	return nil
}

// GenNodeInfo 服务信息生成map[string]any
func (s *Service) GenNodeInfo() NodeInfo {
	nodeInfo := NodeInfo{
		NodeInfoKeyType: s.serviceType,
		NodeInfoKeyID:   s.serviceID,
		NodeInfoKeyPID:  os.Getpid(),
	}
	if hostname, err := os.Hostname(); err == nil {
		nodeInfo[NodeInfoKeyHostname] = hostname
	}
	if s.serviceAddr != "" || s.servicePort != "" {
		addr := s.serviceAddr
		if s.servicePort != "" {
			addr = net.JoinHostPort(addr, s.servicePort)
		}
		nodeInfo[NodeInfoKeyAddrList] = []string{addr}
	}
	return nodeInfo
}

// register 注册服务节点到etcd 返回租约保活channel
func (s *Service) register(leaseCtx context.Context) (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), Timeout)
	defer cancel()
	key := s.formatKey(s.serviceType, s.serviceID)
	// 检查是否有同名节点注册
	rsp, err := s.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if err == nil && len(rsp.Kvs) != 0 {
		return nil, fmt.Errorf("[ETCD] node %v already exist", key)
	}
	// 序列化节点信息为Json串
	infoJson, err := json.Marshal(s.GenNodeInfo())
	if err != nil {
		return nil, fmt.Errorf("[ETCD] node json.Marshal err:%s", err)
	}
	value := string(infoJson)
	ctx, cancel = context.WithTimeout(context.TODO(), Timeout)
	defer cancel()
	// 申请租约
	lease, err := s.client.Grant(ctx, s.TTL)
	if err != nil {
		return nil, fmt.Errorf("[ETCD] grant lease err:%s", err)
	}
	s.leaseID = lease.ID
	ctx, cancel = context.WithTimeout(context.TODO(), Timeout)
	defer cancel()
	log.Debugf("[ETCD] Put(key:%s value:%s leaseID:%d)", key, value, s.leaseID)
	// 注册
	_, err = s.client.Put(ctx, key, value, clientv3.WithLease(s.leaseID))
	if err != nil {
		return nil, fmt.Errorf("[ETCD] Put err:%s", err)
	}
	return s.client.KeepAlive(leaseCtx, s.leaseID)
}

// registrationAndDiscovery 注册服务并获取监听channel
func (s *Service) registrationAndDiscovery() (clientv3.WatchChan, <-chan *clientv3.LeaseKeepAliveResponse, error) {
	cfg := clientv3.Config{
		Endpoints:   s.Addr,
		DialTimeout: Timeout,
		Username:    s.User,
		Password:    s.Password,
	}
	client, err := clientv3.New(cfg)
	if err != nil {
		log.Errorf("[ETCD] clientv3.New err: %s", err)
		return nil, nil, err
	}
	if s.TTL <= 0 {
		s.TTL = DefaultTTL
	}
	s.clearDependence()
	s.client = client
	err = s.fetchRemoteCfg()
	if err != nil {
		return nil, nil, err
	}
	leaseChan, err := s.register(context.Background())
	if err != nil {
		return nil, nil, err
	}
	watchChan := s.client.Watch(clientv3.WithRequireLeader(context.Background()), s.preKey(), clientv3.WithPrefix())
	err = s.fetchRemoteCfg()
	if err != nil {
		return nil, nil, err
	}
	s.setState(StateConnected)
	return watchChan, leaseChan, nil
}

// close 取消租约 关闭etcd连接
func (s *Service) close() {
	if s.client == nil {
		return
	}
	defer func() {
		if err := s.client.Close(); err != nil {
			log.Errorf("[ETCD] clientv3 client Close err:%s", err)
		}
	}()
	log.Debugf("[ETCD] close. leaseID:%d", s.leaseID)
	if s.leaseID > 0 {
		ctx, cancel := context.WithTimeout(context.TODO(), Timeout)
		defer cancel()
		_, err := s.client.Revoke(ctx, s.leaseID)
		if err != nil {
			log.Errorf("[ETCD] lease revoke err:%s", err)
		}
		s.leaseID = 0
	}
}

func (s *Service) GetServiceType() string {
	return s.serviceType
}

func (s *Service) GetServiceID() int64 {
	return s.serviceID
}

//func (s *Service) SetServiceType(t string) {
//	s.serviceType = t
//}
//
//func (s *Service) SetServiceID(id int64) {
//	s.serviceID = id
//}

func (s *Service) SetServiceAddr(addr string) {
	s.serviceAddr = addr
}

func (s *Service) SetServicePort(port string) {
	s.servicePort = port
}

func (s *Service) SetServiceCallback(callback EventHandler) {
	s.callback = callback
}

func SetServiceCallback(callback EventHandler) {
	global.SetServiceCallback(callback)
}

// DefaultEtcdCallback 默认的etcd事件监听回调
func DefaultEtcdCallback(nodeEvent *NodeEvent) {
	log.Warnw("[ETCD] ignored etcd event",
		zap.String("event", nodeEvent.Event.String()),
		zap.Any("nodeInfo", nodeEvent.Node))
}
