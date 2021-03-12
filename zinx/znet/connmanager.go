package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

// ConnManager 连接管理模块
type ConnManager struct {
	connections map[uint32]ziface.IConnection // 管理的连接集合
	connLock    sync.RWMutex                  // 保护连接集合的读写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (cm *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 将conn加入到ConnManager中
	cm.connections[conn.GetConnID()] = conn
	fmt.Println("Add ConnID =", conn.GetConnID(), "to ConnManager success, conn num =", cm.Len())
}

func (cm *ConnManager) Remove(conn ziface.IConnection) {
	// 保护共享资源，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	delete(cm.connections, conn.GetConnID())
	fmt.Println("Remove ConnID =", conn.GetConnID(), "to ConnManager success, conn num =", cm.Len())
}

func (cm *ConnManager) Get(connId uint32) (ziface.IConnection, error) {
	// 保护共享资源，加读锁
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	if conn, ok := cm.connections[connId]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection NOT FOUND")
	}
}

func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

func (cm *ConnManager) Clear() {
	// 保护共享资源，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 删除conn并停止conn的工作
	for connID, conn := range cm.connections {
		// 停止
		conn.Stop()
		// 删除
		delete(cm.connections, connID)
	}
	fmt.Println("Clear all connections success, conn num =", cm.Len())
}
