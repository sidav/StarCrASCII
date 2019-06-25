package main

type pawnResourceInformation struct {
	mineralsCarry, maxMineralsCarry, ticksToMineMineral int
	vespeneCarry, maxVespeneCarry                       int
	receivesResources                                   bool // command centre, whatever
}
