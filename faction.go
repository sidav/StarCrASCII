package main

func getFactionRGB(fn int) (uint8, uint8, uint8) {
	switch fn {
	case 0:
		return 0, 0, 255
	case 1:
		return 255, 0, 0
	case 2:
		return 0, 255, 0
	case 3:
		return 255, 255, 0
	}
	return 32, 32, 32
}

type faction struct {
	cursor                                 *cursor // cursor position
	economy                                *factionEconomy
	tech                                   *factionTech
	factionNumber                          int
	name                                   string
	playerControlled, aiControlled         bool    // used as a stub for now
	aiData                                 *aiData // for AI-controlled factions
	seenTiles, tilesInSight, radarCoverage [][] bool
}

func createFaction(name string, n int, playerControlled, aiControlled bool) *faction { // temporary
	fctn := &faction{
		playerControlled: playerControlled, aiControlled: aiControlled, name: name, factionNumber: n,
		economy: &factionEconomy{minerals: 50, vespene: 0}, cursor: &cursor{},
		tech: &factionTech{},
	}
	if aiControlled {
		fctn.aiData = ai_createAiData()
	}
	fctn.seenTiles = make([][]bool, mapW)
	for i := range fctn.seenTiles {
		fctn.seenTiles[i] = make([]bool, mapH)
	}
	fctn.tilesInSight = make([][]bool, mapW)
	for i := range fctn.tilesInSight {
		fctn.tilesInSight[i] = make([]bool, mapH)
	}
	fctn.radarCoverage = make([][]bool, mapW)
	for i := range fctn.radarCoverage {
		fctn.radarCoverage[i] = make([]bool, mapH)
	}
	return fctn
}

func (f *faction) getResourceReceiverNearCoords(x, y int) *pawn {
	var closestResourceReceiver *pawn
	closestRRDist := 999999
	for _, cc := range CURRENT_MAP.pawns { // TODO: optimize this, do not iterate over every pawn every turn.
		if cc.res != nil && cc.faction == f && cc.currentConstructionStatus == nil && cc.res.receivesResources {
			rx, ry := cc.getCenter()
			dist := (rx-x)*(rx-x) + (ry-y)*(ry-y)
			if dist < closestRRDist {
				closestRRDist = dist
				closestResourceReceiver = cc
			}
		}
	}
	return closestResourceReceiver
}

func (f *faction) recalculateSupplyAndTech() {
	f.economy.cursupply = 0
	f.economy.maxsupply = 0
	f.tech.techBuildings = f.tech.techBuildings[:0] // clears the array keeping the memory allocated

	for _, p := range CURRENT_MAP.pawns {
		if p.faction == f {
			if p.isBuilding() && p.buildingInfo.canBeBuiltInPylonFieldOnly {
				p.isDisabled = !CURRENT_MAP.isPawnInPylonFieldOfFaction(p, f)
				p.order = nil
			}
			if !p.isUnderConstruction() {
				f.economy.maxsupply += p.givesSupply
				if p.isBuilding() {
					f.tech.addTechBuildingToList(p.codename)
				}
			}
			f.economy.cursupply += p.takesSupply
		}
	}
}

func (f *faction) getFactionColor() int {
	//BLACK        = 0
	//DARK_RED     = 1
	//DARK_GREEN   = 2
	//DARK_YELLOW  = 3
	//DARK_BLUE    = 4
	//DARK_MAGENTA = 5
	//DARK_CYAN    = 6
	//BEIGE        = 7
	//DARK_GRAY    = 8
	//RED          = 9
	//GREEN        = 10
	//YELLOW       = 11
	//BLUE         = 12
	//MAGENTA      = 13
	//CYAN         = 14
	//WHITE        = 15
	switch f.factionNumber {
	case 0:
		return 14
	case 1:
		return 9
	case 2:
		return 10
	case 3:
		return 11
	}
	return 7
}
