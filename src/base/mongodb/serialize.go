package mongodb

import (
	"bytes"
	"encoding/gob"
	"sync"
)

var BufferPool = NewPool()

func NewPool() *Pool {
	return &Pool{
		Pool: &sync.Pool{
			New: func() interface{} { return &bytes.Buffer{} },
		},
	}
}

type Pool struct {
	*sync.Pool
}

func (p *Pool) Put(buf *bytes.Buffer) {
	if buf == nil {
		return
	}
	// 大于64K不再存储 直接交给GC
	if buf.Cap() >= 64<<10 {
		return
	}
	buf.Reset()
	p.Pool.Put(buf)
}

func (p *Pool) Get() *bytes.Buffer {
	return p.Pool.Get().(*bytes.Buffer)
}

func GobMarshalWithPool(v interface{}) (*bytes.Buffer, error) {
	buffer := BufferPool.Get()
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(v)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func GobMarshal(v interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(v)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func GobUnmarshal(raw []byte, v interface{}) error {
	decoder := gob.NewDecoder(bytes.NewReader(raw))
	err := decoder.Decode(v)
	return err
}
