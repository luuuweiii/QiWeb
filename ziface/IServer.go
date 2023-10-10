package ziface

type IServer interface {
	// 开启服务器
	Start()
	// 关闭服务器
	Stop()
	// 运行服务器
	Serve()
	//路由功能：给当前的服务注册一个路由方法，供客户端的链接处理使用
	AddRouter(msgID uint32, router IRouter)
	// 获取当前server的链接管理器
	GetConnMgr() IconnManager
	// 注册OnConnStart Hook函数
	SetOnConnStart(func(connection IConnection))
	// 注册OnConnStop Hook函数
	SetOnConnStop(func(connection IConnection))
	// 调用OnConnStart Hook函数
	CallOnConnStart(connection IConnection)
	// 调用OnConnStop Hook函数
	CallOnConnStop(connection IConnection)
}
