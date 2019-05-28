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
			buildingInfo:              &building{w: 4, h: 4, appearance: app, allowsTightPlacement: true},
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
			buildingInfo:              &building{w: 3, h: 3, appearance: app, allowsTightPlacement: true},
			currentConstructionStatus: &underConstructionInformation{maxConstructionAmount: 60, costM: 100},
			productionInfo:            &productionInformation{builderCoeff: 1, allowedUnits: []string{"tmarine"}},
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
			buildingInfo:              &building{w: 4, h: 4, appearance: app, allowsTightPlacement: true},
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
			buildingInfo:              &building{w: 4, h: 4, appearance: app, allowsTightPlacement: true},
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
			buildingInfo:              &building{w: 3, h: 3, appearance: app, allowsTightPlacement: true},
			currentConstructionStatus: &underConstructionInformation{maxConstructionAmount: 60, costM: 150},
			productionInfo:            &productionInformation{builderCoeff: 1, allowedUnits: []string{"pzealot"}},
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
