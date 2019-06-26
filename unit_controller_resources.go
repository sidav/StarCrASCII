package main

import "github.com/sidav/golibrl/geometry"

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
			pathSuccess := p.doMoveOrder(0)
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

func (p *pawn) doGatherVespeneOrder() {
	order := p.order

	// ux, uy := p.getCoords()

	if p.res == nil {
		log.warning("Unit " + p.name + " is trying to gather vespene! Whaaaat the heeeeck?")
		p.order.orderType = order_move
		return
	}
	workerX, workerY := p.getCoords()
	gasMine := order.targetPawn
	mx, my := gasMine.getCenter()
	if p.res.vespeneCarry == 0 {
		gas := CURRENT_MAP.getVespeneAtCoordinates(mx, my)
		if !gasMine.IsCloseupToCoords(workerX, workerY) {
			rx, ry := gasMine.getCoords()
			w, h := gasMine.buildingInfo.w, gasMine.buildingInfo.h
			order.x, order.y = geometry.GetCellNearestToRectFrom(rx, ry, w, h, workerX, workerY)
			pathSuccess := p.doMoveOrder(0)
			if !pathSuccess {
				order.x, order.y = mx, my
				p.doMoveOrder(0)
			}
			return
		}
		if gasMine.isTimeToAct() && !gasMine.isUnderConstruction() {
			if gas > p.res.maxVespeneCarry {
				p.res.vespeneCarry = p.res.maxVespeneCarry
				p.nextTickToAct = CURRENT_TICK + p.res.ticksToMineMineral*p.res.maxVespeneCarry
			} else {
				p.res.vespeneCarry = gas
				p.nextTickToAct = CURRENT_TICK + gas*p.res.maxVespeneCarry
				log.appendMessage("Mineral field depleted.")
			}
			CURRENT_MAP.decreaseVespeneUnderMine(gasMine, p.res.maxVespeneCarry)
			gasMine.nextTickToAct = CURRENT_TICK + 30
			order.orderType = order_return_resources
		}
	}
	p.nextTickToAct = CURRENT_TICK + 5
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
		if closestResourceReceiver.IsCloseupToCoords(ux, uy) { // resources unload
			if p.res.mineralsCarry > 0 {
				p.faction.economy.minerals += p.res.mineralsCarry
				p.res.mineralsCarry = 0
				p.nextTickToAct = CURRENT_TICK + 25
				order.orderType = order_gather_minerals
			}
			if p.res.vespeneCarry > 0 {
				p.faction.economy.vespene += p.res.vespeneCarry
				p.res.vespeneCarry = 0
				p.nextTickToAct = CURRENT_TICK + 25
				order.orderType = order_gather_vespene
			}
		}
		rx, ry := closestResourceReceiver.getCoords()
		w, h := 1, 1
		if closestResourceReceiver.isBuilding() {
			w, h = closestResourceReceiver.buildingInfo.w, closestResourceReceiver.buildingInfo.h
		}
		order.x, order.y = geometry.GetCellNearestToRectFrom(rx, ry, w, h, ux, uy)
		pathSuccess := p.doMoveOrder(0)
		if !pathSuccess {
			order.x, order.y = closestResourceReceiver.getCenter()
			p.doMoveOrder(0)
		}
	}
}
