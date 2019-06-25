package main

import (
	"github.com/sidav/golibrl/astar"
	"github.com/sidav/golibrl/geometry"
)

var (
	mapW int
	mapH int
)

type gameMap struct {
	tileMap  [][]*tile
	factions []*faction
	pawns    []*pawn
}

func (g *gameMap) addPawn(p *pawn) {
	g.pawns = append(g.pawns, p)
}

func (g *gameMap) addBuilding(b *pawn, asAlreadyConstructed bool) {
	if asAlreadyConstructed {
		b.currentConstructionStatus = nil
		b.buildingInfo.hasBeenPlaced = true
	}

	if b.productionInfo != nil && len(b.productionInfo.allowedUnits) > 0 { // sets default rally point for build units.
		b.productionInfo.defaultOrderForUnitBuilt = &order{orderType: order_move, x: b.x + b.buildingInfo.w/2, y: b.y + b.buildingInfo.h + 1}
	}

	g.addPawn(b)
}

func (g *gameMap) removePawn(p *pawn) {
	for i := 0; i < len(g.pawns); i++ {
		if p == g.pawns[i] {
			g.pawns = append(g.pawns[:i], g.pawns[i+1:]...) // ow it's fucking... magic!
		}
	}
}

func (g *gameMap) getPawnAtCoordinates(x, y int) *pawn {
	for _, b := range g.pawns {
		if b.isOccupyingCoords(x, y) {
			return b
		}
	}
	return nil
}

func (g *gameMap) getPawnsInRect(x, y, w, h int) []*pawn {
	var arr []*pawn
	for _, p := range g.pawns {
		cx, cy := p.getCenter()
		if p.isBuilding() {
			if geometry.AreTwoCellRectsOverlapping(x, y, w, h, p.x, p.y, p.buildingInfo.w, p.buildingInfo.h) {
				arr = append(arr, p)
			}
		} else {
			if geometry.AreCoordsInRect(cx, cy, x, y, w, h) {
				arr = append(arr, p)
			}
		}
	}
	return arr
}

//func (g *gameMap) getEnemyPawnsInRadiusFrom(x, y, radius int, f *faction) []*pawn {
//	var arr []*pawn
//	for _, p := range g.pawns {
//		if p.faction != f {
//			if p.isInDistanceFromCoords(x, y, radius) {
//				arr = append(arr, p)
//				continue
//			}
//		}
//	}
//	return arr
//}

func (g *gameMap) getEnemyPawnsInRadiusFromPawn(p *pawn, radius int, f *faction) []*pawn {
	var arr []*pawn
	for _, p2 := range g.pawns {
		if p2.faction != f {
			if p.isInDistanceFromPawn(p2, radius) {
				arr = append(arr, p2)
				continue
			}
		}
	}
	return arr
}

func (g *gameMap) getBuildingAtCoordinates(x, y int) *pawn {
	for _, b := range g.pawns {
		if b.isOccupyingCoords(x, y) {
			return b
		}
	}
	return nil
}

func (g *gameMap) getMineralsAtCoordinates(x, y int) int {
	return g.tileMap[x][y].mineralsAmount
}

func (g *gameMap) getNumberOfMetalDepositsInRect(x, y, w, h int) int {
	total := 0
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			if areCoordsValid(x+i, y+j) {
				total += g.tileMap[x+i][y+j].metalAmount
			}
		}
	}
	return total
}

func (g *gameMap) getNumberOfThermalDepositsInRect(x, y, w, h int) int {
	total := 0
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			if areCoordsValid(x+i, y+j) {
				total += g.tileMap[x+i][y+j].thermalAmount
			}
		}
	}
	return total
}

func (g *gameMap) getNumberOfMetalDepositsUnderBuilding(b *pawn) int {
	return g.getNumberOfMetalDepositsInRect(b.x, b.y, b.buildingInfo.w, b.buildingInfo.h)
}

func (g *gameMap) getNumberOfThermalDepositsUnderBuilding(b *pawn) int {
	return g.getNumberOfThermalDepositsInRect(b.x, b.y, b.buildingInfo.w, b.buildingInfo.h)
}

func (g *gameMap) isPawnInPylonFieldOfFaction(p *pawn, f *faction) bool {
	for _, b := range g.pawns {
		if b.faction == f && b.pylonFieldRadius > 0 {
			if b.isInDistanceFromPawn(p, b.pylonFieldRadius) {
				return true
			}
		}
	}
	return false
}

func (g *gameMap) isConstructionSiteBlockedByUnitOrBuilding(x, y, w, h int, tight bool) bool {
	for _, p := range g.pawns {
		if p.isBuilding() {
			if p.buildingInfo.allowsTightPlacement && tight {
				if geometry.AreTwoCellRectsOverlapping(x, y, w, h, p.x, p.y, p.buildingInfo.w, p.buildingInfo.h) {
					return true
				}
			} else if geometry.AreTwoCellRectsOverlapping(x-1, y-1, w+2, h+2, p.x, p.y, p.buildingInfo.w, p.buildingInfo.h) {
				// -1s and +2s are to prevent tight placement...
				// ..and ensure that there always will be at least 1 cell between buildings.
				return true
			}
		} else {
			cx, cy := p.getCenter()
			if geometry.AreCoordsInRect(cx, cy, x, y, w, h) {
				return true
			}
		}
	}
	return false
}

func (g *gameMap) canBuildingBeBuiltAt(b *pawn, cx, cy int) bool {
	bx := cx - b.buildingInfo.w/2
	by := cy - b.buildingInfo.h/2
	b.x = bx
	b.y = by
	if bx < 0 || by < 0 || bx+b.buildingInfo.w >= mapW || by+b.buildingInfo.h >= mapH {
		return false
	}
	if b.buildingInfo.canBeBuiltInPylonFieldOnly && !g.isPawnInPylonFieldOfFaction(b, b.faction) {
		return false
	}
	if b.buildingInfo.canBeBuiltOnThermalOnly && g.getNumberOfThermalDepositsInRect(bx, by, b.buildingInfo.w, b.buildingInfo.h) == 0 {
		return false
	}
	for x := bx; x < bx+b.buildingInfo.w; x++ {
		for y := by; y < by+b.buildingInfo.h; y++ {
			if !g.tileMap[x][y].isPassable {
				return false
			}
		}
	}
	if g.isConstructionSiteBlockedByUnitOrBuilding(bx, by, b.buildingInfo.w, b.buildingInfo.h, b.buildingInfo.allowsTightPlacement) {
		return false
	}
	return true
}

func (g *gameMap) createCostMapForPathfinding() *[][]int {
	width, height := len(g.tileMap), len((g.tileMap)[0])

	costmap := make([][]int, width)
	for j := range costmap {
		costmap[j] = make([]int, height)
	}
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			// TODO: optimize by iterating through pawns separately
			if !(g.tileMap[i][j].isPassable) || g.getPawnAtCoordinates(i, j) != nil {
				costmap[i][j] = -1
			}
		}
	}
	return &costmap
}

func (g *gameMap) getPathFromTo(fx, fy, tx, ty int) *astar.Cell {
	return astar.FindPath(g.createCostMapForPathfinding(), fx, fy, tx, ty, true, true, false)
}

func (g *gameMap) depleteMineralField(x, y int) {
	currTile := g.tileMap[x][y]
	if currTile.mineralsAmount <= 0 {
		if currTile.mineralsAmount < 0 {
			log.warning("Minerals < 0 !!!")
		}
		currTile.mineralsAmount = 0
		currTile.isPassable = true
		currTile.appearance.char = '.'
	}
}
