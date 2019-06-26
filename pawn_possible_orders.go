package main

func (p *pawn) canConstructBuildings() bool {
	return p.productionInfo != nil && len(p.productionInfo.allowedBuildings) > 0
}

func (p *pawn) canConstructUnits() bool {
	return p.productionInfo != nil && len(p.productionInfo.allowedUnits) > 0
}

func (p *pawn) canCollectMinerals() bool {
	return p.res != nil && p.res.maxMineralsCarry > 0
}

func (p *pawn) canCollectVespene() bool {
	return p.res != nil && p.res.maxVespeneCarry > 0
}

func (p *pawn) canMove() bool {
	return p.moveInfo != nil
}

func (p *pawn) canReleaseContainedPawns() bool {
	return p.canContainPawns() && len(p.containerInfo.pawnsInside) > 0
}

