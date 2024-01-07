// -------------------------------------------
// @file      : stream.go
// @author    : 蔡波
// @contact   : caibo923@gmail.com
// @time      : 2024/1/3 下午2:15
// -------------------------------------------

package network

import (
	"github.com/gorilla/websocket"
	"gogs/base/cberrors"
	"io"
)

// Stream 流 带缓冲
type Stream struct {
	reader        io.Reader
	writer        io.Writer
	protocol      ProtocolType
	websocketConn *websocket.Conn
}

// NewStream 创建流
func NewStream(reader io.Reader, writer io.Writer) *Stream {
	stream := &Stream{
		reader:   reader,
		writer:   writer,
		protocol: ProtocolTCP,
	}
	return stream
}

// NewWebsocketStream 创建websocket流
func NewWebsocketStream(wsConn *websocket.Conn) *Stream {
	stream := &Stream{
		websocketConn: wsConn,
		protocol:      ProtocolWebsocket,
	}
	return stream
}

// Read 读取数据,先读到缓冲区再读取
func (stream *Stream) Read(buf []byte) (int, error) {
	if stream.protocol == ProtocolWebsocket {
		t, r, err := stream.websocketConn.NextReader()
		if t != websocket.BinaryMessage {
			return 0, cberrors.New("invalid websocket message type: %d", t)
		}
		if err != nil {
			return 0, err
		}
		return io.ReadFull(r, buf)
	}
	return io.ReadFull(stream.reader, buf)
}

// Write 写入数据到缓冲区
func (stream *Stream) Write(buf []byte) (int, error) {
	if stream.protocol == ProtocolWebsocket {
		length := len(buf)
		return length, stream.websocketConn.WriteMessage(websocket.BinaryMessage, buf)
	}
	return stream.writer.Write(buf)
}
