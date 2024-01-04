package mongodb

// OpType 操作类型
type OpType int8

const (
	opTypeInsert       OpType = 1
	opTypeUpdate       OpType = 2
	opTypeUpsert       OpType = 3
	opTypeDelete       OpType = 4
	opTypeUploadGridFS OpType = 5
	opTypeReplace      OpType = 6
	opTypeFind         OpType = 7
)

func (t OpType) String() string {
	switch t {
	case opTypeInsert:
		return "Insert"
	case opTypeUpdate:
		return "Update"
	case opTypeUpsert:
		return "Upsert"
	case opTypeDelete:
		return "Insert"
	case opTypeUploadGridFS:
		return "UploadGridFS"
	case opTypeReplace:
		return "Replace"
	}
	return "Unknown"
}

// OpReq mongo操作
type OpReq struct {
	OpType
	Collection string
	Filename   string
	Filter     interface{}
	Doc        interface{}
	NeedRet    bool // 需要异步结果返回
	Context    interface{}
	Event      Event
	//buffer     *bytes.Buffer
}

func NewOpReq(opType OpType, collection, filename string, filter, doc interface{}, needRet bool, event Event, context interface{}) *OpReq {
	op := &OpReq{
		OpType:     opType,
		Collection: collection,
		Filename:   filename,
		Filter:     filter,
		Doc:        doc,
		NeedRet:    needRet,
		Context:    context,
		Event:      event,
	}
	return op
}
