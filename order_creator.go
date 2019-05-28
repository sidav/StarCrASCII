package main

func issueDefaultOrderToUnit(p *pawn, m *gameMap, x, y int) {
	if p.isOccupyingCoords(x, y) {
		p.reportOrderCompletion(p.getCurrentOrderDescription() + " order untouched")
		return
	}
	if m.tileMap[x][y].mineralsAmount > 0 && p.canCollectMinerals() {
		p.setOrder(&order{orderType: order_gather_minerals, xSecondary: x, ySecondary: y})
		log.appendMessage(p.name + ": gathering minerals.")
		return
	}
	target := m.getPawnAtCoordinates(x, y)
	if target != nil {
		if target.faction != p.faction {
			p.setOrder(&order{orderType: order_attack, targetPawn: target, x: x, y: y})
			log.appendMessage(p.name + ": attacking.")
			return
		}
		//if target.isBuilding() && target.currentConstructionStatus.isCompleted() == false {
		//	p.setOrder(&order{orderType: order_build, buildingToConstruct: target})
		//	log.appendMessage(p.name + ": Helps nanolathing")
		//	return
		//}
	}
	if p.canMove() {
		if p.faction.cursor.currentCursorMode == CURSOR_AMOVE {
			p.setOrder(&order{orderType: order_attack_move, x: x, y: y})
		} else {
			p.setOrder(&order{orderType: order_move, x: x, y: y})
		}
	} else if p.canConstructUnits() {
		if CURRENT_MAP.getMineralsAtCoordinates(x, y) > 0 {
			p.productionInfo.defaultOrderForUnitBuilt = &order{orderType: order_gather_minerals, x: x, y: y, xSecondary: x, ySecondary: y}
			p.reportOrderCompletion("mine minerals set")
			return
		}
		if p.faction.cursor.currentCursorMode == CURSOR_AMOVE {
			p.productionInfo.defaultOrderForUnitBuilt = &order{orderType: order_attack_move, x: x, y: y}
			p.reportOrderCompletion("default engage location set")
		} else {
			p.productionInfo.defaultOrderForUnitBuilt = &order{orderType: order_move, x: x, y: y}
			p.reportOrderCompletion("rally point set")
		}
	}
}
