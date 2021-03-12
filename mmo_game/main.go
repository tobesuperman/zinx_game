package main

import (
	"fmt"
	"mmo_game/apis"
	"mmo_game/core"
	"zinx/ziface"
	"zinx/znet"
)

// OnConnectionStart 连接建立之后调用的Hook函数
func OnConnectionStart(conn ziface.IConnection) {
	// 创建一个Player对象
	player := core.NewPlayer(conn)

	// 给客户端发送MsgID为1的消息：同步当前玩家的ID给客户端
	player.SyncPid()

	// 给客户端发送MsgID为200的消息：同步当前玩家的初始位置坐标给客户端
	player.BroadCastStartPosition()

	// 将当前新上线的玩家添加到WorldManager中
	core.WorldMgrObj.AddPlayer(player)

	// 将该连接绑定一个玩家ID
	conn.SetProperty("playerId", player.PlayerID)

	// 同步周围玩家，告知他们当前玩家已经上线，广播当前玩家的位置信息
	player.SyncSurrounding()

	fmt.Println("PlayerID =", player.PlayerID, "has arrived")
}

// OnConnectionStop 连接断开之前调用的Hook函数
func OnConnectionStop(conn ziface.IConnection) {
	// 通过连接属性得到当前连接所绑定玩家ID
	playerId, _ := conn.GetProperty("playerId")
	player := core.WorldMgrObj.GetPlayerByPid(playerId.(int32))

	// 触发玩家下线的业务
	player.Offline()

	fmt.Println("PlayerID =", playerId, "offline...")
}

func main() {
	// 创建zinx Server实例
	s := znet.NewServer()

	// 连接创建和销毁的HOOK钩子函数
	s.SetOnConnStart(OnConnectionStart)
	s.SetOnConnStop(OnConnectionStop)

	// 注册一些路由业务
	s.AddRouter(2, &apis.WorldChatApi{})
	s.AddRouter(3, &apis.MoveApi{})

	// 启动服务
	s.Serve()
}
