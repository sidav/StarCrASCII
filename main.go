package main

import (
	geometry "github.com/sidav/golibrl/geometry"
	cmenu "github.com/sidav/golibrl/console_menu"
	cw "github.com/sidav/golibrl/console"
	"strconv"
	"time"
)

func areCoordsValid(x, y int) bool {
	return geometry.AreCoordsInRect(x, y, 0, 0, mapW, mapH)
}

var (
	GAME_IS_RUNNING                     = true
	log                                 *LOG
	CURRENT_TICK                        = 1
	CURRENT_MAP                         *gameMap
	CURRENT_FACTION_SEEING_THE_SCREEN   *faction // for various rendering crap
	FIRE_WAS_OPENED_ON_SCREEN_THIS_TURN bool     // for killing pewpews overrender.
	CHEAT_IGNORE_FOW                    bool
	DEBUG_OUTPUT                        bool
)

func getCurrentTurn() int {
	return CURRENT_TICK/10 + 1
}

func debug_write(text string) {
	if DEBUG_OUTPUT {
		log.appendMessage("DEBUG: " + text)
	}
}

func main() {
	cw.Init_console("StarCrASCII", cw.TCellRenderer)
	defer cw.Close_console()

	log = &LOG{}

	CURRENT_MAP = &gameMap{}
	CURRENT_MAP.init()
	r_updateBoundsIfNeccessary(true)

	///////////////////////////////
	// uncomment later
	//r_showTitleScreen()
	//showBriefing()
	// comment later
	// endTurnPeriod = 0
	///////////////////////////////////

	for {
		startTime := time.Now()
		for _, f := range CURRENT_MAP.factions {
			f.recalculateSeenTiles()
			checkWinOrLose()
			if !GAME_IS_RUNNING {
				return
			}
			if f.aiControlled {
				ai_controlFaction(f)
			}
			if f.playerControlled {
				CURRENT_FACTION_SEEING_THE_SCREEN = f
				renderFactionStats(f)
				plr_control(f, CURRENT_MAP)
				debug_write("TOTAL FLUSHES: " + strconv.Itoa(cw.GetNumberOfRecentFlushes()))
			}
		}
		for i := 0; i < 10; i++ {
			for _, u := range CURRENT_MAP.pawns {
				if u.hitpoints <= 0 {
					log.appendMessage(u.name + " is destroyed!")
					CURRENT_MAP.removePawn(u)
					continue
				}
				if u.regenPeriod > 0 && CURRENT_TICK%u.regenPeriod == 0 && u.hitpoints < u.maxHitpoints {
					u.hitpoints++
				}
				u.executeOrders(CURRENT_MAP)
				u.openFireIfPossible()
			}
			if FIRE_WAS_OPENED_ON_SCREEN_THIS_TURN {
				cw.Flush_console()
				FIRE_WAS_OPENED_ON_SCREEN_THIS_TURN = false
				time.Sleep(time.Duration(endTurnPeriod/4) * time.Millisecond)
			}
			CURRENT_TICK += 1
		}

		for _, f := range CURRENT_MAP.factions {
			f.recalculateSupply(CURRENT_MAP)
		}
		doAllProduction(CURRENT_MAP)
		CURRENT_MAP.cleanupMinerals()
		timeForTurn := int(time.Since(startTime) / time.Millisecond)
		debug_write("Time for turn: " + strconv.Itoa(timeForTurn) + "ms") // TODO: make it removable
	}

}

func showBriefing() {
	cw.Clear_console()
	text := "Placeholder briefing."
	cmenu.DrawWrappedTextInRect(text, 0, 0, CONSOLE_W, CONSOLE_H)
	cw.Flush_console()
	key := ""
	for key != "ESCAPE" && key != "ENTER" {
		key = cw.ReadKey()
	}
}
