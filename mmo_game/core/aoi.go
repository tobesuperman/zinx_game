package core

import "fmt"

// 定义一些AOI的边界值
const (
	AOI_MIN_X   int = 85
	AOI_MAX_X   int = 410
	AOI_COUNT_X int = 10
	AOI_MIN_Y   int = 75
	AOI_MAX_Y   int = 400
	AOI_COUNT_Y int = 20
)

type AOIManager struct {
	MinX   int           // 区域的左边界坐标
	MaxX   int           //区域的右边界坐标
	CountX int           //X方向格子的数量
	MinY   int           // 区域的上边界坐标
	MaxY   int           // 区域的下边界坐标
	CountY int           // Y方向格子的数量
	grids  map[int]*Grid // 当前区域中有哪些格子map——key=格子的ID，value=格子对象
}

// NewAOIManager 初始化一个AOI管理区域模块
func NewAOIManager(minX, maxX, countX, minY, maxY, countY int) *AOIManager {
	aoiManager := &AOIManager{
		MinX:   minX,
		MaxX:   maxX,
		CountX: countX,
		MinY:   minY,
		MaxY:   maxY,
		CountY: countY,
		grids:  make(map[int]*Grid),
	}

	for y := 0; y < countY; y++ {
		for x := 0; x < countX; x++ {
			// 根据x和y编号计算格子ID
			gridId := y*countX + x

			// 初始化格子
			aoiManager.grids[gridId] = NewGrid(gridId,
				aoiManager.MinX+x*aoiManager.gridWidth(),
				aoiManager.MinX+(x+1)*aoiManager.gridWidth(),
				aoiManager.MinY+y*aoiManager.gridHeight(),
				aoiManager.MinY+(y+1)*aoiManager.gridHeight())
		}
	}
	return aoiManager
}

// 得到每个格子在X轴方向的宽度
func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CountX
}

// 得到每个格子在Y轴方向的高度
func (m *AOIManager) gridHeight() int {
	return (m.MaxY - m.MinY) / m.CountY
}

// 打印格子的基本信息（调试）
func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager:\nMinX: %d, MaxX: %d, CountX: %d, MinY: %d, MaxY: %d, CountY: %d\nGrids in AOIManager:\n",
		m.MinX,
		m.MaxX,
		m.CountX,
		m.MinY,
		m.MaxY,
		m.CountY)
	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}

	return s
}

// GetSurroundGridsByGid 根据格子GID得到周边九宫格格子集合
func (m *AOIManager) GetSurroundGridsByGid(gridId int) (grids []*Grid) {
	// 判断gridId是否在AOIManager中
	if _, ok := m.grids[gridId]; !ok {
		return
	}

	// 初始化grids返回值切片，将当前Grid本身加入到切片中
	grids = append(grids, m.grids[gridId])

	// 当前格子的X轴编号
	idx := gridId % m.CountX

	// 判断左边是否还有格子
	if idx > 0 {
		grids = append(grids, m.grids[gridId-1])
	}

	// 判断右边是否还有格子
	if idx < m.CountX-1 {
		grids = append(grids, m.grids[gridId+1])
	}

	// 将X轴当前的格子都取出进行遍历，再分别得到每个格子上下是否还有格子
	gridsX := make([]int, 0, len(grids))
	for _, v := range grids {
		gridsX = append(gridsX, v.GridID)
	}

	// 遍历gridsX中的每个格子
	for _, v := range gridsX {
		// 当前格子的Y轴编号
		idy := v / m.CountY

		// 判断上边是否还有格子
		if idy > 0 {
			grids = append(grids, m.grids[v-m.CountX])
		}

		// 判断下边是否还有格子
		if idy < m.CountY-1 {
			grids = append(grids, m.grids[v+m.CountX])
		}
	}

	return
}

// GetGidByPos 通过横纵坐标得到当前格子的ID
func (m *AOIManager) GetGidByPos(x, y float32) int {
	idx := (int(x) - m.MinX) / m.gridWidth()
	idy := (int(y) - m.MinY) / m.gridHeight()
	return idy*m.CountX + idx
}

// GetPidByPos 通过横纵坐标得到周边九宫格全部的PlayerIDs
func (m *AOIManager) GetPidByPos(x, y float32) (playerIds []int) {
	// 得到当前玩家的格子ID
	gridId := m.GetGidByPos(x, y)

	// 通过格子ID得到周边九宫格信息
	grids := m.GetSurroundGridsByGid(gridId)

	for _, grid := range grids {
		playerIds = append(playerIds, grid.GetPlayerIDs()...)
	}

	return
}

// AddPidToGrid 添加一个Player到一个格子中
func (m *AOIManager) AddPidToGrid(gridId, playerId int) {
	if _, ok := m.grids[gridId]; !ok {
		return
	}
	m.grids[gridId].Add(playerId)
}

// RemovePidFromGrid 移除一个格子中的Player
func (m *AOIManager) RemovePidFromGrid(gridId, playerId int) {
	if _, ok := m.grids[gridId]; !ok {
		return
	}
	m.grids[gridId].Remove(playerId)
}

// GetPidByGid 通过格子ID获取全部的PlayerID
func (m *AOIManager) GetPidByGid(gridId int) (playerIds []int) {
	if _, ok := m.grids[gridId]; !ok {
		return
	}
	return m.grids[gridId].GetPlayerIDs()
}

// AddPidToGridByPos 通过坐标将Player添加到一个格子中
func (m *AOIManager) AddPidToGridByPos(playerId int, x, y float32) {
	gridId := m.GetGidByPos(x, y)
	m.AddPidToGrid(gridId, playerId)
}

// RemovePidFromGridByPos 通过坐标把一个Player从一个格子中删除
func (m *AOIManager) RemovePidFromGridByPos(playerId int, x, y float32) {
	gridId := m.GetGidByPos(x, y)
	m.RemovePidFromGrid(gridId, playerId)
}
