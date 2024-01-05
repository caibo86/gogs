// 用户信息
struct UserInfo {
	UserID   int64  = 1; 
	Username string = 2; 
	Email    string = 3; 
	Power    int32  = 4; 
}

// 游戏服API
service Game {
	GetServerTime() -> (int64, Err); // 获取服务器时间,毫秒
}

// 用户API
service User(Game) {
	GetUserInfo() -> (UserInfo, Err); // 获取用户信息
}

