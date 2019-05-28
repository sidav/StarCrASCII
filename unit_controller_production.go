package main

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
