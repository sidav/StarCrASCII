package main

type tile struct {
	appearance                          *ccell
	isNaval, isPassable, isMineralField bool
	mineralsAmount                      int
	vespeneAmount                       int
	metalAmount                         int
	thermalAmount                       int
}
