package core

import (
	"fmt"
	"testing"
)

func TestNewAOIManager(t *testing.T) {
	aoiManager := NewAOIManager(100, 300, 4, 200, 450, 5)

	fmt.Println(aoiManager)
}

func TestAOIManager_GetSurroundGridsByGid(t *testing.T) {
	aoiManager := NewAOIManager(0, 250, 5, 0, 250, 5)

	for gridId, _ := range aoiManager.grids {
		grids := aoiManager.GetSurroundGridsByGid(gridId)
		gIDs := make([]int, 0, len(grids))
		for _, grid := range grids {
			gIDs = append(gIDs, grid.GridID)
		}
		fmt.Println("gridId:", gridId, "len:", len(grids))
		fmt.Println(gIDs)
	}
}
