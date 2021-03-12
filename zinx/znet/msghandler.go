package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

// MsgHandler 消息处理模块的实现
type MsgHandler struct {
	APIs           map[uint32]ziface.IRouter // 存放每个MsgID所对应的处理方法
	TaskQueue      []chan ziface.IRequest    // 负责Worker取任务的消息队列
	WorkerPoolSize uint32                    // 业务工作Worker池中的Worker数量
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		APIs:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, // 从全局配置中获取
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

func (m *MsgHandler) DoMsgHandle(request ziface.IRequest) {
	// 1、从Request中找到MsgID
	handler, ok := m.APIs[request.GetMsgID()]
	if !ok {
		fmt.Println("API NOT FOUND! MsgID =", request.GetMsgID())
		return
	}
	// 2、根据MsgID调度对应的Router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (m *MsgHandler) AddRouter(msgId uint32, router ziface.IRouter) {
	// 1、判断当前msg绑定的API处理方法是否已经存在
	if _, ok := m.APIs[msgId]; ok {
		panic("Repeat API, MsgID = " + strconv.Itoa(int(msgId)))
	}
	// 2、添加msg与API的绑定关系
	m.APIs[msgId] = router
	fmt.Println("Add API success, MsgID =", msgId)
}

// StartWorkerPool 启动一个Worker工作池（开启工作池的动作只能发生一次，一个Zinx框架只能有一个Worker工作池）
func (m *MsgHandler) StartWorkerPool() {
	// 根据WorkerPoolSize分别开启Worker，每个Worker用一个Goroutine来承载
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		// 1、当前的Worker对应的channel消息队列开辟空间
		m.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerPoolSize)

		// 2、启动当前的Worker，阻塞等待消息从channel中传递进来
		go m.startOneWorker(i, m.TaskQueue[i])
	}

}

// startOneWorker 启动一个Worker工作流程
func (m *MsgHandler) startOneWorker(workerId int, taskQueue chan ziface.IRequest) {
	fmt.Println("WorkerID =", workerId, "is starting...")

	// 不断阻塞等待对应消息队列的消息
	for {
		select {
		// 如果有消息过来，出队的就是一个客户端Request，执行当前Request所绑定业务
		case request := <-taskQueue:
			m.DoMsgHandle(request)
		}
	}
}

// SendMsgToTaskQueue 将消息交给TaskQueue，由Worker进行处理
func (m *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	// 1、将消息平均分配给不同的Worker
	// 根据客户端建立的ConnID来进行分配
	workerId := request.GetConnection().GetConnID() % m.WorkerPoolSize
	fmt.Println("Add ConnID =", request.GetConnection().GetConnID(),
		"message MsgID =", request.GetMsgID(),
		"to WorkerID =", workerId)

	// 2、将消息发送给对应的Worker的TaskQueue即可
	m.TaskQueue[workerId] <- request
}
