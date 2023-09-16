package ziface

type IServer interface {
	// 开启服务器
	Start()
	// 关闭服务器
	Stop()
	// 运行服务器
	Serve()
}
