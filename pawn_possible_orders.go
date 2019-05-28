package main

func (p *pawn) canConstructBuildings() bool {
	return p.productionInfo != nil && len(p.productionInfo.allowedBuildings) > 0
}

func (p *pawn) canConstructUnits() bool {
	return p.productionInfo != nil && len(p.productionInfo.allowedUnits) > 0
}

func (p *pawn) canMove() bool {
	return p.moveInfo != nil
}
