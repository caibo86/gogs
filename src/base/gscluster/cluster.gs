import (
    "base/gsnet"
)

// cluster系统错误码
@gslang.Error
enum Err {
    OK = 0; // 正确
    RPC = 1; // RPC错误
    Timeout = 2; // 调用  超时
}

struct RProxyMsg {
    SessionID int64 = 1;
    UserID int64 = 2;
    Gate string = 3;
}

// 使用Gate双向 转发的消息
struct TunnelMsg {
    UserID int64 = 1;
    Type gsnet.MessageType = 2;
    Data []byte = 3;
}

// 投递给用户系统的消息
struct UserMsg {
    UserID int64 = 1;
    Data []byte = 2;
}

// Game上运行的服务
service GameServer {
    Login(RProxyMsg) -> (int64, Err); // 用户登录
    Logout(RProxyMsg); // 用户登出
    Tunnel(TunnelMsg); // 用户发送给Game的消息,经过Gate转发
}

// Gate上运行的服务
service GateServer {
    Tunnel(TunnelMsg); // Game发送给用户的消息,经过Gate转发
}

// Game上运行的用户系统服务
service UserSystem {
    UserInvoke(UserMsg) -> (gsnet.Return, Err); // 用户调用
}