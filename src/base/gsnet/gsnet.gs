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

// 序列化之后的参数
struct Param {
    Data bytes = 1;
}

// 一次调用
struct Call {
    ID uint16 = 1;
    Service  uint16 = 2;
    Method uint16 = 3;
    Params []Param = 4;
}

// 返回
struct Return {
    ID uint16 = 1;
    Service  uint16 = 2;
    Params []Param = 3;
}