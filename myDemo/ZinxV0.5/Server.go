package main

import (
	"fmt"
	"myzinx/ziface"
	"myzinx/znet"
)

/*
	基于Zinx框架来开发的 服务器端应用程序
*/

// ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	// 先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client:msgID= ", request.GetMsgID(),
		", data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping\n"))
	if err != nil {
		fmt.Println("call back ping...ping...ping error")
	}
}

func main() {
	s := znet.NewServer("[zinx V0.5]")
	s.AddRouter(&PingRouter{})
	s.Serve()
}
