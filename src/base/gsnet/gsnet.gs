// 状态
enum Status {
    Closed = 1;
}

// 消息类型
enum MessageType {
    WhoAmI = 1;
    Accept = 2;
}

// 消息
struct Message {
    Type MessageType = 1;
    Data bytes = 2;
}

// 一次调用
struct Call {
    ID uint16 = 1;
    Service  uint16 = 2;
    Method uint32 = 3;
    Params []bytes = 4; // 序列化后的入参
}

// 返回
struct Return {
    ID uint16 = 1;
    Service  uint16 = 2;
    Params []bytes = 3; // 序列化后的返回值
}