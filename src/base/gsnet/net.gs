// 会话状态
enum SessionStatus {
	Closed       = 1;
	Connecting   = 2;
	Disconnected = 3;
	InConnected  = 4;
	OutConnected = 5;
}

// 服务状态
enum ServiceStatus {
	Online      = 1;
	Offline     = 2;
	Unreachable = 3;
}

// 消息类型
enum MessageType {
	Handshake = 1; // 握手消息
	Accept    = 2; // 握手成功
	Reject    = 3; // 握手拒绝
	Call      = 4; // 服务调用
	Return    = 5; // 服务调用返回
	Registry  = 6; // 服务注册
}

// 服务注册
struct ServiceRegistry {
	Add         bool   = 1;
	ServiceID   uint32 = 2; // 服务类型名字
	ServiceType string = 3; // 服务名字
	ServiceName string = 4; // 服务ID
}

// 服务注册列表
struct ServiceRegistryData {
	Data []ServiceRegistry = 1;
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

