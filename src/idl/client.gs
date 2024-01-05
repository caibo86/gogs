// 客户端接口,用于服务器调用
service ClientAPI {
	GetClientInfo() -> (ClientInfo); // 获取客户端的设备信息
}

