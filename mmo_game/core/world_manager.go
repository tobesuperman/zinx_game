package core

import "sync"

// WorldManager 当前游戏的实际管理模块
type WorldManager struct {
	AOIMgr  *AOIManager       // 当前世界地图AOI的管理模块
	Players map[int32]*Player // 当前全部在线的玩家集合
	Lock    sync.RWMutex      // 保护玩家集合的锁
}

// WorldMgrObj 提供一个对外的全局的世界管理模块句柄
var WorldMgrObj *WorldManager

// 初始化方法
func init() {
	WorldMgrObj = &WorldManager{
		// 创建世界地图AOI
		AOIMgr: NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_COUNT_X, AOI_MIN_Y, AOI_MAX_Y, AOI_COUNT_Y),
		// 初始化玩家集合
		Players: make(map[int32]*Player),
	}
}

// AddPlayer 添加一个玩家
func (wm *WorldManager) AddPlayer(player *Player) {
	wm.Lock.Lock()
	wm.Players[player.PlayerID] = player
	wm.Lock.Unlock()

	// 将player添加到AOIManager中
	wm.AOIMgr.AddPidToGridByPos(int(player.PlayerID), player.X, player.Z)
}

// RemovePlayerByPid 删除一个玩家
func (wm *WorldManager) RemovePlayerByPid(playerId int32) {
	if player, ok := wm.Players[playerId]; ok {
		// 将player从AOIManager中删除
		wm.AOIMgr.RemovePidFromGridByPos(int(playerId), player.X, player.Z)
	}
	wm.Lock.Lock()
	delete(wm.Players, playerId)
	wm.Lock.Unlock()

}

// GetPlayerByPid 通过玩家ID查询Player对象
func (wm *WorldManager) GetPlayerByPid(playerId int32) *Player {
	wm.Lock.RLock()
	defer wm.Lock.RUnlock()

	return wm.Players[playerId]
}

// GetAllPlayers 获取全部的在线玩家
func (wm *WorldManager) GetAllPlayers() (players []*Player) {
	wm.Lock.RLock()
	defer wm.Lock.RUnlock()

	players = make([]*Player, 0)
	for _, player := range wm.Players {
		players = append(players, player)
	}
	return
}
