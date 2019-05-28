package main

func doAllProduction(m *gameMap) { // does the building itself
	for _, u := range m.pawns {
		// buildings self-construction (zerg and protoss)
		if u.currentConstructionStatus != nil {
			if u.currentConstructionStatus.buildType == buildtype_protoss || u.currentConstructionStatus.buildType == buildtype_zerg {
				u.currentConstructionStatus.currentConstructionAmount += 1
				u.hitpoints += u.maxHitpoints / (u.currentConstructionStatus.maxConstructionAmount)
				if u.hitpoints > u.maxHitpoints {
					u.hitpoints = u.maxHitpoints
				}
				if u.currentConstructionStatus.isCompleted() {
					u.currentConstructionStatus = nil
					u.reportOrderCompletion("Construction completed")
				}
			}
		}
		// buildings construction by units
		if u.order != nil && u.order.orderType == order_build {
			tBld := u.order.buildingToConstruct

			ux, uy := u.getCoords()
			bcx, bcy := tBld.getCenter()

			isInBuildProximity := tBld.buildingInfo.hasBeenPlaced == false &&
				(u.productionInfo.buildType == buildtype_zerg && ux == bcx && uy == bcy || // zerg should build from center
				u.productionInfo.buildType != buildtype_zerg && tBld.IsCloseupToCoords(ux, uy))

			if isInBuildProximity { // place the carcass
				if u.faction.canAffordSpend(tBld.currentConstructionStatus.costM, tBld.currentConstructionStatus.costV) {
					u.reportOrderCompletion("Starts construction")
					tBld.buildingInfo.hasBeenPlaced = true
					tBld.currentConstructionStatus.buildType = u.productionInfo.buildType
					tBld.hitpoints = tBld.maxHitpoints % (tBld.currentConstructionStatus.maxConstructionAmount)
					if tBld.hitpoints == 0 {
						tBld.hitpoints = 1 
					}
					u.faction.spendResources(tBld.currentConstructionStatus.costM, tBld.currentConstructionStatus.costV)
					m.addBuilding(tBld, false)
					if u.productionInfo.buildType == buildtype_protoss {
						u.order = nil
						u.reportOrderCompletion("Warp-in initiated.")
					}
					if u.productionInfo.buildType == buildtype_zerg {
						CURRENT_MAP.removePawn(u)
					}
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

			if u.productionInfo.buildType == buildtype_terran && tBld.IsCloseupToCoords(ux, uy) {
				tBld.currentConstructionStatus.currentConstructionAmount += u.productionInfo.builderCoeff
				tBld.hitpoints += tBld.maxHitpoints / (tBld.currentConstructionStatus.maxConstructionAmount/ u.productionInfo.builderCoeff)
				if tBld.hitpoints > tBld.maxHitpoints {
					tBld.hitpoints = tBld.maxHitpoints
				}
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
				if u.faction.economy.cursupply + uCnst.takesSupply > u.faction.economy.maxsupply {
					u.reportOrderCompletion("Additional supply depots required!")
					continue
				}
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
