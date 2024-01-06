// -------------------------------------------
// @file      : cluster_driver.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/3 下午4:57
// -------------------------------------------

package gsnet

import (
	"fmt"
	"gogs/base/config"
	"gogs/base/gserrors"
	log "gogs/base/logger"
	"hash/crc32"
	"net"
	"runtime"
	"sync"
	"time"
)

// ClusterDriver 集群节点驱动
type ClusterDriver struct {
	sync.RWMutex
	localAddr             string                     // 本地地址
	sessions              map[string]*ClusterSession // 会话集合
	mutexGroup            []sync.Mutex               // 会话互斥锁列表
	sessionHandlerBuilder SessionHandlerBuilder      // 会话处理器构造器
	name                  string                     // 驱动名字
}

// NewClusterDriver 新建tcp集群节点驱动,ClusterDriver有listen服务接受连内连
func NewClusterDriver(localAddr string, builder SessionHandlerBuilder) *ClusterDriver {
	driver := &ClusterDriver{
		localAddr:             localAddr,
		sessions:              make(map[string]*ClusterSession),
		mutexGroup:            make([]sync.Mutex, runtime.NumCPU()*4),
		sessionHandlerBuilder: builder,
		name:                  fmt.Sprintf("ClusterDriver(%s)", localAddr),
	}
	go driver.run()
	return driver
}

// String implements fmt.Stringer
func (driver *ClusterDriver) String() string {
	return driver.name
}

// SetBuilder 设置会话处理器构造器
func (driver *ClusterDriver) SetBuilder(builder SessionHandlerBuilder) {
	driver.sessionHandlerBuilder = builder
}

// Type implements IDriver
func (driver *ClusterDriver) Type() DriverType {
	return DriverTypeCluster
}

// GetSession 获取指定名字的会话
func (driver *ClusterDriver) GetSession(addr string) (ISession, bool) {
	driver.RLock()
	defer driver.RUnlock()
	session, ok := driver.sessions[addr]
	return session, ok
}

// NewSession 新建指定名字的会话
func (driver *ClusterDriver) NewSession(addr string, connectionType ConnectionType) (ISession, error) {
	driver.Lock()
	defer driver.Unlock()
	if session, ok := driver.sessions[addr]; ok {
		return session, gserrors.Newf("%s duplicate session addr: %s", driver, addr)
	}
	return driver.newClusterSession(addr, connectionType)
}

// DelSession 删除指定的会话
func (driver *ClusterDriver) DelSession(session ISession) {
	if driver.Type() != session.DriverType() {
		gserrors.Panicf("session driver type not match: %s != %s", driver.Type(), session.DriverType())
		return
	}
	driver.Lock()
	defer driver.Unlock()
	if _, ok := driver.sessions[session.RemoteAddr()]; ok {
		delete(driver.sessions, session.RemoteAddr())
	}
}

// Close 关闭驱动
func (driver *ClusterDriver) Close() {
	for _, session := range driver.sessions {
		session.Close()
	}
}

// lock 当目标会话在驱动上进行修改时,根据算法获取会话互斥锁,加锁后调用
func (driver *ClusterDriver) lock(session *ClusterSession, callback func()) {
	var hashcode uint32
	if len(session.RemoteAddr()) < 64 {
		scratch := make([]byte, 64)
		copy(scratch, session.RemoteAddr())
		hashcode = crc32.ChecksumIEEE(scratch) % uint32(len(driver.mutexGroup))
	} else {
		hashcode = crc32.ChecksumIEEE([]byte(session.RemoteAddr())) % uint32(len(driver.mutexGroup))
	}
	driver.mutexGroup[hashcode].Lock()
	defer driver.mutexGroup[hashcode].Unlock()
	callback()
}

// inConnection 内连,远程地址发起对本机的连接
func (driver *ClusterDriver) inConnection(whoAmI string, conn net.Conn) (*ClusterSession, chan struct{}) {
	// 根据对方身份新建一个会话
	session, err := driver.NewSession(whoAmI, ConnectionTypeIn)
	// 内连时,可复用以前同地址的断开且未关闭的会话
	if session == nil && err != nil {
		log.Errorf("inConnection(%s) err: %s", whoAmI, err)
		return nil, nil
	}
	clusterSession := session.(*ClusterSession)
	return clusterSession, clusterSession.inConnection(conn)
}

// run 启动驱动
func (driver *ClusterDriver) run() {
	// 使用驱动的本地地址 建立监听
	log.Infof("start cluster listen: %s", driver.localAddr)
	listener, err := net.Listen("tcp", driver.localAddr)
	if err != nil {
		log.Errorf("cluster listen: %s err: %s", driver.localAddr, err)
		time.AfterFunc(config.ListenRetryInterval(), driver.run)
		return
	}
	for {
		var conn net.Conn
		conn, err = listener.Accept()
		if err != nil {
			log.Errorf("cluster driver: %s accept err: %s", driver, err)
			continue
		}
		go driver.handleAccept(conn)
	}
}

// handleAccept 处理新连接
func (driver *ClusterDriver) handleAccept(conn net.Conn) {
	stream := NewStream(conn, conn)
	msg, err := ReadMessage(stream)
	if err != nil {
		log.Errorf("cluster driver: %s read message err: %s", driver, err)
		_ = conn.Close()
		return
	}
	// 第一个必须是握手消息
	if msg.Type != MessageTypeHandshake {
		log.Errorf("cluster driver: %s remote: %s except WhoAmI message, but got: %s", driver, conn.RemoteAddr(), msg.Type)
		_ = conn.Close()
		return
	}
	session, flag := driver.inConnection(string(msg.Data), conn)
	if flag != nil {
		msg.Type = MessageTypeAccept
	} else {
		msg.Type = MessageTypeReject
	}
	err = WriteMessage(stream, msg)
	if err != nil {
		log.Errorf("cluster driver: %s remote: %s write message err: %s", driver, conn.RemoteAddr(), err)
		session.closeConn(conn)
		return
	}
	// 创建会话失败
	if msg.Type == MessageTypeReject {
		_ = conn.Close()
		return
	}
	log.Debugf("cluster driver: %s session: %s new connection established", driver, session)
	go session.recvLoop(conn)
	go session.sendLoop(conn, flag)
}
