package ziface

// IConnManager 连接管理抽象层
type IConnManager interface {
	Add(conn IConnection)                   // 添加连接
	Remove(conn IConnection)                // 删除连接
	Get(connId uint32) (IConnection, error) // 根据ConnID获取连接
	Len() int                               // 得到当前连接总数
	Clear()                                 // 清除并终止所有连接
}
