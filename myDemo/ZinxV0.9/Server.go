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

// hello Zinx test 自定义路由
type HelloZinxRouter struct {
	znet.BaseRouter
}

// Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle...")
	// 先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client:msgID= ", request.GetMsgID(),
		", data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println("call back ping...ping...ping error")
	}
}

// Test Handle
func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle...")
	// 先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client:msgID= ", request.GetMsgID(),
		", data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(201, []byte("hello...hello...hello"))
	if err != nil {
		fmt.Println("call back ping...ping...ping error")
	}
}

// 创建钩子函数
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("====> DoConnectionBegin is Called ...")
	if err := conn.SendMsg(202, []byte("DoConnection Begin")); err != nil {
		fmt.Println(err)
	}
}

func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("====> DoConnectionLost is Called ...")
	fmt.Println("conn ID = ", conn.GetConnID(), "is Lost ...")
}

func main() {
	// 1 创建一个server句柄，使用Zinx的api
	s := znet.NewServer("[zinx V0.9]")

	// 2 注册链接Hook函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	// 3 给当前zinx框架添加router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})

	// 4 启动server
	s.Serve()
}
