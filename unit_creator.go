package main

// TODO: move unit stats to JSON
func createUnit(codename string, x, y int, f *faction, alreadyConstructed bool) *pawn {
	var newUnit *pawn
	switch codename {
	// terrans
	case "tscv":
		newUnit = &pawn{name: "Terran SCV", maxHitpoints: 200, takesSupply: 1,
			unitInfo:       &unit{appearance: ccell{char: 's'}},
			moveInfo:       &pawnMovementInformation{ticksForMoveSingleCell: 10, movesOnLand: true, movesOnSea: true}, regenPeriod: 7, radarRadius: 0,
			currentConstructionStatus: &underConstructionInformation{maxConstructionAmount: 10, costM: 50},
			productionInfo: &productionInformation{builderCoeff: 1, buildType: buildtype_terran,
				allowedBuildings: []string{"tcommand", "tsupply", "tbarracks"},
			},
			res: &pawnResourceInformation{maxMineralsCarry: 5, ticksToMineMineral: 5},
			weapons: []*pawnWeaponInformation{
				{attackDelay: 10, attackRadius: 1, attacksLand: true,
					hitscan: &WeaponHitscan{baseDamage:4},
				},
			},
		}
	case "tmarine":
		newUnit = &pawn{name: "Terran Marine", maxHitpoints: 40, isLight: true, takesSupply: 1,
			unitInfo:       &unit{appearance: ccell{char: 'm'}},
			moveInfo:       &pawnMovementInformation{ticksForMoveSingleCell: 10, movesOnLand: true}, regenPeriod: 7,
			currentConstructionStatus: &underConstructionInformation{maxConstructionAmount: 10, costM: 50},
			weapons: []*pawnWeaponInformation{
				{attackDelay: 10, attackRadius: 5, attacksLand: true,
					hitscan: &WeaponHitscan{baseDamage:5},
				},
			},
		}
	// zergs
	case "zdrone":
		newUnit = &pawn{name: "Zerg Drone", maxHitpoints: 50, takesSupply: 1,
			unitInfo:       &unit{appearance: ccell{char: 'd'}},
			moveInfo:       &pawnMovementInformation{ticksForMoveSingleCell: 10, movesOnLand: true, movesOnSea: true}, regenPeriod: 7, radarRadius: 0,
			currentConstructionStatus: &underConstructionInformation{maxConstructionAmount: 10, costM: 50},
			productionInfo: &productionInformation{builderCoeff: 1, buildType: buildtype_zerg,
				allowedBuildings: []string{"zhatchery"},
			},
			res: &pawnResourceInformation{maxMineralsCarry: 5, ticksToMineMineral: 5},
			weapons: []*pawnWeaponInformation{
				{attackDelay: 10, attackRadius: 1, attacksLand: true,
					hitscan: &WeaponHitscan{baseDamage:4},
				},
			},
		}
	// protossi
	case "pprobe":
		newUnit = &pawn{name: "Protoss Probe", maxHitpoints: 50, takesSupply: 1,
			unitInfo:       &unit{appearance: ccell{char: 'p'}},
			moveInfo:       &pawnMovementInformation{ticksForMoveSingleCell: 10, movesOnLand: true, movesOnSea: true}, regenPeriod: 7, radarRadius: 0,
			currentConstructionStatus: &underConstructionInformation{maxConstructionAmount: 10, costM: 50},
			productionInfo: &productionInformation{builderCoeff: 1, buildType: buildtype_protoss,
				allowedBuildings: []string{"pnexus"},
			},
			res: &pawnResourceInformation{maxMineralsCarry: 5, ticksToMineMineral: 5},
			weapons: []*pawnWeaponInformation{
				{attackDelay: 10, attackRadius: 1, attacksLand: true,
					hitscan: &WeaponHitscan{baseDamage:4},
				},
			},
		}

	default:
		newUnit = &pawn{name: "UNKNOWN UNIT " + codename,
			moveInfo:                  &pawnMovementInformation{ticksForMoveSingleCell: 10, movesOnLand: true},
			unitInfo:                  &unit{appearance: ccell{char: '?'}},
			currentConstructionStatus: &underConstructionInformation{maxConstructionAmount: 10, costM: 250},
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
	case "tmarine":
		description += "Hey, how to get out this chickenshit outfit?"
	default:
		description += "No description."
	}
	return name, description
}
