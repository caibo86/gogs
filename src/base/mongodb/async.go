package mongodb

import (
	"fmt"
	"go.uber.org/zap"
	log "gogs/base/logger"
	"runtime/debug"
)

const (
	DefaultAsyncReqChanLength = 1024 // Req
	DefaultAsyncRetChanLength = 1024 // Ret
)

type Event int32

const (
	EventSaveUser Event = iota + 1
)

// AsyncRet 异步操作的结果
type AsyncRet struct {
	Event   Event
	Context interface{}
	Result  interface{}
	Err     error
}

func (mc *MongoClient) AsyncStartWork() {
	mc.reqChan = make(chan *OpReq, DefaultAsyncReqChanLength)
	mc.closeChan = make(chan struct{}, 1)
	mc.async = true
	mc.AsyncRetChan = make(chan *AsyncRet, DefaultAsyncRetChanLength)
	go mc.loop()
}

func (mc *MongoClient) AsyncEndWork() {
	close(mc.reqChan)
	<-mc.closeChan
	close(mc.closeChan)
}

func (mc *MongoClient) loop() {
	defer func() {
		if err := recover(); err != nil {
			log.Errorw("mongo client panic", zap.Stack("stack"))
		}
		mc.async = false
		mc.closeChan <- struct{}{}
	}()
	for {
		select {
		case req, ok := <-mc.reqChan:
			if !ok {
				log.Infof("mongo client req channel is closed")
				return
			}
			ret := mc.handleOpReq(req)
			if ret.Err != nil {
				log.Errorf("mongo client op(%s:%s) err:%s stack:%s",
					req.OpType, req.Collection, ret.Err, debug.Stack())
			}
			if req.NeedRet {
				mc.AsyncRetChan <- ret
			}

		}
	}
}

func (mc *MongoClient) handleOpReq(req *OpReq) *AsyncRet {
	ret := &AsyncRet{
		Event:   req.Event,
		Context: req.Context,
	}
	switch req.OpType {
	case opTypeInsert:
		ret.Result, ret.Err = mc.doInsertOne(req.Collection, req.Doc)
	case opTypeUpdate:
		ret.Result, ret.Err = mc.doUpdateOne(req.Collection, req.Filter, req.Doc)
	case opTypeUpsert:
		ret.Result, ret.Err = mc.doUpsertOne(req.Collection, req.Filter, req.Doc)
	case opTypeReplace:
		ret.Result, ret.Err = mc.doReplaceOne(req.Collection, req.Filter, req.Doc)
	case opTypeDelete:
		ret.Result, ret.Err = mc.doDeleteOne(req.Collection, req.Filter)
	case opTypeFind:
		ret.Result, ret.Err = mc.doFindOne(req.Collection, req.Filter)
	case opTypeUploadGridFS:
		data, ok := req.Doc.([]byte)
		if !ok {
			ret.Err = fmt.Errorf("invalid data for upload grid fs")
		} else {
			ret.Err = mc.doWriteGridFile(req.Filename, req.Collection, data)
		}
	default:
		ret.Err = fmt.Errorf("unknown db op type:%d", req.OpType)
	}
	return ret
}
