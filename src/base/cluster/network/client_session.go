// -------------------------------------------
// @file      : client_session.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/3 下午7:12
// -------------------------------------------

package network

import (
	"fmt"
	"github.com/gorilla/websocket"
	"gogs/base/cberrors"
	log "gogs/base/logger"
	"net"
)

// ClientSession 客户端会话
type ClientSession struct {
	driver    *ClientDriver   // 所属驱动
	handler   ISessionHandler // 会话处理器
	name      string          // 会话名字,在driver中唯一
	key       []byte          // AES加密密钥
	conn      net.Conn        // tcp连接
	websocket *websocket.Conn // websocket连接
	exit      chan struct{}   // 关闭信号
	cached    chan *Message   // 发送消息队列
	status    SessionStatus   // 状态
	protocol  ProtocolType    // 协议类型
}

// String implements fmt.Stringer
func (session *ClientSession) String() string {
	return fmt.Sprintf("ClientSession(%s)", session.name)
}

// Name implements ISession
func (session *ClientSession) Name() string {
	return session.name
}

// Status implements ISession
func (session *ClientSession) Status() SessionStatus {
	return session.status
}

// DriverType implements ISession
func (session *ClientSession) DriverType() DriverType {
	return session.driver.Type()
}

// Close implements ISession
func (session *ClientSession) Close() {
	session.disconnect()
	session.driver.DelSession(session)
}

// Handler implements ISession
func (session *ClientSession) Handler() ISessionHandler {
	return session.handler
}

// disconnect 断开连接
func (session *ClientSession) disconnect() {
	if session.protocol == ProtocolTCP && session.conn != nil {
		_ = session.conn.Close()
	} else if session.protocol == ProtocolWebsocket && session.websocket != nil {
		_ = session.websocket.Close()
	}
	// 加锁回调
	session.driver.lock(session, func() {
		switch session.status {
		case SessionStatusOutConnected, SessionStatusConnecting:
			if session.exit != nil {
				close(session.exit)
			}
			session.conn = nil
			session.websocket = nil
			session.key = nil
			session.exit = nil
			session.status = SessionStatusDisconnected
			session.handler.SessionStatusChanged(session.status)
		}
	})
}

// outConnection 外连后的设置
func (session *ClientSession) outConnection(key []byte, conn net.Conn, websocketConn *websocket.Conn) chan struct{} {
	var exit chan struct{}
	// 加锁回调
	session.driver.lock(session, func() {
		switch session.status {
		case SessionStatusDisconnected, SessionStatusConnecting:
			if session.protocol == ProtocolWebsocket {
				session.websocket = websocketConn
			} else if session.protocol == ProtocolTCP {
				session.conn = conn
			} else {
				cberrors.Panic("unsupported protocol: %s", session.protocol)
			}
			session.conn = conn
			session.key = key
			session.status = SessionStatusOutConnected
			session.exit = make(chan struct{})
			session.handler.SessionStatusChanged(session.status)
			exit = session.exit
		}
	})
	return exit
}

// connect 连接
func (session *ClientSession) connect() {
	ok := false
	session.driver.lock(session, func() {
		if session.status == SessionStatusDisconnected {
			ok = true
			session.status = SessionStatusConnecting
			session.handler.SessionStatusChanged(session.status)
		}
	})
	if !ok {
		log.Debug("client session: %s status: %s skip connect", session, session.status)
		return
	}
	// 连接
	var conn net.Conn
	var err error
	var stream *Stream
	var websocketConn *websocket.Conn
	switch session.driver.protocol {
	case ProtocolTCP:
		conn, err = net.Dial("tcp", session.driver.remoteAddr)
		if err != nil {
			log.Errorf("client session: %s dial: %s err: %s", session, session.driver.remoteAddr, err)
			session.disconnect()
			return
		}
		stream = NewStream(conn, conn)
	case ProtocolWebsocket:
		websocketConn, _, err = websocket.DefaultDialer.Dial(session.driver.remoteAddr, nil)
		if err != nil {
			log.Errorf("client session: %s dial: %s err: %s", session, session.driver.remoteAddr, err)
			session.disconnect()
			return
		}
		stream = NewWebsocketStream(websocketConn)
	default:
		cberrors.Panic("unsupported protocol: %s", session.driver.protocol)
	}
	// TODO 密钥交换
	_ = stream
	var key []byte
	exit := session.outConnection(key, conn, websocketConn)
	if exit == nil {
		log.Debugf("client session: %s drop out connection: %s", session, conn)
		session.disconnect()
		return
	}
	log.Infof("client session: %s connected", session)
	// 启动会话
	go session.recvLoop()
	go session.sendLoop()
}

// recvLoop 接收循环
func (session *ClientSession) recvLoop() {
	var stream *Stream
	if session.protocol == ProtocolTCP {
		stream = NewStream(session.conn, session.conn)
	} else if session.protocol == ProtocolWebsocket {
		stream = NewWebsocketStream(session.websocket)
	} else {
		cberrors.Panic("unsupported protocol: %s", session.protocol)
	}
	for {
		msg, err := ReadMessage(stream)
		if err != nil {
			session.disconnect()
			log.Errorf("client session: %s read message err: %s", session, err)
			break
		}
		session.handler.Read(session, msg)
	}
}

// sendLoop 发送循环
func (session *ClientSession) sendLoop() {
	var stream *Stream
	if session.protocol == ProtocolTCP {
		stream = NewStream(session.conn, session.conn)
	} else if session.protocol == ProtocolWebsocket {
		stream = NewWebsocketStream(session.websocket)
	} else {
		cberrors.Panic("unsupported protocol: %s", session.protocol)
	}
	for {
		select {
		case msg := <-session.cached:
			err := WriteMessage(stream, msg)
			if err != nil {
				session.disconnect()
				log.Errorf("client session: %s write message err: %s", session, err)
				return
			}
		case <-session.exit:
			return
		}
	}
}

// Write 发送一个Message
func (session *ClientSession) Write(msg *Message) error {
	if session.status == SessionStatusClosed {
		return cberrors.New("client session: %s closed", session)
	}
	select {
	case session.cached <- msg:
		return nil
	default:
		return cberrors.New("client session: %s sending queue overflow: %d", session, len(session.cached))
	}
}
