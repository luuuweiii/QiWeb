package znet

import (
	"errors"
	"fmt"
	"myzinx/ziface"
	"sync"
)

/*
	链接管理模块
*/

type ConnManager struct {
	connections map[uint32]ziface.IConnection //管理的链接集合
	connLock    sync.RWMutex                  //保护链接集合的读写锁
}

// NewConnManager 创建当前链接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// Add 添加链接
func (cm *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 将conn加入到connManager中
	cm.connections[conn.GetConnID()] = conn
	fmt.Println("connID = ", conn.GetConnID(), ", add to ConnManager successfully: conn num = ", cm.Len())
}

// Remove 删除链接
func (cm *ConnManager) Remove(conn ziface.IConnection) {
	// 保护共享资源map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 删除链接信息
	delete(cm.connections, conn.GetConnID())
	fmt.Println("connID = ", conn.GetConnID(), ", remove from ConnManager successfully: conn num = ", cm.Len())
}

// Get 根据connID获取链接
func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	// 保护共享资源map，加读锁
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	if conn, ok := cm.connections[connID]; ok {
		// 链接存在
		return conn, nil
	} else {
		return nil, errors.New("connection not FOUND!")
	}
}

// Len 得到当前链接总数
func (cm *ConnManager) Len() int {
	// 保护共享资源map，加读锁 （为什么这里加锁会出问题）
	//cm.connLock.RLock()
	//defer cm.connLock.RUnlock()

	return len(cm.connections)
}

// ClearConn 清除并终止所有链接
func (cm *ConnManager) ClearConn() {
	// 保护共享资源map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Lock()

	// 删除conn并停止conn的工作
	for connID, conn := range cm.connections {
		// 停止
		conn.Stop()
		// TODO 这里读进程可能无法被关闭

		// 删除
		delete(cm.connections, connID)
	}
	fmt.Println("Clear All connections succ! conn num = ", cm.Len())
}
