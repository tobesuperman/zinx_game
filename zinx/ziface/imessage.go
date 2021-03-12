package ziface

// IMessage 定义抽象接口，将请求的消息封装到一个Message中
type IMessage interface {
	GetMsgID() uint32          // 获取消息的ID
	GetDataLen() uint32        // 获取消息的长度
	GetData() []byte           // 获取消息的内容
	SetMsgID(id uint32)        // 设置消息的ID
	SetDataLen(dataLen uint32) // 设置消息的内容
	SetData(data []byte)       // 设置消息的长度
}
