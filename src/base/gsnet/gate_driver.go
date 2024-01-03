// -------------------------------------------
// @file      : gate_driver.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/3 上午10:02
// -------------------------------------------

package gsnet

import (
	"fmt"
	"github.com/gorilla/websocket"
	"gogs/base/config"
	"gogs/base/gserrors"
	log "gogs/base/logger"
	"net"
	"net/http"
	"runtime"
	"sync"
	"time"
)

// upgrader websocket默认参数
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// GateDriver 网关驱动
type GateDriver struct {
	sync.RWMutex
	localAddr             string                  // 本地地址
	remotes               map[string]*GateSession // 远程会话
	mutexGroup            []sync.Mutex            // 会话互斥锁列表
	sessionHandlerBuilder SessionHandlerBuilder   // 会话处理器建造者
	name                  string                  // 驱动名字
	protocol              ProtocolType            // 协议类型
}

// NewGateDriver 新建网关驱动
func NewGateDriver(localAddr string, builder SessionHandlerBuilder, protocol ProtocolType) *GateDriver {
	driver := &GateDriver{
		localAddr:             localAddr,
		remotes:               make(map[string]*GateSession),
		mutexGroup:            make([]sync.Mutex, runtime.NumCPU()),
		sessionHandlerBuilder: builder,
		name:                  fmt.Sprintf("GateDriver(%s)", localAddr),
		protocol:              protocol,
	}
	go driver.run()
	return driver
}

// Close 关闭驱动
func (driver *GateDriver) Close() {
}

// String 获取网关驱动名字 GateDriver(localAddr)
func (driver *GateDriver) String() string {
	return driver.name
}

// Name 获取网关驱动名字 GateDriver(localAddr)
func (driver *GateDriver) Name() string {
	return driver.name
}

// SetBuilder 设置会话处理器建造者
func (driver *GateDriver) SetBuilder(builder SessionHandlerBuilder) {
	driver.sessionHandlerBuilder = builder
}

// Type 获取驱动类型
func (driver *GateDriver) Type() DriverType {
	return DriverTypeGate
}

// Protocol 协议类型
func (driver *GateDriver) Protocol() ProtocolType {
	return driver.protocol
}

// GetSession 获取指定名字的会话
func (driver *GateDriver) GetSession(name string) (ISession, bool) {
	driver.RLock()
	defer driver.RUnlock()
	channel, ok := driver.remotes[name]
	return channel, ok
}

// NewSession 仅实现IDriver接口,GateDriver不支持手动创建会话
func (driver *GateDriver) NewSession(name string, connectionType byte) (ISession, error) {
	return nil, gserrors.Newf("gate driver: %s not support manual new channel", driver)
}

// DelSession 删除指定的会话
func (driver *GateDriver) DelSession(channel ISession) {
	driver.Lock()
	defer driver.Unlock()
	if driver.remotes[channel.Name()] == channel {
		delete(driver.remotes, channel.Name())
	}
}

// run 启动网关驱动
func (driver *GateDriver) run() {
	switch driver.protocol {
	case ProtocolTCP:
		driver.runTCP()
	case ProtocolWebsocket:
		http.HandleFunc("/gate", driver.handleAcceptWebsocket)
		driver.runWebsocket()
	default:
		gserrors.Panicf("unsupported protocol: %s", driver.protocol)
	}
}

// runTCP 以tcp协议启动网关驱动
func (driver *GateDriver) runTCP() {
	// 启动监听
	log.Infof("start tcp listen: %s", driver.localAddr)
	listener, err := net.Listen("tcp", driver.localAddr)
	if err != nil {
		log.Errorf("tcp listen: %s err: %s", driver.localAddr, err)
		time.AfterFunc(config.ListenRetryInterval(), driver.run)
		return
	}
	for {
		var conn net.Conn
		conn, err = listener.Accept()
		if err != nil {
			log.Errorf("tcp accept: %s err: %s", driver.localAddr, err)
			continue
		}
		stream := NewStream(conn, conn)
		go driver.handleAccept(stream, conn, nil)
	}
}

// handleAcceptTCP 处理新连接
func (driver *GateDriver) handleAccept(stream *Stream, conn net.Conn, websocketConn *websocket.Conn) {
	// TODO 交换密钥
	// msg, err := ReadMessage(stream)
	var key []byte
	channel, err := driver.newGateSession(key, conn, websocketConn)
	if err != nil {
		log.Errorf("driver: %s new channel: %s err: %s", driver, conn.RemoteAddr(), err)
		_ = conn.Close()
		return
	}
	// websocket 特殊处理
	if driver.protocol == ProtocolWebsocket {
		<-channel.exit
	}
}

// runWebsocket 以websocket协议启动网关驱动
func (driver *GateDriver) runWebsocket() {
	log.Infof("start websocket listen: %s", driver.localAddr)
	if err := http.ListenAndServe(driver.localAddr, nil); err != nil {
		log.Errorf("websocket listen: %s err: %s", driver.localAddr, err)
		time.AfterFunc(config.ListenRetryInterval(), driver.run)
		return
	}
}

// handleAcceptWebSocket 处理新的websocket连接
func (driver *GateDriver) handleAcceptWebsocket(w http.ResponseWriter, r *http.Request) {
	websocketConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("websocket upgrade: %s err: %s", driver.localAddr, err)
		return
	}
	defer func() {
		if err = websocketConn.Close(); err != nil {
			log.Errorf("websocket close remote: %s err: %s", websocketConn.RemoteAddr(), err)
		}
	}()
	stream := NewWebsocketStream(websocketConn)
	go driver.handleAccept(stream, nil, websocketConn)
}
