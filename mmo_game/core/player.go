package core

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"mmo_game/pb"
	"sync"
	"zinx/ziface"
)

type Player struct {
	PlayerID int32              // 玩家ID
	Conn     ziface.IConnection // 当前玩家的连接（用于和客户端的连接）
	X        float32            // 平面的x坐标
	Y        float32            // 高度
	Z        float32            // 平面的y坐标（注意不是Y）
	V        float32            // 旋转的0-360角度
}

var (
	PidGen int32      = 1 // 用来生成玩家ID的计数器
	Lock   sync.Mutex     // 保护PidGen的锁
)

// NewPlayer 创建一个玩家
func NewPlayer(conn ziface.IConnection) *Player {
	// 生成一个玩家ID
	Lock.Lock()
	id := PidGen
	PidGen++
	Lock.Unlock()

	return &Player{
		PlayerID: id,
		Conn:     conn,
		X:        float32(160 + rand.Intn(10)), // 随机在160坐标点，基于x轴若干偏移
		Y:        0,
		Z:        float32(160 + rand.Intn(10)), // 随机在140坐标点，基于y轴若干偏移
		V:        0,
	}
}

// SendMsg 发送给客户端消息，主要是将pb的protobuf数据序列化之后，再调用zinx的SendMsg方法
func (p *Player) SendMsg(msgId uint32, data proto.Message) {
	// 将proto.Message结构体序列化转换成二进制
	msg, err := proto.Marshal(data)

	if err != nil {
		fmt.Println("Marshal error:", err)
		return
	}

	// 将二进制数据通过zinx框架的SendMsg发送给客户端
	if p.Conn == nil {
		fmt.Println("Connection in player is nil")
		return
	}
	if err := p.Conn.SendMsg(msgId, msg); err != nil {
		fmt.Println("SendMsg error:", err)
		return
	}
}

// SyncPid 将PlayerID同步给客户端
func (p *Player) SyncPid() {
	// 组建MsgID为1的proto数据
	protoMsg := &pb.SyncPid{
		Pid: p.PlayerID,
	}

	// 将消息发送给客户端
	p.SendMsg(1, protoMsg)
}

// BroadCastStartPosition 将Player上线的初始位置同步给客户端
func (p *Player) BroadCastStartPosition() {
	// 组建MsgID为200的proto数据
	protoMsg := &pb.BroadCast{
		Pid: p.PlayerID,
		Tp:  2, // 2代表广播位置坐标
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	// 将消息发送给客户端
	p.SendMsg(200, protoMsg)
}

// Talk 玩家广播世界聊天消息
func (p *Player) Talk(content string) {
	// 1、组建MsgID为200的proto数据
	protoMsg := &pb.BroadCast{
		Pid: p.PlayerID,
		Tp:  1, // 1代表世界聊天
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}

	// 2、得到当前世界所有的在线玩家
	players := WorldMgrObj.GetAllPlayers()

	// 3、向所有的玩家（包括自己）发送MsgID为200的消息
	for _, player := range players {
		// 分别给对应的客户端发送消息
		player.SendMsg(200, protoMsg)
	}
}

// SyncSurrounding 同步玩家上线的位置消息
func (p *Player) SyncSurrounding() {
	// 1、获取当前玩家周围的玩家有哪些（九宫格）
	players := p.GetSurroundPlayers()

	// 2、将当前玩家的位置信息通过MsgID为200的消息发送给周围的玩家（让其他玩家看到自己）
	// 2.1、组建MsgID为200的proto数据
	protoMsg := &pb.BroadCast{
		Pid: p.PlayerID,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	// 2.2、全部周围的玩家都向各自的客户端发送200消息
	for _, player := range players {
		player.SendMsg(200, protoMsg)
	}

	// 3、将周围的全部玩家的位置信息发送给当前的玩家（让自己看到其他玩家）
	// 3.1、组建MsgID为202的proto数据
	// 3.1.1、制作pb.Player切片
	playersProtoMsg := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		// 制作一个pb.Player
		p := &pb.Player{
			Pid: player.PlayerID,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		playersProtoMsg = append(playersProtoMsg, p)
	}
	// 3.1.2、封装MsgID为202的proto数据
	syncPlayersProtoMsg := &pb.SyncPlayers{
		Ps: playersProtoMsg[:],
	}

	// 3.2、将组建好的数据发送给当前玩家的客户端
	p.SendMsg(202, syncPlayersProtoMsg)
}

// UpdatePosition 广播当前玩家的位置移动信息
func (p *Player) UpdatePosition(x, y, z, v float32) {
	// 更新当前玩家的坐标
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v

	// 组建MsgID为200的proto数据（Tp为3）
	protoMsg := &pb.BroadCast{
		Pid: p.PlayerID,
		Tp:  3,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	// 获取当前玩家的周围玩家
	players := p.GetSurroundPlayers()

	// 依次给每个玩家对应的客户端发送当前玩家位置更新的消息
	for _, player := range players {
		player.SendMsg(200, protoMsg)
	}
}

// Offline 玩家下线业务
func (p *Player) Offline() {
	// 得到当前玩家的周围玩家
	players := p.GetSurroundPlayers()

	//给周围玩家广播MsgID的消息
	protoMsg := &pb.SyncPid{
		Pid: p.PlayerID,
	}

	for _, player := range players {
		player.SendMsg(201, protoMsg)
	}

	WorldMgrObj.RemovePlayerByPid(p.PlayerID)
}

// GetSurroundPlayers 获取当前玩家的周围玩家（AOI九宫格内的玩家）
func (p *Player) GetSurroundPlayers() []*Player {
	playerIds := WorldMgrObj.AOIMgr.GetPidByPos(p.X, p.Z)
	players := make([]*Player, 0, len(playerIds))
	for _, playerId := range playerIds {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(playerId)))
	}
	return players
}
