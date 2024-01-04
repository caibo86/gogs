// 账号类型
enum AccountType {
	Blank    = 0; // 空
	Test     = 1; // 测试账号
	Platform = 2; // 平台账号
}

// 登录请求
struct LoginReq {
	AccountID   int64       = 1; // 账号ID
	Token       string      = 2; // 登录token
	UserID      int64       = 3; // 用户ID
	ServerID    int64       = 4; // 服务器ID
	AccountType AccountType = 5; // 账号类型
}

// 登录应答
struct LoginAck {
	AccountID int64 = 1; 
	UserID    int64 = 2; 
	ServerID  int64 = 3; 
}

// 客户端信息
struct ClientInfo {
	OpenUDID             string = 1; // 设备唯一标识
	Language             string = 2; // 语言
	OS                   string = 3; // 操作系统
	ClientVersion        string = 4; // 客户端版本
	ClientMemoryCapacity int32  = 5; // 客户端内存容量
	ClientDeviceLevel    int32  = 6; // 客户端设备等级
	ClientChannel        string = 7; // 客户端渠道
	IsAndroidEmulator    bool   = 8; // 是否安卓模拟器
}

// 网关API
service Gate {
	Login(LoginReq, ClientInfo) -> (LoginAck, Err); // 登录
}

