package main

import "fmt"

type underConstructionInformation struct {
	// for buildings which are under construction right now
	currentConstructionAmount, maxConstructionAmount int
	buildType                                        BUILD_TYPE_ENUM
	costM, costV                                     int
}

func (ci *underConstructionInformation) isCompleted() bool {
	return ci.currentConstructionAmount >= ci.maxConstructionAmount
}

func (ci *underConstructionInformation) getCompletionPercent() int {
	return ci.currentConstructionAmount * 100 / ci.maxConstructionAmount
}

func (ci *underConstructionInformation) getDescriptionString() string {
	return fmt.Sprintf("Metal: %d ENERGY: %d Base build time: %d", ci.costM, 50, ci.maxConstructionAmount)
}
