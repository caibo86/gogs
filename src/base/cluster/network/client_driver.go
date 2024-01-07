// -------------------------------------------
// @file      : client_driver.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/3 下午7:12
// -------------------------------------------

package network

import (
	"fmt"
	"gogs/base/cberrors"
	"gogs/base/config"
	"hash/crc32"
	"runtime"
	"sync"
)

// ClientDriver 客户端驱动
type ClientDriver struct {
	sync.RWMutex
	remoteAddr            string                    // 远程地址
	userSessions          map[string]*ClientSession // 用户会话
	mutexGroup            []sync.Mutex              // 会话互斥锁列表
	sessionHandlerBuilder SessionHandlerBuilder     // 会话处理器构造器
	name                  string                    // 驱动名字
	protocol              ProtocolType              // 协议类型
}

// NewClientDriver 新建客户端驱动
func NewClientDriver(remoteAddr string, builder SessionHandlerBuilder, protocol ProtocolType) *ClientDriver {
	driver := &ClientDriver{
		remoteAddr:            remoteAddr,
		userSessions:          make(map[string]*ClientSession),
		mutexGroup:            make([]sync.Mutex, runtime.NumCPU()),
		sessionHandlerBuilder: builder,
		name:                  fmt.Sprintf("ClientDriver(%s)", remoteAddr),
		protocol:              protocol,
	}
	return driver
}

// Close 关闭驱动
func (driver *ClientDriver) Close() {

}

// String 获取客户端驱动名字 ClientDriver(remoteAddr)
func (driver *ClientDriver) String() string {
	return driver.name
}

// Name 获取驱动名字
func (driver *ClientDriver) Name() string {
	return driver.name
}

// Type 获取驱动类型
func (driver *ClientDriver) Type() DriverType {
	return DriverTypeClient
}

// Protocol 获取协议类型
func (driver *ClientDriver) Protocol() ProtocolType {
	return driver.protocol
}

// SetBuilder 设置会话处理器构造器
func (driver *ClientDriver) SetBuilder(builder SessionHandlerBuilder) {
	driver.sessionHandlerBuilder = builder
}

// lock 会话加锁回调
func (driver *ClientDriver) lock(session *ClientSession, callback func()) {
	var hashCode uint32
	name := session.name
	length := len(name)
	if length < 64 {
		scratch := make([]byte, 64)
		copy(scratch, name)
		hashCode = crc32.ChecksumIEEE(scratch[:length]) % uint32(len(driver.mutexGroup))
	} else {
		hashCode = crc32.ChecksumIEEE([]byte(name)) % uint32(len(driver.mutexGroup))
	}

	driver.mutexGroup[hashCode].Lock()
	defer driver.mutexGroup[hashCode].Unlock()
	callback()
}

// GetSession 获取指定名字的会话
func (driver *ClientDriver) GetSession(name string) (ISession, bool) {
	driver.RLock()
	defer driver.RUnlock()
	session, ok := driver.userSessions[name]
	return session, ok
}

// DelSession 删除指定的会话
func (driver *ClientDriver) DelSession(session ISession) {
	driver.Lock()
	defer driver.Unlock()
	if driver.userSessions[session.Name()] == session {
		delete(driver.userSessions, session.Name())
	}
}

// NewSession 创建一个客户端会话
func (driver *ClientDriver) NewSession(name string, ct ConnectionType) (ISession, error) {
	// 客户端仅支持外连
	if ct != ConnectionTypeOut {
		return nil, cberrors.New("client session only support out connection")
	}
	driver.RLock()
	if session, ok := driver.userSessions[name]; ok {
		driver.RUnlock()
		return session, cberrors.New("client session(%s) already exists", name)
	}
	driver.RUnlock()
	// 创建会话
	session := &ClientSession{
		driver:   driver,
		name:     name,
		status:   SessionStatusDisconnected,
		cached:   make(chan *Message, config.ClientSessionCache()),
		protocol: driver.protocol,
	}
	handler, err := driver.sessionHandlerBuilder(session)
	if err != nil {
		return nil, err
	}
	session.handler = handler
	// 异步连接会话
	go session.connect()
	driver.Lock()
	driver.userSessions[session.name] = session
	driver.Unlock()
	return session, nil
}
