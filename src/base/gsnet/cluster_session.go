// -------------------------------------------
// @file      : cluster_channel.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/3 下午4:58
// -------------------------------------------

package gsnet

import (
	"fmt"
	"gogs/base/config"
	"gogs/base/gserrors"
	log "gogs/base/logger"
	"net"
	"sync"
	"time"
)

type ConnectionType byte

const (
	ConnectionTypeIn  ConnectionType = 1 // 内连,外部发起对本地连接
	ConnectionTypeOut ConnectionType = 2 // 外连,本地发起对外部连接
)

// ClusterSession 集群节点会话
type ClusterSession struct {
	sync.WaitGroup
	remoteAddr     string          // 远程地址
	conn           net.Conn        // tcp连接
	exit           chan struct{}   // 关闭信号
	driver         *ClusterDriver  // 所属驱动
	status         SessionStatus   // 状态
	handler        ISessionHandler // 会话处理器
	cached         chan *Message   // 发送消息队列
	name           string          // 会话名字
	connectionType ConnectionType  // 连接类型
}

// newClusterSession 在指定驱动上创建一个集群节点会话,外连会话,注意此函数外层已经加锁
func (driver *ClusterDriver) newClusterSession(addr string, ct ConnectionType) (*ClusterSession, error) {
	session := &ClusterSession{
		remoteAddr:     addr,
		driver:         driver,
		status:         SessionStatusDisconnected,
		cached:         make(chan *Message, config.ClusterSessionCache()),
		name:           fmt.Sprintf("ClusterSession(%s)", addr),
		connectionType: ct,
	}
	handler, err := driver.sessionHandlerBuilder(session)
	if err != nil {
		return nil, err
	}
	session.handler = handler
	// 外连需要初始化连接,并定时检查,内连不需要
	if ct == ConnectionTypeOut {
		if session.connect() == SessionStatusClosed {
			log.Debugf("session: %s, heart beat exit return status closed", session)
		}
		go func() {
			tick := time.Tick(config.ClusterSessionHeartbeat())
			for _ = range tick {
				if session.connect() == SessionStatusClosed {
					log.Debugf("session: %s, heart beat exit return status closed", session)
					return
				}
			}
		}()
	}
	driver.sessions[addr] = session
	return session, nil
}

// String implement fmt.Stringer
func (session *ClusterSession) String() string {
	return session.name
}

// Name 获取会话名字
func (session *ClusterSession) Name() string {
	return session.name
}

// RemoteAddr 获取远程地址
func (session *ClusterSession) RemoteAddr() string {
	return session.remoteAddr
}

// DriverType 获取驱动类型
func (session *ClusterSession) DriverType() DriverType {
	return session.driver.Type()
}

// Close 关闭会话
func (session *ClusterSession) Close() {
	session.driver.lock(session, func() {
		log.Debugf("cluster session: %s closing", session)
		session.status = SessionStatusClosed
		if session.exit != nil {
			close(session.exit)
		}
		session.Wait()
		if session.conn != nil {
			_ = session.conn.Close()
		}
		session.changeStatus(SessionStatusClosed)
		session.driver.DelSession(session)
	})
}

// Handler 获取会话处理器
func (session *ClusterSession) Handler() ISessionHandler {
	return session.handler
}

// Write 向会话写入一个消息
func (session *ClusterSession) Write(msg *Message) error {
	if session.status == SessionStatusClosed {
		return gserrors.Newf("cluster %s session: %s closed", session)
	}
	select {
	case session.cached <- msg:
		return nil
	default:
		return gserrors.Newf("cluster session: %s sending queue overflow: %d", session, len(session.cached))
	}
}

// Status 获取会话状态
func (session *ClusterSession) Status() SessionStatus {
	return session.status
}

// changeStatus 修改会话状态
func (session *ClusterSession) changeStatus(status SessionStatus) {
	session.status = status
	session.handler.SessionStatusChanged(status)
	log.Debugf("cluster session: %s status changed: %s", session, status)
}

// connect 连接该会话,加锁异步外连
func (session *ClusterSession) connect() SessionStatus {
	session.driver.lock(session, func() {
		switch session.status {
		case SessionStatusDisconnected:
			// 连接
			log.Infof("cluster session: %s connecting", session)
			session.changeStatus(SessionStatusConnecting)
			go session.outConnect()
		}
	})
	return session.status
}

// outConnect 外连
func (session *ClusterSession) outConnect() {
	// 连接
	conn, err := net.Dial("tcp", session.remoteAddr)
	if err != nil {
		log.Errorf("cluster session: %s out connect err: %s", session, err)
		session.closeConn(nil)
		return
	}
	stream := NewStream(conn, conn)
	// 发送握手消息
	msg := NewMessage()
	msg.Type = MessageTypeHandshake
	msg.Data = []byte(conn.LocalAddr().String())
	// 发送
	err = WriteMessage(stream, msg)
	if err != nil {
		log.Errorf("cluster session: %s handshake err: %s", session, err)
		session.closeConn(conn)
		return
	}
	// 读取一个消息
	msg, err = ReadMessage(stream)
	if err != nil {
		log.Errorf("cluster session: %s handshake err: %s", session, err)
		session.closeConn(conn)
		return
	}
	if msg.Type != MessageTypeAccept {
		log.Errorf("cluster session: %s handshake err: %s", session, msg.Type)
		session.closeConn(conn)
		return
	}
	// 完成外连握手后的设置
	exit := session.outConnection(conn)
	if exit == nil {
		log.Errorf("cluster session: %s drop out connection: %s", session, conn)
		session.closeConn(conn)
		return
	}
	log.Debugf("cluster session: %s out connection established", session)
	// 开始读写
	go session.recvLoop(conn)
	go session.sendLoop(conn, exit)
}

// recvLoop 接收循环
func (session *ClusterSession) recvLoop(conn net.Conn) {
	stream := NewStream(conn, conn)
	for {
		msg, err := ReadMessage(stream)
		if err != nil {
			if session.connectionType == ConnectionTypeOut {
				session.closeConn(conn)
			} else if session.connectionType == ConnectionTypeIn {
				session.Close()
			}
			break
		}
		// 通知处理器,读取到一个消息
		session.handler.Read(session, msg)
	}
}

// sendLoop 发送循环
func (session *ClusterSession) sendLoop(conn net.Conn, exit chan struct{}) {
	session.Add(1)
	defer session.Done()
	stream := NewStream(conn, conn)
	for {
		select {
		case msg, ok := <-session.cached:
			if !ok {
				log.Debugf("cluster session: %s cache closed", session)
				return
			}
			err := WriteMessage(stream, msg)
			if err != nil {
				session.closeConn(conn)
				log.Debugf("cluster session: %s send loop err: %s", session, err)
				return
			}
		case <-exit:
			return
		}
	}
}

// closeConn 关闭会话的远程连接
func (session *ClusterSession) closeConn(conn net.Conn) {
	if conn != nil {
		_ = conn.Close()
	}
	session.driver.lock(session, func() {
		switch session.status {
		case SessionStatusOutConnected, SessionStatusInConnected:
			if conn == session.conn {
				if session.exit != nil {
					close(session.exit)
				}
				session.conn = nil
				session.exit = nil
				session.changeStatus(SessionStatusDisconnected)
			}
		case SessionStatusConnecting:
			if session.conn != nil || session.exit != nil {
				gserrors.Panic("check outConnection or inConnection implement")
			}
			session.changeStatus(SessionStatusDisconnected)
		}
	})
}

// outConnection 外连成功,加锁异步设置
func (session *ClusterSession) outConnection(conn net.Conn) chan struct{} {
	var exit chan struct{}
	session.driver.lock(session, func() {
		switch session.status {
		case SessionStatusDisconnected, SessionStatusConnecting:
			session.conn = conn
			session.exit = make(chan struct{})
			exit = session.exit
			session.changeStatus(SessionStatusOutConnected)
		default:
			log.Errorf("session: %s outConnection when status: %s", session, session.status)
		}
	})
	return exit
}

// inConnection 内连成功,加锁异步设置
func (session *ClusterSession) inConnection(conn net.Conn) chan struct{} {
	var exit chan struct{}
	session.driver.lock(session, func() {
		switch session.status {
		case SessionStatusDisconnected:
			session.conn = conn
			session.exit = make(chan struct{})
			exit = session.exit
			session.changeStatus(SessionStatusInConnected)
		case SessionStatusConnecting:
			if session.driver.localAddr < session.remoteAddr {
				session.conn = conn
				session.exit = make(chan struct{})
				exit = session.exit
				session.changeStatus(SessionStatusInConnected)
			}
		default:
			log.Errorf("session: %s inConnection when status: %s", session, session.status)
		}
	})
	return exit
}
