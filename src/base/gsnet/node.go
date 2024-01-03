// -------------------------------------------
// @file      : node.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/29 下午6:52
// -------------------------------------------

package gsnet

import (
	"fmt"
	"gogs/base/gserrors"
	log "gogs/base/logger"
	"sync"
)

// DriverType 驱动类型
type DriverType int32

const (
	DriverTypeCluster DriverType = 1
	DriverTypeClient  DriverType = 2
	DriverTypeGate    DriverType = 3
)

// ProtocolType 协议类型
type ProtocolType int32

const (
	ProtocolTCP       ProtocolType = 1
	ProtocolWebsocket ProtocolType = 2
)

// ISession 会话接口
type ISession interface {
	fmt.Stringer
	Write(*Message) error
	Status() SessionStatus    // 状态
	DriverType() DriverType   // 驱动类型
	Close()                   // 关闭会话
	Handler() ISessionHandler // 会话句柄
	Name() string             // 会话标识符
	RemoteAddr() string       // 远程地址
}

// ISessionHandler 会话处理器
type ISessionHandler interface {
	Read(ISession, *Message)     // 从会话读取一个消息,实现请注意线程安全
	StatusChanged(SessionStatus) // 通知一个新的会话状态
}

// SessionHandlerBuilder 会话句柄建造者,指定会话生成会话句柄
type SessionHandlerBuilder func(ISession) (ISessionHandler, error)

// IDriver 传输层协议提供者 驱动 一种Driver对应于一个协议名
type IDriver interface {
	fmt.Stringer
	Type() DriverType                                    // 驱动类型
	GetSession(string) (ISession, bool)                  // 通过名字获取会话
	NewSession(string, ConnectionType) (ISession, error) // 创建会话
	DelSession(ISession)                                 // 删除指定的会话
	SetBuilder(SessionHandlerBuilder)                    // 设置会话句柄建造者
	Close()                                              // 关闭驱动
}

// Node 集群节点
type Node struct {
	sync.RWMutex
	drivers map[DriverType]IDriver // 按协议索引的驱动集合
}

// NewNode 新建集群节点
func NewNode() *Node {
	return &Node{
		drivers: make(map[DriverType]IDriver),
	}
}

// Close 关闭节点
func (node *Node) Close() {
	for _, driver := range node.drivers {
		driver.Close()
	}
}

// NewDriver 将驱动注册到节点
func (node *Node) NewDriver(driver IDriver) error {
	node.Lock()
	defer node.Unlock()
	driverType := driver.Type()
	if _, ok := node.drivers[driverType]; ok {
		// 同个节点下单个类型的驱动只能有一个
		return gserrors.Newf("duplicate driver type support: %s", driverType)
	}
	node.drivers[driverType] = driver
	return nil
}

// GetSession 获取指定类型驱动中指定名字的会话
func (node *Node) GetSession(driverType DriverType, name string) (ISession, bool) {
	node.RLock()
	defer node.RUnlock()
	if driver, ok := node.drivers[driverType]; ok {
		return driver.GetSession(name)
	}
	log.Warnf("get channel failed: driver type not found: %s", driverType)
	return nil, false
}

// NewSession 在指定类型驱动上新建指定名字的会话
func (node *Node) NewSession(driverType DriverType, name string, ct ConnectionType) (ISession, error) {
	node.RLock()
	defer node.RUnlock()
	if driver, ok := node.drivers[driverType]; ok {
		return driver.NewSession(name, ct)
	}
	return nil, gserrors.Newf("new channel failed: driver type not found: %s", driverType)
}

// DelSession 在指定类型驱动上删除指定名字的会话
func (node *Node) DelSession(channel ISession) {
	node.RLock()
	defer node.RUnlock()
	if driver, ok := node.drivers[channel.DriverType()]; ok {
		driver.DelSession(channel)
	}
}
