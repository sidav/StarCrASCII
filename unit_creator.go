package main

// TODO: move unit stats to JSON
func createUnit(codename string, x, y int, f *faction, alreadyConstructed bool) *pawn {
	var newUnit *pawn
	switch codename {
	case "tscv":
		newUnit = &pawn{name: "Terran SCV", maxHitpoints: 200,
			unitInfo:       &unit{appearance: ccell{char: 's'}},
			moveInfo:       &pawnMovementInformation{ticksForMoveSingleCell: 10, movesOnLand: true, movesOnSea: true}, regenPeriod: 7, radarRadius: 0,
			currentConstructionStatus: &constructionInformation{maxConstructionAmount: 10, costM: 50},
			// res:            &pawnResourceInformation{metalIncome: 2, energyIncome: 30, metalStorage: 250, energyStorage: 1000},
			productionInfo: &productionInformation{builderCoeff: 1, allowedBuildings: []string{
				"tcommand"},
			},
			res: &pawnResourceInformation{maxMineralsCarry: 5, ticksToMineMineral: 5},
			weapons: []*pawnWeaponInformation{
				{attackDelay: 10, attackEnergyCost: 1, attackRadius: 1, attacksLand: true,
					hitscan: &WeaponHitscan{baseDamage:4},
				},
			},
		}
	case "tmarine":
		newUnit = &pawn{name: "Terran Marine", maxHitpoints: 200, isHeavy: true, isCommander: true,
			unitInfo:       &unit{appearance: ccell{char: 'm'}},
			moveInfo:       &pawnMovementInformation{ticksForMoveSingleCell: 10, movesOnLand: true}, regenPeriod: 7,
			weapons: []*pawnWeaponInformation{
				{attackDelay: 10, attackEnergyCost: 1, attackRadius: 5, attacksLand: true,
					hitscan: &WeaponHitscan{baseDamage:5},
				},
			},
		}
	default:
		newUnit = &pawn{name: "UNKNOWN UNIT " + codename,
			moveInfo:                  &pawnMovementInformation{ticksForMoveSingleCell: 10, movesOnLand: true},
			unitInfo:                  &unit{appearance: ccell{char: '?'}},
			currentConstructionStatus: &constructionInformation{maxConstructionAmount: 10, costM: 250},
		}
	}
	if newUnit.maxHitpoints == 0 {
		newUnit.maxHitpoints = 1
		log.appendMessage("No hitpoints set for "+newUnit.name)
	}
	newUnit.hitpoints = newUnit.maxHitpoints
	newUnit.x = x
	newUnit.y = y
	newUnit.faction = f
	newUnit.codename = codename
	if newUnit.sightRadius == 0 {
		newUnit.sightRadius = 8
	}
	if alreadyConstructed {
		newUnit.currentConstructionStatus = nil
	}
	return newUnit
}

func getUnitNameAndDescription(code string) (string, string) {
	unit := createUnit(code, 0, 0, nil, false)
	name := unit.name
	var description string
	if unit.currentConstructionStatus != nil {
		constr := unit.currentConstructionStatus
		description += constr.getDescriptionString() + " \\n "
	}
	if len(unit.weapons) > 0 {
		for _, wpn := range unit.weapons {
			description += wpn.getDescriptionString() + " \\n "
		}
	}
	switch code {
	case "tscv":
		description += "Space Construction Vehicle."
	default:
		description += "No description."
	}
	return name, description
}
