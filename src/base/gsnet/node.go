// -------------------------------------------
// @file      : node.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2023/12/29 下午6:52
// -------------------------------------------

package gsnet

import (
	"errors"
	"fmt"
	"gogs/base/gserrors"
	log "gogs/base/logger"
	"sync"
)

var (
	ErrNode = errors.New("cluster node error") // 集群节点错误
)

// IChannel 通道接口
type IChannel interface {
	fmt.Stringer
	Write(*Message) error
	Status() Status         // 状态
	Protocol() string       // 协议
	Close()                 // 关闭通道
	Handle() IChannelHandle // 通道句柄
	Name() string           // 通道标识符
}

// IChannelHandle 通道句柄
type IChannelHandle interface {
	Read(IChannel, *Message) // 从通道读取一个消息,实现请注意线程安全
	StatusChanged(Status)    // 通知一个新的通道状态
}

// ChannelHandleBuilder 通道句柄建造者,指定通道生成通道句柄
type ChannelHandleBuilder func(IChannel) (IChannelHandle, error)

// IDriver 传输层协议提供者 驱动 一种Driver对应于一个协议名
type IDriver interface {
	fmt.Stringer
	Protocol() string                          // 驱动协议名
	GetChannel(string) (IChannel, bool)        // 通过名字获取通道
	NewChannel(string, byte) (IChannel, error) // 创建通道
	DelChannel(IChannel)                       // 删除指定的通道
	SetBuilder(ChannelHandleBuilder)           // 设置通道句柄建造者
	Close()                                    // 关闭驱动
}

// Node 集群节点
type Node struct {
	sync.RWMutex
	drivers map[string]IDriver // 按协议索引的驱动集合
}

// NewNode 新建集群节点
func NewNode() *Node {
	return &Node{
		drivers: make(map[string]IDriver),
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
	protocol := driver.Protocol()
	if _, ok := node.drivers[protocol]; ok {
		// 同个节点下不能有相同协议的驱动
		return gserrors.Newf("duplicate driver support: %s", protocol)
	}
	node.drivers[protocol] = driver
	return nil
}

// GetChannel 在节点注册的指定协议驱动中查找指定名字的通道
func (node *Node) GetChannel(protocol string, name string) (IChannel, bool) {
	node.RLock()
	defer node.RUnlock()
	if driver, ok := node.drivers[protocol]; ok {
		return driver.GetChannel(name)
	}
	log.Warnf("get channel failed: driver protocol not found: %s", protocol)
	return nil, false
}

// NewChannel 在节点内指定协议的驱动上新建指定名字的通道
func (node *Node) NewChannel(protocol string, name string, connectionType byte) (IChannel, error) {
	node.RLock()
	defer node.RUnlock()
	if driver, ok := node.drivers[protocol]; ok {
		return driver.NewChannel(name, connectionType)
	}
	return nil, gserrors.Newf("new channel failed: driver protocol not found: %s", protocol)
}

// DelChannel 在节点内指定协议的驱动上删除指定名字的通道
func (node *Node) DelChannel(channel IChannel) {
	node.RLock()
	defer node.RUnlock()
	if driver, ok := node.drivers[channel.Protocol()]; ok {
		driver.DelChannel(channel)
	}
}
