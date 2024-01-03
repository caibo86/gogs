// 状态
enum SessionStatus {
	Closed       = 1;
	Connecting   = 2;
	Disconnected = 3;
	InConnected  = 4;
	OutConnected = 5;
}

// 消息类型
enum MessageType {
	WhoAmI = 1;
	Accept = 2;
	Reject = 3;
    Call   = 4;
    Return = 5;
}

// 消息
struct Message {
	Type MessageType = 1;
	Data bytes       = 2;
}

// 一次调用
struct Call {
	ID      uint32  = 1;
	Method  uint32  = 2;
	Params  []bytes = 3; // 序列化后的入参
	Service uint16  = 4;
}

// 返回
struct Return {
	ID      uint32  = 1;
	Params  []bytes = 2; // 序列化后的返回值
	Service uint16  = 3;
}

