// 状态
enum SessionStatus {
	Closed       = 1;
	Connecting   = 2;
	Disconnected = 3;
	InConnected  = 4;
	OutConnected = 5;
}

enum ServiceStatus {
    Online = 1;
    Offline  = 2;
    Unreachable   = 3;
}

// 消息类型
enum MessageType {
	WhoAmI = 1; // 握手消息
	Accept = 2; // 握手成功
	Reject = 3; // 握手拒绝
	Call   = 4; // 调用
	Return = 5; // 调用返回
	CER = 6; // 服务注册
}

//
struct CER {
    Add bool = 1;
    ID uint32 = 2; // 服务类型名字
    Type string = 3; // 服务名字
    Name string = 4; // 服务ID
}

struct CERs {
    Data []CER = 1;
}

// 消息
struct Message {
	Type MessageType = 1;
	Data bytes       = 2;
}

// 一次调用
struct Call {
	ID        uint32  = 1; // 流水号
	ServiceID uint32  = 2; // 服务ID
	MethodID  uint32  = 3; // 方法ID
	Params    []bytes = 4; // 序列化后的入参

}

// 返回
struct Return {
	ID        uint32  = 1; // 流水号
	ServiceID uint32  = 2; // 服务ID
	Params    []bytes = 3; // 序列化后的返回值

}

