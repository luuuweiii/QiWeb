package znet

import (
	"fmt"
	"myzinx/utils"
	"myzinx/ziface"
	"net"
)

// IServer接口的具体实现，定义一个Server服务器模块
type Server struct {
	// 服务器名称
	Name string
	// 服务器绑定的IP版本
	IPVersion string
	// 服务器监听的IP
	IP string
	// 服务器监听的端口
	Port int

	// 当前server的消息管理模块，用来绑定MsgID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandle
}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name :%s, listenner at IP: %s, Port:%d is starting...\n", s.Name, s.IP, s.Port)
	fmt.Printf("[Zinx] Version %s, MaxConn:%d, MaxPacketSize:%d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPackageSize)

	// 1.获取一个TCP的Addr
	go func() {
		// 将地址解析为一个结构体
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addt error :", err)
			return
		}

		// 2.监听服务器的地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("lisnten", s.IPVersion, " err", err)
			return
		}

		fmt.Println("start Zinx server succ, ", s.Name, " succ, Listenning...")
		// cid写在这里合适吗
		var cid uint32
		cid = 0
		// 3.阻塞的等待客户端连接，处理客户端的业务（读写）
		for {
			// 如果有客户端连接过来，阻塞会返回 （握手）
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			// 已经与客户端建立连接基本读写操作
			// 将处理新链接的业务方法和conn进行绑定 得到我们的链接模块
			dealConn := NewConnection(conn, cid, s.MsgHandler)
			cid++

			//启动当前的链接业务处理
			dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	// TODO 将服务器的资源、状态或者一些已经开辟的连接信息 进行停止或者回收
}

func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()

	// TODO 做一些启动服务之后的额外业务

	// 阻塞状态
	select {}
}

// 添加路由方法
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddTouter(msgID, router)
	fmt.Println("Add Router Success")
}

// 初始化Server模块的方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
	}
	return s
}
