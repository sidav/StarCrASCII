package main

type factionTech struct {
	techBuildings []string // contains list of buildings the faction currently has
}

func (ftt *factionTech) doesHaveTechBuilding(code string) bool {
	if code == "" {
		return true
	}
	for i := range ftt.techBuildings {
		if ftt.techBuildings[i] == code {
			return true
		}
	}
	return false
}

func (ftt *factionTech) addTechBuildingToList(code string) {
	if !ftt.doesHaveTechBuilding(code) {
		ftt.techBuildings = append(ftt.techBuildings, code)
	}
}

func (f *faction) checkTechBuildings() {
	f.tech.techBuildings = f.tech.techBuildings[:0] // clears the array keeping the memory allocated
	for _, u := range CURRENT_MAP.pawns {
		if u.faction == f && u.isBuilding() && !u.isUnderConstruction() {
			f.tech.addTechBuildingToList(u.codename)
		}
	}
}
