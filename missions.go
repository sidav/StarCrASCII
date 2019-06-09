package main

import (
	rnd "github.com/sidav/golibrl/random"
	"strconv"
)

//var MIS1_MAP = &[]string {
//	"~~~.........................^............................................",
//	"~~..;;......................^............................................",
//	"~~..;..............;;.......^^......;;;..................................",
//	"~.....;.....................^^^.....;;;;.................................",
//	"~~..........................^^^.....;;;..................................",
//	"~~~..........................^...........................................",
//	"~~.......................................................................",
//	"~~.......................................................................",
//	"~.....................................^..................................",
//	"~..................;............^^^.^^^^^................................",
//	"~..........$........;;......^^^^^.^^^^..^^^^^^...........................",
//	"~..........$........;.......^^..........^................................",
//	"~~...........................^^^.....$...................................",
//	"~..............................^....$$$..................................",
//	"~~..................................$$...................................",
//	"~~~.........................^^...........................................",
//	"~~~.........;................^^..........................................",
//	"~~~~.......;;;................^^.........................................",
//	"~~~~........................^^^..........................................",
//	"~~~~~.....................^^^^...........................................",
//}

var MIS1_MAP = &[]string {
	"........................................^^^^^^^^^^^^^^^^^^^^........................................",
	"....................................................^^^^^^^^........................................",
	"....................................................^^^^^^^^........................................",
	"........................................^^^^^^^^^...^^^^^^^^........................................",
	"........................................^^^^^^^^^...^^^^^^^^........................................",
	"..........*.....................***.....^^^^^^^^^...^^^^^^^^..........................*.............",
	".........***....................***.....^^^^^^^^^....................................***............",
	"..........*......................**.....^^^^^^^^^.....................................*.............",
	"........................................^^^^^^^^^^^^^^^^^^^^........................................",
	"........................................^^^^^^^^^^^^^^^^^^^^........................................",
	"........................................^^^^^^^^^^^^^^^^^^^^........................................",
	"........................................^^^^^^^^^^^^^^^^^^^^........................................",
	"........................................^^^^^^^^^^^^^^^^^^^^........................................",
	"...**...................................^^^^^^^^^^^^^^^^^^^^........................................",
	".*****..................................^^^^^^^^^^^^^^^^^^^^..................................***...",
	".*****..................................^^^^^^^^^^^^^^^^^^^^..................................***...",
	".*****..................................^^^^^^^^^^^^^^^^^^^^..................................***...",
	"...**...................................^^^^^^^^^^^^^^^^^^^^........................................",
	"........................................^^^^^^^^^^^^^^^^^^^^........................................",
	"........................................^^^^^^^^^^^^^^^^^^^^........................................",
	"........................................^^^^^^^^^^^^^^^^^^^^........................................",
	"........................................^^^^^^^^^^^^^^^^^^^^........................................",
	"....................................................^^^^^^^^.....**.................................",
	"........*...........................................^^^^^^^^.....***................................",
	".......***..............................^^^^^^^^^...^^^^^^^^.....***.......................***......",
	"........................................^^^^^^^^^...^^^^^^^^................................*.......",
	"........................................^^^^^^^^^...^^^^^^^^........................................",
	"........................................^^^^^^^^^...................................................",
	"........................................^^^^^^^^^...................................................",
	"........................................^^^^^^^^^^^^^^^^^^^^........................................",
}

func initMapForMission(g *gameMap, missionNumber int) {
	g.initTileMap(MIS1_MAP)

	ai_write("Seed is " + strconv.Itoa(rnd.Randomize()))

	g.factions = append(g.factions, createFaction("AI 1", 0,true, false))
	g.addBuilding(createBuilding("tcommand", 8, mapH/2, g.factions[0]), true)
	g.addPawn(createUnit("tscv", 8, mapH/2+4, g.factions[0], true))
	g.addPawn(createUnit("zdrone", 9, mapH/2+4, g.factions[0], true))
	g.addPawn(createUnit("pprobe", 10, mapH/2+4, g.factions[0], true))
	com := createUnit("ccommander", 11, mapH/2+4, g.factions[0], true)
	com.hitpoints = 1
	g.addPawn(com)

	g.factions[0].cursor.centralizeCamera()

	g.factions = append(g.factions, createFaction("AI 2", 1, false, true))
	g.addPawn(createUnit("tscv", mapW - 10, mapH/2, g.factions[1], true))
	CHEAT_IGNORE_FOW = false
}

func checkWinOrLose() { // TEMPORARY
	//if getCurrentTurn() % 10 != 0 {
	//	return
	//}
	//plrAlive := false
	//enemyAlive := false
	//for _, p := range CURRENT_MAP.pawns {
	//	if p.isCommander {
	//		if p.faction.playerControlled {
	//			plrAlive = true
	//		}
	//		if p.faction.aiControlled {
	//			enemyAlive = true
	//		}
	//	}
	//}
	//if !plrAlive {
	//	GAME_IS_RUNNING = false
	//	r_gamelostScreen()
	//	return
	//}
	//if !enemyAlive {
	//	GAME_IS_RUNNING = false
	//	r_gameWonScreen()
	//	return
	//}
}
