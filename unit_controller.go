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
	}

	// move

}

func (u *pawn) doMoveOrder() { // TODO: rewrite
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
		u.order = nil
		return
	}

	if u.coll_canMoveByVector(vx, vy) {

		u.x += vx
		u.y += vy

		u.nextTickToAct = CURRENT_TICK + u.moveInfo.ticksForMoveSingleCell

		if u.x == ox && u.y == oy {
			u.reportOrderCompletion("Arrived")
			u.order = nil
			return
		}
	}
}

func (p *pawn) doGatherMineralsOrder() {
	order := p.order

	ux, uy := p.getCoords()
	mx, my := order.xSecondary, order.ySecondary

	if p.res.mineralsCarry == 0 {
		if !geometry.AreCoordsInRange(ux, uy, mx, my, 1) {
			order.x = mx
			order.y = my
			p.doMoveOrder()
			return
		}
		mins := CURRENT_MAP.getMineralsAtCoordinates(mx, my)
		if mins > 0 {
			if mins > p.res.maxMineralsCarry {
				CURRENT_MAP.tileMap[mx][my].mineralsAmount -= p.res.maxMineralsCarry
				p.res.mineralsCarry = p.res.maxMineralsCarry
				p.nextTickToAct = CURRENT_TICK + p.res.ticksToMineMineral*p.res.maxMineralsCarry
			} else {
				CURRENT_MAP.tileMap[mx][my].mineralsAmount = 0
				p.res.mineralsCarry = mins
				p.nextTickToAct = CURRENT_TICK + mins*p.res.maxMineralsCarry
			}
		}
	} else { // return to command center or whatever
		var closestResourceReceiver *pawn
		closestRRDist := 999999
		for _, cc := range CURRENT_MAP.pawns {
			if cc.res != nil && cc.res.receivesResources {
				rx, ry := cc.getCenter()
				dist := (rx-ux)*(rx-ux) + (ry-uy)*(ry-uy)
				if dist < closestRRDist {
					closestRRDist = dist
					closestResourceReceiver = cc
				}
			}
		}
		if closestResourceReceiver == nil {
			p.reportOrderCompletion("Nowhere to return the resources.")
			p.nextTickToAct = CURRENT_TICK + 10
			return
		} else {
			if closestRRDist < 16 { // resources unload
				p.faction.economy.minerals += p.res.mineralsCarry
				p.res.mineralsCarry = 0
				p.nextTickToAct = CURRENT_TICK + 10
			}
			order.x, order.y = closestResourceReceiver.getCenter()
			p.doMoveOrder()
		}
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
		if pawnInOrder != nil && geometry.AreCoordsInRange(attackerCenterX, attackerCenterY, pawnInOrder.x, pawnInOrder.y, radius) {
			target = pawnInOrder
		} else {
			potential_targets := CURRENT_MAP.getEnemyPawnsInRadiusFrom(attackerCenterX, attackerCenterY, radius, attacker.faction)
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
	ox, oy := tBld.getCenter()

	building_w := tBld.buildingInfo.w + 1
	building_h := tBld.buildingInfo.h + 1
	sqdistance := (ox-ux)*(ox-ux) + (oy-uy)*(oy-uy)

	if tBld == nil {
		log.appendMessage(u.name + " NIL BUILD")
		return
	}

	if tBld.currentConstructionStatus == nil {
		u.reportOrderCompletion("Construction is finished by another unit")
		u.order = nil
		return
	}

	if sqdistance < building_w*building_w || sqdistance < building_h*building_h { // is in building range
		// u.res.metalSpending = u.productionInfo.builderCoeff * tBld.currentConstructionStatus.costM / tBld.currentConstructionStatus.maxConstructionAmount
	} else { // out of range, move to the construction site
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

func doAllProduction(m *gameMap) { // does the building itself
	for _, u := range m.pawns {
		// buildings construction
		if u.order != nil && u.order.orderType == order_build {
			tBld := u.order.buildingToConstruct

			ux, uy := u.getCoords()
			ox, oy := tBld.getCenter()
			building_w := tBld.buildingInfo.w + 1
			building_h := tBld.buildingInfo.h + 1
			sqdistance := (ox-ux)*(ox-ux) + (oy-uy)*(oy-uy)

			if tBld.buildingInfo.hasBeenPlaced == false && (sqdistance <= building_w*building_w || sqdistance <= building_h*building_h) { // place the carcass
				if u.faction.canAffordSpend(tBld.currentConstructionStatus.costM, tBld.currentConstructionStatus.costV) {
					u.reportOrderCompletion("Starts construction")
					tBld.buildingInfo.hasBeenPlaced = true
					u.faction.spendResources(tBld.currentConstructionStatus.costM, tBld.currentConstructionStatus.costV)
					m.addBuilding(tBld, false)
				} else {
					u.reportOrderCompletion("Awaiting resources...")
					continue
				}
			}

			if tBld.currentConstructionStatus == nil {
				u.reportOrderCompletion("Construction interrupted")
				u.order = nil
				continue
			}
			if tBld.hitpoints <= 0 {
				u.reportOrderCompletion("Construction interrupted by hostile action")
				u.order = nil
				continue
			}

			if sqdistance < building_w*building_w || sqdistance < building_h*building_h {
				tBld.currentConstructionStatus.currentConstructionAmount += u.productionInfo.builderCoeff
				if tBld.currentConstructionStatus.isCompleted() {
					tBld.currentConstructionStatus = nil
					u.order = nil
					u.reportOrderCompletion("Construction completed")
				}
			}
		}

		// units construction
		if u.order != nil && u.order.orderType == order_construct {
			uCnst := u.order.constructingQueue[0]

			ux, _ := u.getCenter()

			if uCnst.currentConstructionStatus == nil {
				u.reportOrderCompletion("WTF CONSTRUCTION STATUS IS NIL FOR " + uCnst.name)
				continue
			}
			if uCnst.currentConstructionStatus.currentConstructionAmount == 0 {
				if u.faction.canAffordSpend(uCnst.currentConstructionStatus.costM, uCnst.currentConstructionStatus.costV) {
					u.faction.spendResources(uCnst.currentConstructionStatus.costM, uCnst.currentConstructionStatus.costV)
				} else {
					u.reportOrderCompletion("Awaiting resources...")
					continue
				}
			}
			uCnst.currentConstructionStatus.currentConstructionAmount += u.productionInfo.builderCoeff
			if uCnst.currentConstructionStatus.isCompleted() {
				uCnst.currentConstructionStatus = nil
				uCnst.x, uCnst.y = ux, u.y+u.buildingInfo.h
				uCnst.order = &order{}
				uCnst.order.cloneFrom(u.productionInfo.defaultOrderForUnitBuilt)
				m.addPawn(uCnst)
				u.order.constructingQueue = u.order.constructingQueue[1:]
				if u.repeatConstructionQueue {
					u.order.constructingQueue = append(u.order.constructingQueue, createUnit(uCnst.codename, 0, 0, u.faction, false))
				}
				u.reportOrderCompletion("Construction completed")
			}

		}
	}
}

func (u *pawn) reportOrderCompletion(verb string) {
	if u.faction.playerControlled {
		log.appendMessage(u.name + ": " + verb + ".")
	}
}
