package main

type tile struct {
	appearance          *ccell
	isNaval, isPassable bool
	mineralsAmount      int
	vespeneAmount       int
	metalAmount         int
	thermalAmount       int
}
