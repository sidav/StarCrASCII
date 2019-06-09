package main

func createBuilding(codename string, x, y int, f *faction) *pawn {
	var b *pawn
	switch codename {
	// terran
	case "tcommand":
		colors := []int{
			7, 7, 7, 7,
			7, -1, -1, 7,
			7, -1, -1, 7,
			7, 7, 7, 7}
		app := &buildingAppearance{chars: "" +
			"/==\\" +
			"|xx|" +
			"|xx|" +
			"\\--/", colors: colors}
		b = &pawn{name: "Command Center", maxHitpoints: 1000, givesSupply: 10,
			buildingInfo:              &building{w: 4, h: 4, appearance: app, allowsTightPlacement: false},
			currentConstructionStatus: &underConstructionInformation{maxConstructionAmount: 75, costM: 400},
			productionInfo:            &productionInformation{builderCoeff: 1, allowedUnits: []string{"tscv"}},
			res:                       &pawnResourceInformation{receivesResources: true},
		}

	case "tsupply":
		colors := []int{
			-1, 7,
			7, -1}
		app := &buildingAppearance{chars: "" +
			"o=" +
			"=o", colors: colors}
		b = &pawn{name: "Supply Depot", maxHitpoints: 200, givesSupply: 8,
			buildingInfo:              &building{w: 2, h: 2, appearance: app, allowsTightPlacement: true},
			currentConstructionStatus: &underConstructionInformation{maxConstructionAmount: 30, costM: 100},
		}

	case "tbarracks":
		colors := []int{
			-1, 7, -1,
			7, -1, 7,
			-1, 7, -1,
		}
		app := &buildingAppearance{chars: "" +
			"#=#" +
			"=%=" +
			"#=#", colors: colors}
		b = &pawn{name: "Barracks", maxHitpoints: 500,
			buildingInfo:              &building{w: 3, h: 3, appearance: app, allowsTightPlacement: false},
			currentConstructionStatus: &underConstructionInformation{maxConstructionAmount: 60, costM: 100},
			productionInfo:            &productionInformation{builderCoeff: 1, allowedUnits: []string{"tmarine"}},
		}
	case "tbunker":
		colors := []int{
			-1, 7, -1,
			-1, 7, -1,
		}
		app := &buildingAppearance{chars: "" +
			"=%=" +
			"=%=", colors: colors}
		b = &pawn{name: "Bunker", maxHitpoints: 300,
			buildingInfo:              &building{w: 3, h: 2, appearance: app, allowsTightPlacement: true},
			currentConstructionStatus: &underConstructionInformation{maxConstructionAmount: 60, costM: 100},
			containerInfo:             &pawnContainerInformation{allowFireFromInside: true, maxSize: 4},
			//weapons: []*pawnWeaponInformation{{attackDelay: 10, attackRadius: 6, attacksLand: true,
			//	hitscan: &WeaponHitscan{baseDamage: 5, lightMod: 2},
			//}},
		}
	case "tautoturret":
		colors := []int{
			-1, 7,
			7, -1,
		}
		app := &buildingAppearance{chars: "" +
			"T\\" +
			"\\T", colors: colors}
		b = &pawn{name: "Auto turret", maxHitpoints: 150, regenPeriod: 20, sightRadius: 7,
			buildingInfo:              &building{w: 2, h: 2, appearance: app, allowsTightPlacement: true},
			currentConstructionStatus: &underConstructionInformation{maxConstructionAmount: 50, costM: 125},
			weapons: []*pawnWeaponInformation{{attackDelay: 10, attackRadius: 6, attacksLand: true,
				hitscan: &WeaponHitscan{baseDamage: 5, lightMod: 2},
			}},
		}
		// zerg
	case "zhatchery":
		colors := []int{
			7, 7, 7, 7,
			7, -1, -1, 7,
			7, -1, -1, 7,
			7, 7, 7, 7}
		app := &buildingAppearance{chars: "" +
			"/||\\" +
			"=/\\=" +
			"=\\/=" +
			"\\||/", colors: colors}
		b = &pawn{name: "Hatchery", maxHitpoints: 1000, givesSupply: 10,
			buildingInfo:              &building{w: 4, h: 4, appearance: app, allowsTightPlacement: false},
			currentConstructionStatus: &underConstructionInformation{maxConstructionAmount: 75, costM: 300},
			productionInfo:            &productionInformation{builderCoeff: 1, allowedUnits: []string{"zdrone"}},
			res:                       &pawnResourceInformation{receivesResources: true},
		}

		//protoss
	case "pnexus":
		colors := []int{
			7, 7, 7, 7,
			7, -1, -1, 7,
			7, -1, -1, 7,
			7, 7, 7, 7}
		app := &buildingAppearance{chars: "" +
			"\\||/" +
			"-xx-" +
			"-xx-" +
			"/||\\", colors: colors}
		b = &pawn{name: "Nexus", maxHitpoints: 1000, givesSupply: 10,
			buildingInfo:              &building{w: 4, h: 4, appearance: app, allowsTightPlacement: false},
			currentConstructionStatus: &underConstructionInformation{maxConstructionAmount: 75, costM: 400},
			productionInfo:            &productionInformation{builderCoeff: 1, allowedUnits: []string{"pprobe"}},
			res:                       &pawnResourceInformation{receivesResources: true},
		}
	case "pgateway":
		colors := []int{
			-1, 7, -1,
			7, -1, 7,
			-1, 7, -1,
		}
		app := &buildingAppearance{chars: "" +
			"#|#" +
			">*<" +
			"#|#", colors: colors}
		b = &pawn{name: "Gateway", maxHitpoints: 500,
			buildingInfo:              &building{w: 3, h: 3, appearance: app, allowsTightPlacement: false},
			currentConstructionStatus: &underConstructionInformation{maxConstructionAmount: 60, costM: 150},
			productionInfo:            &productionInformation{builderCoeff: 1, allowedUnits: []string{"pzealot"}},
		}
	case "pphotoncannon":
		colors := []int{
			-1, 7,
			7, -1,
		}
		app := &buildingAppearance{chars: "" +
			"P\\" +
			"\\P", colors: colors}
		b = &pawn{name: "Photon Cannon", maxHitpoints: 200, regenPeriod: 20, sightRadius: 7,
			buildingInfo:              &building{w: 2, h: 2, appearance: app, allowsTightPlacement: true},
			currentConstructionStatus: &underConstructionInformation{maxConstructionAmount: 60, costM: 150},
			weapons: []*pawnWeaponInformation{{attackDelay: 12, attackRadius: 6, attacksLand: true,
				hitscan: &WeaponHitscan{baseDamage: 20},
			}},
		}
	// test and whatever
	case "testsmall":
		colors := []int{
			-1,
		}
		app := &buildingAppearance{chars: "#", colors: colors}
		b = &pawn{name: "Test 1x1 Block", maxHitpoints: 200, regenPeriod: 20, sightRadius: 7,
			buildingInfo:              &building{w: 1, h: 1, appearance: app, allowsTightPlacement: true},
			currentConstructionStatus: &underConstructionInformation{maxConstructionAmount: 60, costM: 150},
			//weapons: []*pawnWeaponInformation{{attackDelay: 12, attackRadius: 6, attacksLand: true,
			//	hitscan: &WeaponHitscan{baseDamage: 20},
			//}},
		}
	case "testbig":
		colors := []int{
			-1, -1, -1, -1, -1,
			-1, -1, -1, -1, -1,
			-1, -1, -1, -1, -1,
			-1, -1, -1, -1, -1,
			-1, -1, -1, -1, -1,
		}
		app := &buildingAppearance{chars: "#########################", colors: colors}
		b = &pawn{name: "Test 5x5 Block", maxHitpoints: 200, regenPeriod: 20, sightRadius: 7,
			buildingInfo:              &building{w: 5, h: 5, appearance: app, allowsTightPlacement: true},
			currentConstructionStatus: &underConstructionInformation{maxConstructionAmount: 60, costM: 150},
			//weapons: []*pawnWeaponInformation{{attackDelay: 12, attackRadius: 6, attacksLand: true,
			//	hitscan: &WeaponHitscan{baseDamage: 20},
			//}},
		}
	default:
		colors := []int{
			-1,
		}
		app := &buildingAppearance{chars: "?", colors: colors}
		b = &pawn{name: "Unknown building " + codename, maxHitpoints: 200, regenPeriod: 20, sightRadius: 7,
			buildingInfo:              &building{w: 1, h: 1, appearance: app, allowsTightPlacement: true},
			currentConstructionStatus: &underConstructionInformation{maxConstructionAmount: 60, costM: 150},
			//weapons: []*pawnWeaponInformation{{attackDelay: 12, attackRadius: 6, attacksLand: true,
			//	hitscan: &WeaponHitscan{baseDamage: 20},
			//}},
		}
	}
	if b.maxHitpoints == 0 {
		b.maxHitpoints = 25
	}
	b.hitpoints = b.maxHitpoints
	b.x = x
	b.y = y
	b.faction = f
	b.codename = codename
	if b.sightRadius == 0 {
		b.sightRadius = b.buildingInfo.w + 2
	}
	if b.productionInfo != nil && b.res == nil {
		b.res = &pawnResourceInformation{} // adds zero-value resource info struct for spendings usage.
	}
	return b
}

func getBuildingNameAndDescription(code string) (string, string) {
	bld := createBuilding(code, 0, 0, nil)
	name := bld.name
	var description string
	if bld.currentConstructionStatus != nil {
		constr := bld.currentConstructionStatus
		description += constr.getDescriptionString() + " \\n "
	}
	if len(bld.weapons) > 0 {
		for _, wpn := range bld.weapons {
			description += wpn.getDescriptionString() + " \\n "
		}
	}
	switch code {
	default:
		description += "No description."
	}
	return name, description
}
