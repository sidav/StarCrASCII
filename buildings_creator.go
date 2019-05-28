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
			currentConstructionStatus: &constructionInformation{maxConstructionAmount: 100, costM: 400},
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
			currentConstructionStatus: &constructionInformation{maxConstructionAmount: 60, costM: 100},
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
			currentConstructionStatus: &constructionInformation{maxConstructionAmount: 60, costM: 100},
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
			currentConstructionStatus: &constructionInformation{maxConstructionAmount: 99, costM: 300},
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
			currentConstructionStatus: &constructionInformation{maxConstructionAmount: 99, costM: 400},
			productionInfo:            &productionInformation{builderCoeff: 1, allowedUnits: []string{"pprobe"}},
			res:                       &pawnResourceInformation{receivesResources: true},
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
	case "metalmaker":
		description += "A very complicated device which uses quantum fluctuations to convert large energy amounts to some metal."
	case "mextractor":
		description += "A basic ore extraction and purification device. Should be placed on metal deposits."
	case "geo":
		description += "Classic heat-to-electricity conversion device. Should be placed on thermal vents."
	case "mstorage":
		description += "Allows to store more metal."
	case "estorage":
		description += "Allows to store more energy."
	case "quark": // cheating building, useful for debugging
		description += "You should not see that text."
	case "armhq": // cheating building for stub AI
		description += "Enemy HQ. Can gather metal, generate energy and produce Tech 1 land forces."
	case "armvehfactory", "corevehfactory":
		description += "A basic nanolathing facility which is designed to construct wheeled or tracked vehicles. "
	case "armkbotlab", "corekbotlab":
		description += "A basic nanolathing facility which is designed to construct the Kinetic Bio-Organic Technology " +
			"mechs, or KBots. "
	case "coret2kbotlab":
		description += "A more advanced nanolathing facility which is designed to construct Tech 2 KBots."
	case "solar":
		description += "A classic solar battery array. The heavy use of superconductors and wireless energy transfer technologies " +
			"made this energy acqurement devices much more efficient than ever."
	case "lturret":
		description += "A basic yet quite universal base defense structure. Its only weapon uses EM-waves amplified by stimulated emission of radiation."
	case "guardian":
		description += "A stationary plasma artillery with great range and damage, but slow rate of fire."
	case "railgunturret":
		description += "A stationary defense structure which fires projectiles accelerated with Lorenz' force to hypersound velocities. " +
			"Has great range and damage, but slow rate of fire."
	case "radar":
		description += "A radar facility. Reveals enemy units' locations in an area arount itself. Drains energy."
	case "wall":
		description += "A hard metal block designed to block enemy movement."
	default:
		description += "No description."
	}
	return name, description
}
