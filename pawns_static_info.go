package main

// Everything which is not required to be personally stored in building structs.

func PSI_getRequiredTechBuilding(code string) *[]string {
	switch code {
	// terran
	case "tbarracks":
		return &([]string {"tsupply"})
	case "tautoturret":
		return &([]string {"tbarracks"})
	//case "tmarine":
	//	return &([]string {"tbarracks"})
	// zerg
	// protoss
	default:
		return nil
	}
}

func PSI_getPylonFieldRadius(code string) int {
	switch code {
	// terran
	case "ppylon":
		return 5
	default:
		return 0
	}
}

