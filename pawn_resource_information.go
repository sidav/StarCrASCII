package main

type pawnResourceInformation struct {
	mineralsCarry, maxMineralsCarry int
	vespeneCarry, maxVespeneCarry   int
	// TODO: delete fields below
	energyDrain                   int  //
	metalSpending, energySpending int  // both are unconditional only
	isMetalExtractor              bool // gives metal only when placed at metal deposit
	isGeothermalPowerplant        bool // gives energy only when placed at thermal vents
}
