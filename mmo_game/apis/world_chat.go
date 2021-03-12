package apis

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"mmo_game/core"
	"mmo_game/pb"
	"zinx/ziface"
	"zinx/znet"
)

// WorldChatApi 世界聊天业务
type WorldChatApi struct {
	znet.BaseRouter
}

// Handle 世界聊天业务的具体逻辑
func (wc *WorldChatApi) Handle(request ziface.IRequest) {
	// 1、解析客户端传递进来的proto协议
	protoMsg := &pb.Talk{}
	err := proto.Unmarshal(request.GetData(), protoMsg)
	if err != nil {
		fmt.Println("Unmarshal error:", err)
		return
	}

	// 2、当前的聊天数据是属于哪个玩家发起的
	playerId, err := request.GetConnection().GetProperty("playerId")
	if err != nil {
		fmt.Println("GetProperty playerId error:", err)
		return
	}

	// 3、根据playerId得到对应的Player对象
	player := core.WorldMgrObj.GetPlayerByPid(playerId.(int32))

	// 4、将这个消息广播给其他全部在线的玩家
	player.Talk(protoMsg.Content)
}
