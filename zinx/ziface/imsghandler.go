package ziface

// IMsgHandler 消息处理模块的抽象接口
type IMsgHandler interface {
	DoMsgHandle(request IRequest)           // 调度/执行对应的Router消息处理方法
	AddRouter(msgId uint32, router IRouter) // 为消息添加具体的处理逻辑
	StartWorkerPool()                       // 启动Worker工作池
	SendMsgToTaskQueue(request IRequest)    // 将消息发送给消息任务队列处理
}
