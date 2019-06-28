package main

const CREEP_SPREAD_PERIOD = 100 // in ticks

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
	// noCreepersNearby := true
	totalCreepers := len(*creepers)
	for x := 0; x < mapW; x++ {
	nextCell:
		for y := 0; y < mapH; y++ {
			if totalCreepers > 0 {
				for _, creeper := range *creepers {
					if creeper.isInDistanceFromCoords(x, y, creeper.creepSpreadRadius) && !creeper.isUnderConstruction() &&
						(g.isCellNeighbouringACreep(x, y) || creeper.IsCloseupToCoords(x, y)) {
						if (x+y+CURRENT_TICK/CREEP_SPREAD_PERIOD)%3 == 0 {
							g.tileMap[x][y].creep = true
						}
						continue nextCell
					}
				}
			}
			if  g.tileMap[x][y].creep && g.getNumberOfCreepInRect(x-1, y-1, 3, 3) <= 5 && (x+y+CURRENT_TICK/CREEP_SPREAD_PERIOD) % 3 == 0 {
				g.tileMap[x][y].creep = false
			}
		}
	}
}
//
//func (g *gameMap) spreadCreep(creeper *pawn) bool {
//	for r := 1; r <= creeper.creepSpreadRadius; r++ {
//
//	}
//}
