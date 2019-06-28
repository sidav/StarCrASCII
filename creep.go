package main

const (
	CREEP_SPREAD_PERIOD         = 150 // in ticks
	CREEP_APPEAR_PROBABILITY    = 2   // more == less likely; 1 == always
	CREEP_DISAPPEAR_PROBABILITY = 3   // same
)

func isTimeToSpreadCreep() bool {
	return CURRENT_TICK%CREEP_SPREAD_PERIOD == 0
}

func (g *gameMap) getNumberOfCreepInRect(x, y, w, h int) int {
	total := 0
	for i := x; i < x+w; i++ {
		for j := y; j < y+h; j++ {
			if areCoordsValid(i, j) && g.tileMap[i][j].creep {
				total++
			}
		}
	}
	return total
}

func (g *gameMap) isCellNeighbouringACreep(x, y int) bool {
	if areCoordsValid(x-1, y) && g.tileMap[x-1][y].creep {
		return true
	}
	if areCoordsValid(x+1, y) && g.tileMap[x+1][y].creep {
		return true
	}
	if areCoordsValid(x, y-1) && g.tileMap[x][y-1].creep {
		return true
	}
	if areCoordsValid(x, y+1) && g.tileMap[x][y+1].creep {
		return true
	}
	return false
}

func (g *gameMap) renewCreepSpread(creepers *[]*pawn) {
	newCreepTiles := make([]*tile, 0)
	removeCreepTiles := make([]*tile, 0)
	for x := 0; x < mapW; x++ {
	nextCell:
		for y := 0; y < mapH; y++ {
			curTile := g.tileMap[x][y]
			isInDistance := false
			for _, creeper := range *creepers {
				isInDistance = creeper.isInDistanceFromCoords(x, y, creeper.creepSpreadRadius) && !creeper.isUnderConstruction()
				if isInDistance && curTile.creep {
					continue nextCell
				}
				if isInDistance && (g.isCellNeighbouringACreep(x, y) || creeper.IsCloseupToCoords(x, y)) && (x+y+CURRENT_TICK/CREEP_SPREAD_PERIOD)%CREEP_APPEAR_PROBABILITY == 0 {
					newCreepTiles = append(newCreepTiles, curTile)
					continue nextCell
				}
			}
			if !isInDistance && curTile.creep && g.getNumberOfCreepInRect(x-1, y-1, 3, 3) <= 5 && (x+y+CURRENT_TICK/CREEP_SPREAD_PERIOD)%CREEP_DISAPPEAR_PROBABILITY == 0 {
				removeCreepTiles = append(removeCreepTiles, curTile)
			}
		}
	}
	for _, c := range newCreepTiles {
		c.creep = true
	}
	for _, c := range removeCreepTiles {
		c.creep = false
	}
}
