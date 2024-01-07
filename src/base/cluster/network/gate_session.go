// -------------------------------------------
// @file      : gate_session.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/3 下午3:33
// -------------------------------------------

package network

import (
	"fmt"
	"github.com/gorilla/websocket"
	"gogs/base/cberrors"
	"gogs/base/config"
	log "gogs/base/logger"
	"net"
)

// GateSession 网关会话
type GateSession struct {
	conn          net.Conn        // tcp连接
	websocketConn *websocket.Conn // websocket连接
	driver        *GateDriver     // 所属驱动
	status        SessionStatus   // 状态
	handler       ISessionHandler // 会话处理器
	cached        chan *Message   // 发送消息队列
	name          string          // 会话名字
	key           []byte          // AES密钥
	exit          chan struct{}   // 结束信号
}

// newGateSession 新建网关会话
func (driver *GateDriver) newGateSession(key []byte, conn net.Conn, websocketConn *websocket.Conn) (*GateSession, error) {
	remoteAddr := conn.RemoteAddr().String()
	session := &GateSession{
		conn:          conn,
		websocketConn: websocketConn,
		driver:        driver,
		status:        SessionStatusInConnected,
		cached:        make(chan *Message, config.GateSessionCache()),
		name:          fmt.Sprintf("GateSession(%s->%s)", remoteAddr, driver.localAddr),
		key:           key,
		exit:          make(chan struct{}, 1),
	}
	// 创建会话处理器
	handler, err := driver.sessionHandlerBuilder(session)
	if err != nil {
		return nil, cberrors.New("%s session(%s) create session handler err: %s", driver, session, err)
	}
	session.handler = handler
	// 启动发送和接收
	go session.recvLoop()
	go session.sendLoop()
	// 加入到所属驱动
	driver.Lock()
	driver.remotes[remoteAddr] = session
	driver.Unlock()
	// 通知处理器,会话状态变更
	session.handler.SessionStatusChanged(session.status)
	return session, nil
}

// String implement fmt.Stringer
func (session *GateSession) String() string {
	return session.name
}

// Name 获取会话名字
func (session *GateSession) Name() string {
	return session.name
}

// Close 关闭会话
func (session *GateSession) Close() {
	// 修改状态
	session.status = SessionStatusClosed
	// 关闭连接
	_ = session.conn.Close()
	// 取消注册
	session.driver.DelSession(session)
	// 通知处理器,会话状态变更
	session.handler.SessionStatusChanged(session.status)
}

// Handler 获取会话处理器
func (session *GateSession) Handler() ISessionHandler {
	return session.handler
}

// Write 向会话写入一个消息
func (session *GateSession) Write(msg *Message) error {
	if session.status == SessionStatusClosed {
		return cberrors.New("cluster session: %s closed", session)
	}
	select {
	case session.cached <- msg:
		return nil
	default:
		return cberrors.New("cluster session: %s sending queue overflow: %d", session, len(session.cached))
	}
}

// Status 获取会话状态
func (session *GateSession) Status() SessionStatus {
	return session.status
}

// changeStatus 修改会话状态
func (session *GateSession) changeStatus(status SessionStatus) {
	session.status = status
	session.handler.SessionStatusChanged(status)
}

// DriverType 获取会话所属驱动类型
func (session *GateSession) DriverType() DriverType {
	return session.driver.Type()
}

// recvLoop 接收循环
func (session *GateSession) recvLoop() {
	var stream *Stream
	if session.driver.protocol == ProtocolTCP {
		stream = NewStream(session.conn, session.conn)
	} else {
		stream = NewWebsocketStream(session.websocketConn)
	}
	for {
		msg, err := ReadMessage(stream)
		log.Infof("session: %s recv msg: %+v", session, msg)
		if err != nil {
			session.Close()
			log.Debugf("%s session: %s recv loop err: %s", session.driver, session, err)
			break
		}
		// 通知处理器,读取到一个消息
		session.handler.Read(session, msg)
	}
	// 关闭信号
	session.exit <- struct{}{}
}

// sendLoop 发送循环
func (session *GateSession) sendLoop() {
	var stream *Stream
	if session.driver.protocol == ProtocolTCP {
		stream = NewStream(session.conn, session.conn)
	} else {
		stream = NewWebsocketStream(session.websocketConn)
	}
	for msg := range session.cached {
		err := WriteMessage(stream, msg)
		if err != nil {
			session.Close()
			log.Debugf("%s session: %s send loop err: %s", session.driver, session, err)
			break
		}
	}
	// 关闭信号
	session.exit <- struct{}{}
}
