package main

// TODO: move unit stats to JSON
func createUnit(codename string, x, y int, f *faction, alreadyConstructed bool) *pawn {
	var newUnit *pawn
	switch codename {
	case "tscv":
		newUnit = &pawn{name: "Terran SCV", maxHitpoints: 200,
			unitInfo:       &unit{appearance: ccell{char: 's'}},
			moveInfo:       &pawnMovementInformation{ticksForMoveSingleCell: 10, movesOnLand: true, movesOnSea: true}, regenPeriod: 7, radarRadius: 15,
			currentConstructionStatus: &constructionInformation{maxConstructionAmount: 10, costM: 50},
			// res:            &pawnResourceInformation{metalIncome: 2, energyIncome: 30, metalStorage: 250, energyStorage: 1000},
			productionInfo: &productionInformation{builderCoeff: 1, allowedBuildings: []string{
				"tcommand"},
			},
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
	case "armcommander":
		description += "A Commander Unit of the Arm Rebellion."
	case "corecommander":
		description += "A Commander Unit of the Core Corporation."
	case "protocommander":
		description += "The brain and heart of any modern military operation, Command Unit is a massive bipedal amphibious hulk " +
			"equipped with anything needed for forward operations base establishment and support - a radar, production, " +
			"quantum generator, metal synthesizer, and light weaponry. " +
			"Although the Command Unit is extremely heavily armored and is capable of self-repair, it have to be well protected, " +
			"because its loss means inevitable defeat. \\n " +
			"This old prototype lacks some more modern equipment, such as the Disintegrator Gun."
	case "coreck":
		description += "An engineering KBot equipped with production. Can build more advanced buildings than Commander do."
	case "corecv":
		description += "An engineering vehicle equipped with production. Can build more advanced buildings than Commander do."
	case "coreweasel":
		description += "Fast recon vehicle. It has very weak attack, but is equipped with advanced visual sensors array " +
			"which is providing quite huge vision range."
	case "armjeffy":
		description += "Fast recon vehicle. It has very weak attack, but is equipped with advanced visual sensors array " +
			"which is providing quite huge vision range."
	case "coreraider":
		description += "Fast light tank. Its speed and armor regeneration ability makes it useful for hit-and-run tactics."
	case "armflash":
		description +=  "Fast light tank. Its speed and armor regeneration ability makes it useful for hit-and-run tactics."
	case "coreak":
		description += "A basic assault KBot effective against light armor."
	case "armpeewee":
		description += "A cheap and relatively fast basic assault KBot effective against light armor."
	case "armhammer":
		description += "A basic artillery KBot. Effective against heavy armor. Designed to take out buildings. "
	case "corethud":
		description += "A basic artillery KBot. Effective against heavy armor. Designed to take out buildings. "
	case "corethecan":
		description += "Slow and clunky, The Can is designed to take part in front-line assault. Although its " +
			"armor can sustain significant amount of punishment, this KBot should be supported due to its short range."
	default:
		description += "No description."
	}
	return name, description
}
