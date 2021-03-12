package apis

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"mmo_game/core"
	"mmo_game/pb"
	"zinx/ziface"
	"zinx/znet"
)

// MoveApi 玩家移动业务
type MoveApi struct {
	znet.BaseRouter
}

// Handle 玩家移动业务的具体逻辑
func (m *MoveApi) Handle(request ziface.IRequest) {
	// 1、解析客户端传递进来的proto协议
	protoMsg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), protoMsg)
	if err != nil {
		fmt.Println("Unmarshal error:", err)
		return
	}

	// 2、当前的位置信息是属于哪个玩家发起的
	playerId, err := request.GetConnection().GetProperty("playerId")
	if err != nil {
		fmt.Println("GetProperty playerId error:", err)
		return
	}

	fmt.Printf("PlayerID = %d move(%f, %f, %f, %f)\n",
		playerId,
		protoMsg.X,
		protoMsg.Y,
		protoMsg.Z,
		protoMsg.V,
	)

	// 3、根据playerId得到对应的Player对象
	player := core.WorldMgrObj.GetPlayerByPid(playerId.(int32))

	// 4、将这个位置信息广播给其他全部在线的玩家
	player.UpdatePosition(protoMsg.X, protoMsg.Y, protoMsg.Z, protoMsg.V)
}
