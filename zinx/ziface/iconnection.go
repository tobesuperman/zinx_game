package ziface

import "net"

// IConnection 定义连接模块的抽象层
type IConnection interface {
	Start()                                      // 启动连接（让当前的连接准备开始工作）
	Stop()                                       // 停止连接（结束当前连接的工作）
	GetTCPConnection() *net.TCPConn              // 获取当前连接所绑定的socket
	GetConnID() uint32                           // 获取当前连接模块的ID
	RemoteAddr() net.Addr                        // 获取远程客户端的TCP状态（包括IP和端口）
	SendMsg(msgId uint32, data []byte) error     // 发送数据（将数据发送给远程的客户端）
	SetProperty(key string, value interface{})   // 设置连接属性
	GetProperty(key string) (interface{}, error) // 获取连接属性
	RemoveProperty(key string)                   // 删除连接属性
}

// HandleFunc 定义一个处理连接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
