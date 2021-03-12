package ziface

/*
	针对Message进行TLV格式的封包：
	1、先写消息内容的长度和类型
	2、再写消息的内容
*/

/*
	针对Message进行TLV格式的拆包
	1、先读取固定长度的首部，获取消息内容的长度和消息的类型
	2、再根据消息内容的长度进行一次读写，从连接中读取消息的内容
*/

// IDataPack 定义一个解决TCP粘包问题的封包拆包模块
// 直接面向TCP连接的数据流，用于处理TCP粘包问题
type IDataPack interface {
	GetHeadLen() uint32                    // 获取包的head长度
	Pack(message IMessage) ([]byte, error) // 封包
	Unpack(data []byte) (IMessage, error)  // 拆包
}
