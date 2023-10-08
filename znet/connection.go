package znet

import (
	"errors"
	"fmt"
	"io"
	"myzinx/ziface"
	"net"
)

// 链接模块
type Connection struct {
	// 当前链接的socket TCP套接字
	Conn *net.TCPConn

	// 链接的ID
	ConnID uint32

	// 当前的链接状态
	isClosed bool

	// 告知当前链接已经退出/停止的 channel
	ExitChan chan bool

	// 消息的管理MsgID 和对应的业务API关系
	MsgHandler ziface.IMsgHandle
}

// 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, MsgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: MsgHandler,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
	}

	return c
}

// 读数据
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutinr is running...")

	defer fmt.Println("connID = ", c.ConnID, "Reader is exit, remote addr is ", c.Conn.RemoteAddr())
	defer c.Stop()

	for {
		// 创建一个拆包解包的对象
		dp := NewDataPack()

		// 读取客户端的Msg Head二进制流的8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error ", err)
			break
		}

		// 拆包 得到MsgID和msgDatalen放在一个msg消息中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			break
		}

		// 根据datalen读取Data，放在Msg.Data中
		var Data []byte
		if msg.GetMsgLen() > 0 {
			Data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), Data); err != nil {
				fmt.Println("read msg data error ", err)
				break
			}
		}
		msg.SetData(Data)

		// 得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		// 从路由中，找到注册绑定的Conn对应的router调用
		go c.MsgHandler.DoMsgHandler(&req)
	}
}

// 启动链接 让当前连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID= ", c.ConnID)

	// 启动从当前链接的读数据的业务
	go c.StartReader()
	//TODO 启动从当前链接写数据的业务

}

// 停止链接 结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID= ", c.ConnID)

	//如果当前链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//关闭socket链接
	c.Conn.Close()

	//关闭管道
	close(c.ExitChan)
}

// 获取当前链接绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前链接模块的链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端的TCP状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 提供一个SendMsg方法 将我们要发送给客户端的数据先进行封包再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}

	// 将data进行封包 MsgDataLen/MsgID/Data
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack eeror msg")
	}

	// 将数据发送给客户端
	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("Write msg id = ", msgId, "error :", err)
		return errors.New("conn Write error")
	}
	return nil
}
