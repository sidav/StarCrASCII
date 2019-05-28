package main

type BUILD_TYPE_ENUM uint8

const (
	buildtype_terran  BUILD_TYPE_ENUM = iota
	buildtype_zerg    BUILD_TYPE_ENUM = iota
	buildtype_protoss BUILD_TYPE_ENUM = iota
)

type productionInformation struct {
	builderCoeff             int
	buildType                BUILD_TYPE_ENUM // terran-like, zerg-like or protoss-like
	allowedBuildings         []string
	allowedUnits             []string
	defaultOrderForUnitBuilt *order
}
