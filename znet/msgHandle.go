package znet

import (
	"fmt"
	"myzinx/utils"
	"myzinx/ziface"
	"strconv"
)

/*
	消息处理模块的实现
*/

type MsgHandle struct {
	// 存放每个MsgID 对应的处理方法
	Apis map[uint32]ziface.IRouter

	// 负责Worker取任务的消息队列
	TaskQueue []chan ziface.IRequest

	// 业务工作Worker池的worker数量
	WorkerPoolSize uint32
}

// 初始化/创建MsgHandle方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, // 从全局配置中获取
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

// 启动一个Worker工作池 (开启工作池的动作只能发生一次，一个框架只能有一个工作池)
func (m *MsgHandle) StartWorkerPool() {
	// 根据workerPoolSize分别开启Worker，每个Worker用一个go来承载
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		// 一个worker被启动
		// 1 给当前的worker对应的channel消息队列 开辟空间 第0个worker用第0个管道
		m.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)

		// 2 启动当前的Worker，阻塞等待消息从channel传递过来
		go m.StartOneWorker(i, m.TaskQueue[i])
	}
}

// 启动一个Worker工作流程
func (m *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID =", workerID, " is started ...")

	// 不断地阻塞等待对应消息队列的消息
	for {
		select {
		// 如果有消息过俩，出列的就是一个客户端的Request，执行当前Request所绑定的业务
		case request := <-taskQueue:
			m.DoMsgHandler(request)
		}
	}
}

// 将消息交给TaskQueue，由Worker进行处理
func (m *MsgHandle) SendMessageToTaskQueue(request ziface.IRequest) {
	// 1 将消息平均分配给不同的worker
	// 根据客户端建立的ConnID来分配
	workerID := request.GetConnection().GetConnID() % m.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(),
		"request MsgID =", request.GetMsgID(),
		"to WorkerID =", workerID)

	// 2 将消息发送给对应的worker的TaskQueue即可
	m.TaskQueue[workerID] <- request
}
