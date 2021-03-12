package core

import (
	"fmt"
	"sync"
)

// Grid 一个AOI地图中的格子对象
type Grid struct {
	GridID  int          // 格子ID
	MinX    int          // 格子的左边边界坐标
	MaxX    int          // 格子的右边边界坐标
	MinY    int          // 格子的上边边界坐标
	MaxY    int          // 格子的下边边界坐标
	players map[int]bool // 当前格子内玩家或者物体成员的ID集合
	lock    sync.RWMutex // 保护当前集合的锁
}

// NewGrid 初始化当前格子
func NewGrid(gridId, minX, maxX, miny, maxY int) *Grid {
	return &Grid{
		GridID:  gridId,
		MinX:    minX,
		MaxX:    maxX,
		MinY:    miny,
		MaxY:    maxY,
		players: make(map[int]bool),
	}
}

// Add 给格子添加一个玩家
func (g *Grid) Add(playerId int) {
	g.lock.Lock()
	defer g.lock.Unlock()

	g.players[playerId] = true
}

// Remove 从格子中删除一个玩家
func (g *Grid) Remove(playerId int) {
	g.lock.Lock()
	defer g.lock.Unlock()

	delete(g.players, playerId)
}

// GetPlayerIDs 得到当前格子中所有的玩家ID
func (g *Grid) GetPlayerIDs() (playerIds []int) {
	g.lock.RLock()
	defer g.lock.RUnlock()

	for k, _ := range g.players {
		playerIds = append(playerIds, k)
	}
	return
}

// 打印格子的基本信息（调试）
func (g *Grid) String() string {
	return fmt.Sprintf("GridID: %d, MinX: %d, MaxX: %d, MinY: %d, MaxY: %d, players: %v",
		g.GridID,
		g.MinX,
		g.MaxX,
		g.MinY,
		g.MaxY,
		g.players)
}
