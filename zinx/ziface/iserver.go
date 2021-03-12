package ziface

// IServer 定义一个服务器接口
type IServer interface {
	Start()                                 // 启动服务器
	Stop()                                  // 停止服务器
	Serve()                                 // 运行服务器
	AddRouter(msgId uint32, router IRouter) // 给当前的服务注册一个Router，供客户端的连接处理使用
	GetConnManager() IConnManager           // 获取当前Server的连接管理模块
	SetOnConnStart(func(conn IConnection))  // 注册OnConnStart钩子函数的方法
	SetOnConnStop(func(conn IConnection))   // 注册OnConnStart钩子函数的方法
	CallOnConnStart(conn IConnection)       // 调用OnConnStart钩子函数的方法
	CallOnConnStop(conn IConnection)        // 调用OnConnStart钩子函数的方法
}
