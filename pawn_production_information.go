package main

type productionInformation struct {
	builderCoeff int
	buildingType int // terran-like, zerg-like or protoss-like
	allowedBuildings []string
	allowedUnits []string
	defaultOrderForUnitBuilt *order
}
