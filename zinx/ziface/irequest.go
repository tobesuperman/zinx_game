package ziface

// IRequest 实际上是把客户端请求的连接信息和请求的数据包装起来
type IRequest interface {
	GetConnection() IConnection // 得到当前连接
	GetData() []byte            // 得到请求的消息数据
	GetMsgID() uint32           // 得到请求的消息ID
}
