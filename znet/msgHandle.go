package znet

import (
	"fmt"
	"myzinx/ziface"
	"strconv"
)

/*
	消息处理模块的实现
*/

type MsgHandle struct {
	// 存放每个MsgID 对应的处理方法
	Apis map[uint32]ziface.IRouter
}

// 初始化/创建MsgHandle方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

// 调度/执行对应的Router消息处理方法
func (m *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	// 1 从Request中找到msgID
	handler, ok := m.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), " is NOT FOUND! Need Register!")
	}
	//2 更具MSgID 调度对应的router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (m *MsgHandle) AddTouter(msgID uint32, router ziface.IRouter) {
	// 1 判断当前msg绑定的api处理方法是否存在
	if _, ok := m.Apis[msgID]; ok {
		// id 已经注册
		panic("repeat api, msgID = " + strconv.Itoa(int(msgID)))
	}
	// 2 添加msg与API的绑定关系
	m.Apis[msgID] = router
	fmt.Println("Add api MsgID =", msgID, " succ!")
}
