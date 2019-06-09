package main

import (
	cw "github.com/sidav/golibrl/console"
	"github.com/sidav/golibrl/geometry"
)

func (p *pawn) isTimeToAct() bool {
	return p.nextTickToAct <= CURRENT_TICK
}

func (u *pawn) executeOrders(m *gameMap) {
	if !u.isTimeToAct() {
		return
	}

	order := u.order
	if order == nil {
		return
	}

	switch order.orderType {
	case order_move:
		u.doMoveOrder()
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
	case order_return_resources:
		u.doReturnResourcesOrder()
	case order_enter_container:
		u.doEnterContainerOrder()
	}

	// move

}

func (u *pawn) doMoveOrder() bool { // Returns true if route exists. TODO: rewrite
	order := u.order

	ox, oy := order.x, order.y
	ux, uy := u.getCoords()
	var vx, vy int

	//vector := geometry.CreateVectorByStartAndEndInt(ux, uy, ox, oy)
	//vector.TransformIntoUnitVector()
	//vx, vy := vector.GetRoundedCoords()
	path := CURRENT_MAP.getPathFromTo(ux, uy, ox, oy)
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

func (p *pawn) doGatherMineralsOrder() {
	order := p.order

	ux, uy := p.getCoords()
	mx, my := order.xSecondary, order.ySecondary

	if p.res == nil {
		log.warning("Unit " + p.name + " is trying to gather minerals! Whaaaat the heeeeck?")
		p.order.orderType = order_move
		return
	}

	if p.res.mineralsCarry == 0 {
		mins := CURRENT_MAP.getMineralsAtCoordinates(mx, my)
		if mins <= 0 {
			p.switchToAnotherMineralFieldNearby()
			return
		}
		if !geometry.AreCoordsInRange(ux, uy, mx, my, 1) {
			order.x = mx
			order.y = my
			pathSuccess := p.doMoveOrder()
			if !pathSuccess {
				p.switchToAnotherMineralFieldNearby()
			}
			return
		}
		if mins > p.res.maxMineralsCarry {
			CURRENT_MAP.tileMap[mx][my].mineralsAmount -= p.res.maxMineralsCarry
			p.res.mineralsCarry = p.res.maxMineralsCarry
			p.nextTickToAct = CURRENT_TICK + p.res.ticksToMineMineral*p.res.maxMineralsCarry
		} else {
			CURRENT_MAP.tileMap[mx][my].mineralsAmount = 0
			p.res.mineralsCarry = mins
			p.nextTickToAct = CURRENT_TICK + mins*p.res.maxMineralsCarry
			CURRENT_MAP.depleteMineralField(mx, my)
			log.appendMessage("Mineral field depleted.")
		}
	}
	order.orderType = order_return_resources
}

func (p *pawn) doReturnResourcesOrder() {
	order := p.order
	ux, uy := p.getCoords()
	closestResourceReceiver := p.faction.getResourceReceiverNearCoords(ux, uy)
	if closestResourceReceiver == nil {
		p.reportOrderCompletion("Nowhere to return the resources.")
		p.nextTickToAct = CURRENT_TICK + 10
		return
	} else {
		if closestResourceReceiver.IsCloseupToCoords(p.x, p.y) { // resources unload
			if p.res.mineralsCarry > 0 {
				p.faction.economy.minerals += p.res.mineralsCarry
				p.res.mineralsCarry = 0
				p.nextTickToAct = CURRENT_TICK + 10
				order.orderType = order_gather_minerals
			}
			if p.res.vespeneCarry > 0 {
				p.faction.economy.vespene += p.res.vespeneCarry
				p.res.vespeneCarry = 0
				p.nextTickToAct = CURRENT_TICK + 10
				// order.orderType = order_gather_minerals // TODO: order_gather_vespene
			}
		}
		order.x, order.y = closestResourceReceiver.getCenter()
		p.doMoveOrder()
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
		p.doMoveOrder()
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
		if !cont.isUnderConstruction() {
			cont.containerInfo.addPawnToContainer(p)
			CURRENT_MAP.removePawn(p)
			p.reportOrderCompletion("entered the container")
			p.order = nil
			return
		}
	} else {
		order.x, order.y = cont.getCenter()
		p.doMoveOrder()
		return
	}
}

func (attacker *pawn) openFireIfPossible() { // does the firing, does NOT necessary mean execution of attack order (but can be)
	if attacker.currentConstructionStatus != nil || !attacker.hasWeapons() || attacker.order != nil && attacker.order.orderType == order_build {
		return
	}
	var pawnInOrder *pawn
	if attacker.order != nil && attacker.order.targetPawn != nil {
		pawnInOrder = attacker.order.targetPawn
	}
	attackerCenterX, attackerCenterY := attacker.getCenter()
	for _, wpn := range attacker.weapons {
		//if attacker.faction.economy.currentEnergy < wpn.attackEnergyCost {
		//	continue
		//}
		if (wpn.canBeFiredOnMove && wpn.nextTurnToFire > CURRENT_TICK) || (!wpn.canBeFiredOnMove && !attacker.isTimeToAct()) {
			// log.appendMessage(fmt.Sprintf("Skipping fire: TtA:%b CBFoM:%b TRN: %b", attacker.isTimeToAct() ,wpn.canBeFiredOnMove, wpn.nextTurnToFire > CURRENT_TICK))
			continue
		}
		var target *pawn
		radius := wpn.attackRadius
		if pawnInOrder != nil && attacker.isInDistanceFromPawn(pawnInOrder, radius) {
			target = pawnInOrder
		} else {
			potential_targets := CURRENT_MAP.getEnemyPawnsInRadiusFromPawn(attacker, radius, attacker.faction)
			for _, potentialTarget := range potential_targets {
				ptx, pty := potentialTarget.getCoords()
				if attacker.faction.areCoordsInSight(ptx, pty) || attacker.faction.areCoordsInRadarRadius(ptx, pty) {
					target = potentialTarget
				}
			}
		}
		if target != nil {
			if wpn.canBeFiredOnMove {
				wpn.nextTurnToFire = CURRENT_TICK + wpn.attackDelay
			} else {
				attacker.nextTickToAct = CURRENT_TICK + wpn.attackDelay
			}
			// draw the pew pew laser TODO: move this crap somewhere already
			if areGlobalCoordsOnScreenForFaction(attackerCenterX, attackerCenterY, CURRENT_FACTION_SEEING_THE_SCREEN) || areGlobalCoordsOnScreenForFaction(target.x, target.y, CURRENT_FACTION_SEEING_THE_SCREEN) {
				cw.SetFgColor(cw.RED)
				cx, cy := target.getCenter()
				camx, camy := CURRENT_FACTION_SEEING_THE_SCREEN.cursor.getCameraCoords()
				renderLine(attackerCenterX, attackerCenterY, cx, cy, false, camx, camy)
				FIRE_WAS_OPENED_ON_SCREEN_THIS_TURN = true
			}
			dealDamageToTarget(attacker, wpn, target)
		}
	}
}

func (p *pawn) doAttackMoveOrder() {
	if p.isTimeToAct() {
		p.openFireIfPossible()
	}
	if p.isTimeToAct() {
		p.doMoveOrder()
	}
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
		u.doMoveOrder()
		log.appendMessage(u.name + " MOVES TO BUILD")
		return
	}

	if !tBld.IsCloseupToCoords(ux, uy) {
		// out of range, move to the construction site
		order.x, order.y = tBld.getCenter()
		u.doMoveOrder()
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
	//
	//uCnst := order.constructingQueue[0]
	//
	//p.res.metalSpending = p.productionInfo.builderCoeff * uCnst.currentConstructionStatus.costM / uCnst.currentConstructionStatus.maxConstructionAmount
}

func (u *pawn) reportOrderCompletion(verb string) {
	if u.faction.playerControlled {
		log.appendMessage(u.name + ": " + verb + ".")
	}
}
