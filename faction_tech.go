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

func (ftt *factionTech) areRequirementsSatisfiedForCode(code string) bool {
	reqs := PSI_getRequiredTechBuilding(code)
	if reqs == nil {
		return true
	}
	for _, req := range *reqs {
		if !ftt.doesHaveTechBuilding(req) {
			return false
		}
	}
	return true
}
