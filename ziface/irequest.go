package ziface

/*
	IRequest接口：
	实际上把客户端请求的链接信息和请求数据包装到一个Request中
*/

type IRequest interface {
	// 得到当前链接
	GetConnection() IConnection

	// 得到
	GetData() []byte
}
