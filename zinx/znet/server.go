package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

// Server IServer的接口实现，定义一个Server的服务器模块
type Server struct {
	Name        string                        // 服务器的名称
	IPVersion   string                        // 服务器绑定的IP版本
	IP          string                        // 服务器监听的IP
	Port        int                           // 服务器监听的端口
	MsgHandler  ziface.IMsgHandler            // 当前Server的消息管理模块，用来绑定MsgID和对应的处理业务API关系
	ConnManager ziface.IConnManager           // 当前Server的连接管理模块
	OnConnStart func(conn ziface.IConnection) // 当前Server创建连接之后自动调用的Hook函数
	OnConnStop  func(conn ziface.IConnection) // 当前Server创建连接之后自动调用的Hook函数
}

// NewServer 初始化Server模块
func NewServer() ziface.IServer {
	s := &Server{
		Name:        utils.GlobalObject.Name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalObject.IP,
		Port:        8999,
		MsgHandler:  NewMsgHandler(),
		ConnManager: NewConnManager(),
	}
	return s
}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name: %s, Listener at IP: %s, Port: %d, is starting\n", s.Name, s.IP, s.Port)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d, MaxPackageSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPackageSize,
	)

	go func() {
		// 0、开启消息队列及Worker工作池
		s.MsgHandler.StartWorkerPool()

		// 1、获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("Resolve TCP addr error:", err)
			return
		}

		// 2、监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("Listen", s.IPVersion, " error:", err)
			return
		}
		fmt.Println("Start Zinx server", s.Name, "success")
		var connID uint32
		connID = 0

		// 3、阻塞地等待客户端连接，处理客户端连接业务（读写）
		for {
			// 如果有客户端连接，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error:", err)
				continue
			}

			// 设置最大连接个数的判断，如果超过最大连接，则关闭此新的连接
			if s.ConnManager.Len() >= utils.GlobalObject.MaxConn {
				// TODO 给客户端响应一个超过最大连接的错误包
				fmt.Println("Too many connections, MaxConn =", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}
			// 将处理新连接的业务方法和conn进行绑定，得到连接模块
			dealConn := NewConnection(s, conn, connID, s.MsgHandler)
			connID++

			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	// TODO 将一些服务器的资源、状态或者已经开辟的连接信息，进行停止或回收
	s.ConnManager.Clear()
	fmt.Println("Stop Zinx server", s.Name, "success")
}

func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()

	// TODO 做一些启动服务器之外的额外业务

	// 阻塞状态
	select {}
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
	fmt.Println("Add router success")
}

func (s *Server) GetConnManager() ziface.IConnManager {
	return s.ConnManager
}

func (s *Server) SetOnConnStart(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("Call OnConnStart()...")
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("Call OnConnStop()...")
		s.OnConnStop(conn)
	}
}
