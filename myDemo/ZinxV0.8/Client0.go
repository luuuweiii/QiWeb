package main

import (
	"fmt"
	"io"
	"myzinx/znet"
	"net"
	"time"
)

// 模拟客户端

func main() {
	fmt.Println("client0 start...")

	time.Sleep(1 * time.Second)

	// 1 连接远程服务器，得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err,exit!")
	}

	//2 调用Write 写数据
	for {
		// 发送封包的message消息 MsgId:0
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("Hello Zinx client0 V0.8..")))
		if err != nil {
			fmt.Println("Pack error", err)
			return
		}

		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("Pack error", err)
			return
		}

		// 服务器应该给我们回复一个message数据， Msg：1 pingping (客户端不能用send方法所以只能重新写)

		// 1 先读取流中的head部分 得到ID 和 dataLen
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error", err)
		}

		// 将二进制的head拆包到msg结构体中
		msgHead, err2 := dp.UnPack(binaryHead)
		if err2 != nil {
			fmt.Println("client unpack msgHead error", err)
		}

		if msgHead.GetMsgLen() > 0 {
			// 2 msg里面有数据，根据DataLen进行第二次读取，将data读取出来
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error ,", err)
				return
			}

			fmt.Println("---> Recv Server Msg :ID=", msg.GetMsgId(),
				", len = ", msg.GetMsgLen(),
				", dat = ", string(msg.GetData()))
		}

		// cpu阻塞
		time.Sleep(1 * time.Second)
	}
}
