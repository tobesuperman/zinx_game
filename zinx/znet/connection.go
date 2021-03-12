package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	Server         ziface.IServer         // 当前Connection隶属于哪个Server
	Conn           *net.TCPConn           // 当前连接的socket TCP套接字
	ConnID         uint32                 // 当前连接的ID
	isClosed       bool                   // 当前连接的状态
	ExitChan       chan bool              // 告知当前连接已经退出（停止）的channel（由Reader告知Writer退出）
	msgChan        chan []byte            // 无缓冲通道，用户读写goroutine之间的消息通信
	MsgHandler     ziface.IMsgHandler     // 消息管理模块
	properties     map[string]interface{} // 连接属性集合
	propertiesLock sync.RWMutex           // 保护连接属性的锁
}

// NewConnection 初始化链接模块
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, MsgHandler ziface.IMsgHandler) *Connection {
	connection := &Connection{
		Server:     server,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: MsgHandler,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
		properties: make(map[string]interface{}),
	}

	// 将conn加入到ConnManager中
	connection.Server.GetConnManager().Add(connection)
	return connection
}

// StartReader 连接的读数据业务方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader goroutine is running]")
	defer fmt.Println("ConnID =", c.ConnID, "RemoteAddr =", c.Conn.RemoteAddr().String(), "reader exit...")
	defer c.Stop()

	for {
		dp := NewDataPack()

		// 读取客户端的msgHead，二进制流，8字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.Conn, headData); err != nil {
			fmt.Println("Read msg head error:", err)
			break
		}

		// 拆包，得到msgID和msgDataLen，存放在消息对象msg中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("Server unpack error:", err)
			break
		}

		// 根据msgDataLen，存放在msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			// 第二次读消息内容
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.Conn, data); err != nil {
				fmt.Println("Server unpack error:", err)
				break
			}
		}
		msg.SetData(data)

		// 得到当前Conn的Request
		req := &Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启了工作池，将消息发送给Worker工作池处理即可
			c.MsgHandler.SendMsgToTaskQueue(req)
		} else {
			// 从路由中，找到注册绑定的Conn对应的Router调用
			go c.MsgHandler.DoMsgHandle(req)
		}
	}
}

// StartWriter 专门将数据发送给客户端
func (c *Connection) StartWriter() {
	fmt.Println("[Writer goroutine is running]")
	defer fmt.Println("ConnID =", c.ConnID, "RemoteAddr =", c.Conn.RemoteAddr().String(), "writer exit...")

	// 不断地阻塞地等待channel的数据，如果有数据则发送给客户端
	for {
		select {
		// 有数据要发送给客户端
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error:", err)
				return
			}
		// 代表Reader已经退出，说明Writer也要退出
		case <-c.ExitChan:
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("ConnID =", c.ConnID, "start...")

	// 启动从当前连接读数据的业务
	go c.StartReader()
	// 启动从当前连接写数据的业务
	go c.StartWriter()

	// 按照开发者传递进来的创建连接之后需要调用的处理业务，执行对应的Hook函数
	c.Server.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	fmt.Println("ConnID =", c.ConnID, "stop...")

	// 如果当前连接已经关闭
	if c.isClosed {
		return
	}
	c.isClosed = true

	// 按照开发者传递进来的销毁连接之前需要调用的处理业务，执行对应的Hook函数
	c.Server.CallOnConnStop(c)

	// 关闭socket连接
	c.Conn.Close()

	// 告知Writer关闭
	c.ExitChan <- true

	// 将当前连接从ConnManager中删除
	c.Server.GetConnManager().Remove(c)

	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SendMsg 将要发送给客户端的数据先进行封包，再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed when sending msg")
	}

	// 进行封包
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("Pack ID =", msgId, "error")
		return errors.New("pack msg error")
	}

	// 将数据发送给客户端
	c.msgChan <- binaryMsg
	return nil
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertiesLock.Lock()
	defer c.propertiesLock.Unlock()

	c.properties[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertiesLock.RLock()
	defer c.propertiesLock.RUnlock()

	if value, ok := c.properties[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("property NOT FOUND")
	}
}

func (c *Connection) RemoveProperty(key string) {
	c.propertiesLock.Lock()
	defer c.propertiesLock.Unlock()

	delete(c.properties, key)
}
