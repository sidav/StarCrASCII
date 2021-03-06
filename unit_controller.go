package main

import (
	"github.com/sidav/golibrl/geometry"
)

func (p *pawn) isTimeToAct() bool {
	return p.nextTickToAct <= CURRENT_TICK
}

func (u *pawn) executeOrders(m *gameMap) {
	if !u.isTimeToAct() || u.isDisabled {
		return
	}

	order := u.order
	if order == nil {
		return
	}

	switch order.orderType {
	case order_move:
		u.doMoveOrder(0)
	case order_attack:
		u.doAttackOrder()
	case order_attack_move:
		u.doAttackMoveOrder()
	case order_build:
		u.doBuildOrder(m)
	case order_construct:
		u.doConstructOrder(m)
	case order_gather_minerals:
		u.doGatherMineralsOrder()
	case order_gather_vespene:
		u.doGatherVespeneOrder()
	case order_return_resources:
		u.doReturnResourcesOrder()
	case order_enter_container:
		u.doEnterContainerOrder()
	case order_unload:
		u.doUnloadOrder()
	}

	// move

}

func (u *pawn) doMoveOrder(desiredAccuracy int) bool { // Returns true if route exists. TODO: rewrite
	order := u.order

	ox, oy := order.x, order.y
	ux, uy := u.getCoords()
	var vx, vy int

	//vector := geometry.CreateVectorByStartAndEndInt(ux, uy, ox, oy)
	//vector.TransformIntoUnitVector()
	//vx, vy := vector.GetRoundedCoords()
	path := CURRENT_MAP.getPathFromTo(ux, uy, ox, oy, desiredAccuracy)
	if path != nil {
		vx, vy = path.GetNextStepVector()
	}

	if vx == 0 && vy == 0 && (ux != ox || uy != oy) { // path stops not at the target
		u.reportOrderCompletion("Can't find route to target. Arrived to closest position.") // can be dangerous if order is not move
		if order.orderType == order_move || order.orderType == order_attack_move {
			u.order = nil
		}
		return false
	}

	if u.coll_canMoveByVector(vx, vy) {

		u.x += vx
		u.y += vy

		u.nextTickToAct = CURRENT_TICK + u.moveInfo.ticksForMoveSingleCell

		if u.x == ox && u.y == oy {
			if order.orderType == order_move || order.orderType == order_attack_move {
				u.reportOrderCompletion("Arrived")
				u.order = nil
			}
			return true
		}
	}
	return true
}

func (p *pawn) switchToAnotherMineralFieldNearby() { // find pseudorandom mineral field close to the ordered one
	order := p.order
	mx, my := order.xSecondary, order.ySecondary
	const MINERAL_SEARCH_AREA = 4
	mineralsLocsx := make([]int, 0)
	mineralsLocsy := make([]int, 0)
	for i := mx - MINERAL_SEARCH_AREA; i <= mx+MINERAL_SEARCH_AREA; i++ {
		for j := my - MINERAL_SEARCH_AREA; j <= my+MINERAL_SEARCH_AREA; j++ {
			if areCoordsValid(i, j) && CURRENT_MAP.getMineralsAtCoordinates(i, j) > 0 {
				mineralsLocsx = append(mineralsLocsx, i)
				mineralsLocsy = append(mineralsLocsy, j)
			}
		}
	}
	if len(mineralsLocsx) == 0 {
		p.reportOrderCompletion("No minerals nearby. Going on standby.")
		p.order = nil
		return
	} else {
		num := CURRENT_TICK % len(mineralsLocsx)
		order.xSecondary = mineralsLocsx[num]
		order.ySecondary = mineralsLocsy[num]
	}
}

func (p *pawn) doAttackOrder() { // Only moves the unit to a firing position. The firing itself is in openFireIfPossible()
	order := p.order

	ux, uy := p.getCoords()
	if order.targetPawn.hitpoints <= 0 {
		p.reportOrderCompletion("target destroyed. Now standing by")
		p.order = nil
		return
	}
	targetX, targetY := order.targetPawn.getCenter()

	if !geometry.AreCoordsInRange(ux, uy, targetX, targetY, p.getMaxRadiusToFire()) {
		order.x = targetX
		order.y = targetY
		p.doMoveOrder(0)
		return
	}
}

func (p *pawn) doEnterContainerOrder() {
	order := p.order
	// get contaier pawn
	cont := order.targetPawn
	if cont.isInDistanceFromPawn(p, 1) {
		if cont.hitpoints <= 0 {
			p.reportOrderCompletion("container destroyed. Now standing by")
			p.order = nil
			return
		}
		if !cont.containerInfo.canAddPawn(p) || cont.isUnderConstruction() {
			p.reportOrderCompletion("can't enter. Now  standing by")
			p.order = nil
			return
		}

		cont.containerInfo.addPawnToContainer(p)
		CURRENT_MAP.removePawn(p)
		p.reportOrderCompletion("entered.")
		p.order = nil
		return

	} else {
		order.x, order.y = cont.getCenter()
		p.doMoveOrder(0)
		return
	}
}

func (attacker *pawn) openFireIfPossible() { // does the firing, does NOT necessary mean execution of attack order (but can be)
	if attacker.currentConstructionStatus != nil ||
		attacker.order != nil && attacker.order.orderType == order_build ||
		!attacker.hasWeapons() && !(attacker.canContainPawns() && attacker.containerInfo.allowFireFromInside) ||
		attacker.isDisabled {
		return
	}
	var pawnInOrder *pawn
	if attacker.order != nil && attacker.order.targetPawn != nil {
		pawnInOrder = attacker.order.targetPawn
	}

	for _, wpn := range attacker.weapons {
		fired := attackWithWeapon(wpn, attacker, pawnInOrder)
		if fired {
			if wpn.canBeFiredOnMove {
				wpn.nextTurnToFire = CURRENT_TICK + wpn.attackDelay
			} else {
				attacker.nextTickToAct = CURRENT_TICK + wpn.attackDelay
			}
		}
	}

	// attack from inside of a bunker
	if attacker.canContainPawns() && attacker.containerInfo.allowFireFromInside {
		for _, unitInside := range attacker.containerInfo.pawnsInside {
			if unitInside.isTimeToAct() {
				for _, wpn := range unitInside.weapons {
					fired := attackWithWeapon(wpn, attacker, pawnInOrder)
					if fired {
						if wpn.canBeFiredOnMove {
							wpn.nextTurnToFire = CURRENT_TICK + wpn.attackDelay
						} else {
							unitInside.nextTickToAct = CURRENT_TICK + wpn.attackDelay
						}
					}
				}
			}
		}
	}
}

func (p *pawn) doAttackMoveOrder() {
	if p.isTimeToAct() {
		p.openFireIfPossible()
	}
	if p.isTimeToAct() {
		p.doMoveOrder(0)
	}
}

func (p *pawn) doUnloadOrder() {
	if len(p.containerInfo.pawnsInside) > 0 {
		for i, curr_unit := range p.containerInfo.pawnsInside {
			curr_unit.x = p.x + i
			curr_unit.y = p.y - 1
			curr_unit.order = nil
			CURRENT_MAP.addPawn(curr_unit)
		}
		p.containerInfo.pawnsInside = nil
	}
	p.order = nil
}

func (u *pawn) doBuildOrder(m *gameMap) { // only moves to location and/or sets the spendings. Building itself is in doAllProduction()
	// TODO: rewrite the heck out of it. Tip: implement and use doCircleAndRectangleIntersect() with the build radius
	order := u.order
	tBld := order.buildingToConstruct
	ux, uy := u.getCoords()

	if tBld == nil {
		log.appendMessage(u.name + " NIL BUILD")
		return
	}

	if tBld.currentConstructionStatus == nil {
		u.reportOrderCompletion("Construction is finished by another unit")
		u.order = nil
		return
	}

	if u.productionInfo.buildType == buildtype_zerg { // zerg will always move to center
		// out of range, move to the construction site
		order.x, order.y = tBld.getCenter()
		u.doMoveOrder(0)
		log.appendMessage(u.name + " MOVES TO BUILD")
		return
	}

	if !tBld.IsCloseupToCoords(ux, uy) {
		// out of range, move to the construction site
		order.x, order.y = tBld.getCenter()
		u.doMoveOrder(0)
		log.appendMessage(u.name + " MOVES TO BUILD")
		return
	}
}

func (p *pawn) doConstructOrder(m *gameMap) {
	if len(p.order.constructingQueue) == 0 {
		p.reportOrderCompletion("Construction queue finished")
		p.order = nil
		return
	}
}

func (u *pawn) reportOrderCompletion(verb string) {
	if u.faction.playerControlled {
		log.appendMessage(u.name + ": " + verb + ".")
	}
}
