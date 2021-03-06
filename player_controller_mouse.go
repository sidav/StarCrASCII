package main

import (
	"fmt"
	cw "github.com/sidav/golibrl/console"
	cmenu "github.com/sidav/golibrl/console_menu"
	"time"
)

func plr_selectPawnWithMouse(f *faction, m *gameMap) *[]*pawn { // returns a pointer to an array of selected pawns.
	f.cursor.currentCursorMode = CURSOR_SELECT
	for {
		if reRenderNeeded {
			r_renderScreenForFaction(f, m, nil, true)
		}
		keyPressed := cw.ReadKeyAsync()
		click := cw.GetMouseClickedButton()
		reRenderNeeded = true

		if plr_moveCameraOrCursorWithMouseIfNeeded(f) {
			return nil
		}
		u := f.cursor.snappedPawn
		if cw.GetMouseHeldButton() == "LEFT" {
			return plr_bandboxSelectionWithMouse(f)
		}
		if u != nil && click == "LEFT" {
			if u.faction.factionNumber != f.factionNumber {
				log.appendMessage("Enemy units can't be selected, Commander.")
				return nil
			}
			return &[]*pawn{f.cursor.snappedPawn}
		}

		switch keyPressed {
		case "NOTHING", "NON-KEY":
			if !IS_PAUSED && isTimeToAutoEndTurn() {
				last_time = time.Now()
				PLR_LOOP = false // end turn
				return nil
			} else {
				reRenderNeeded = false
			}
		case ".": // end turn without unpausing the game
			if IS_PAUSED {
				PLR_LOOP = false
				return nil
			}
		case "`":
			mouseEnabled = !mouseEnabled
			if mouseEnabled {
				log.appendMessage("Mouse controls enabled.")
			} else {
				log.appendMessage("Mouse controls disabled.")
			}
		case "SPACE", " ":
			IS_PAUSED = !IS_PAUSED
			if IS_PAUSED {
				log.appendMessage("Tactical pause engaged.")
			} else {
				log.appendMessage("Switched to real-time mode.")
			}
		case "=":
			if endTurnPeriod > 100 {
				endTurnPeriod -= 100
				log.appendMessagef("Game speed increased to %d", 10-(endTurnPeriod/100))
			} else {
				log.appendMessage("Can't increase game speed any further.")
			}
		case "-":
			if endTurnPeriod < 2000 {
				endTurnPeriod += 100
				log.appendMessagef("Game speed decreased to %d", 10-(endTurnPeriod/100))
			} else {
				log.appendMessage("Can't decrease game speed any further.")
			}

		case "ENTER", "RETURN":
			u := f.cursor.snappedPawn //m.getUnitAtCoordinates(cx, cy)
			if u == nil {
				return plr_bandboxSelectionWithMouse(f) // select multiple units
			}
			if u.faction.factionNumber != f.factionNumber {
				log.appendMessage("Enemy units can't be selected, Commander.")
				return nil
			}
			return &[]*pawn{f.cursor.snappedPawn}
		case "TAB":
			trySelectNextIdlePawn(f)
		case "C":
			trySnapCursorToCommander(f)
			return &[]*pawn{f.cursor.snappedPawn}
		case "?":
			if f.cursor.snappedPawn != nil {
				renderPawnInfo(f.cursor.snappedPawn)
			}
		case "ESCAPE":
			if cmenu.ShowSimpleYNChoiceModalWindow("Are you sure you want to quit?") {
				GAME_IS_RUNNING = false
				PLR_LOOP = false
				return nil
			}

		case "DELETE": // cheat
			f.economy.minerals += 10000
			f.economy.vespene += 10000
		case "E": // test
			CURRENT_MAP.addBuilding(createBuilding("testsmall", f.cursor.x, f.cursor.y, CURRENT_MAP.factions[1]), true)
			log.appendMessage("Test enemy building created.")
		case "B": // test
			CURRENT_MAP.addBuilding(createBuilding("testbig", f.cursor.x, f.cursor.y, CURRENT_MAP.factions[1]), true)
			log.appendMessage("LARGE test enemy building created.")
		case "INSERT": // cheat
			for _, fac := range CURRENT_MAP.factions {
				if fac != f {
					fac.economy.minerals += 500
					fac.economy.vespene += 500
				}
			}
			log.appendMessage("Added 500 minerals and gas to the enemies.")
		case "HOME": // cheat
			// CURRENT_MAP.addBuilding(createBuilding("lturret", f.cursor.x, f.cursor.y, CURRENT_MAP.factions[0]), true)
			endTurnPeriod = 0
		case "END": // cheat
			CHEAT_IGNORE_FOW = !CHEAT_IGNORE_FOW

		default:
			plr_moveCursor(f, keyPressed)
		}
	}
}

func plr_bandboxSelectionWithMouse(f *faction) *[]*pawn {
	f.cursor.currentCursorMode = CURSOR_MULTISELECT
	f.cursor.xorig, f.cursor.yorig = f.cursor.getCoords()
	reRenderNeeded = true
	for {
		if reRenderNeeded {
			r_renderScreenForFaction(f, CURRENT_MAP, nil, true)
		}
		keyPressed := cw.ReadKeyAsync()
		if keyPressed == "ESCAPE" {
			return nil
		}
		if cw.GetMouseHeldButton() == "NONE" {
			reRenderNeeded = true
			fromx, fromy := f.cursor.xorig, f.cursor.yorig
			tox, toy := f.cursor.getCoords()
			if fromx > tox {
				t := fromx
				fromx = tox
				tox = t
			}
			if fromy > toy {
				t := fromy
				fromy = toy
				toy = t
			}
			unitsInSelection := CURRENT_MAP.getPawnsInRect(fromx, fromy, tox-fromx+1, toy-fromy+1)
			unitsToReturn := make([]*pawn, 0)
			for _, p := range unitsInSelection {
				// select only the pawns if current faction which are capable to move AND attack and are not commanders.
				if p.faction != nil && p.faction == f && p.hasWeapons() && p.canMove() && !p.isCommander {
					unitsToReturn = append(unitsToReturn, p)
				}
			}
			log.appendMessage(fmt.Sprintf("%d units selected from %d", len(unitsToReturn), len(unitsInSelection)))
			return &unitsToReturn
		}
		if cw.WasMouseMovedSinceLastEvent() {
			plr_moveCursorWithMouse(f)
		} else {
			reRenderNeeded = false
		}
	}
}

func plr_giveOrderWithMouse(selection *[]*pawn, f *faction) {
	selectedPawn := (*selection)[0] //m.getUnitAtCoordinates(cx, cy)
	log.appendMessage(selectedPawn.name + " is awaiting orders.")
	f.cursor.currentCursorMode = CURSOR_MOVE
	reRenderNeeded = true
	for {
		equivKey := "NONE" // mouse clicked menu result
		click := cw.GetMouseClickedButton()
		cx, cy := f.cursor.getCoords()
		if reRenderNeeded {
			r_renderScreenForFaction(f, CURRENT_MAP, selection, false)
		}
		equivKey = pcm_mouseOrderSelectMenu(selectedPawn)
		if reRenderNeeded {
			cw.Flush_console()
		}

		keyPressed := cw.ReadKeyAsync()
		if plr_moveCameraOrCursorWithMouseIfNeeded(f) {
			continue
		}
		if click == "LEFT" {
			if equivKey == "NONE" && areGlobalCoordsOnScreen(cx, cy) {
				reRenderNeeded = true
				issueDefaultOrderToUnit(selectedPawn, CURRENT_MAP, cx, cy)
				return
			} else {
				keyPressed = equivKey
			}
		}
		if click == "RIGHT" {
			reRenderNeeded = true
			return
		}

		switch keyPressed {
		case "a": // attack-move
			if selectedPawn.hasWeapons() || selectedPawn.canConstructUnits() {
				f.cursor.currentCursorMode = CURSOR_AMOVE
				reRenderNeeded = true
			}
		case "m": // move
			f.cursor.currentCursorMode = CURSOR_MOVE
			reRenderNeeded = true
		case "b": // build
			if selectedPawn.canConstructBuildings() {
				code := plr_selectBuidingToConstruct(selectedPawn)
				if code != "" {
					plr_selectBuildingSiteWithMouse(selectedPawn, createBuilding(code, cx, cy, f), CURRENT_MAP)
					return
				}
			}
		case "c": // construct units
			if selectedPawn.canConstructUnits() {
				plr_selectUnitsToConstruct(selectedPawn)
				reRenderNeeded = true
			}
		case "r": // repeat construction queue
			if selectedPawn.canConstructUnits() {
				selectedPawn.repeatConstructionQueue = !selectedPawn.repeatConstructionQueue
				reRenderNeeded = true
			}
		case "u": // unload units inside
			if selectedPawn.canReleaseContainedPawns() {
				selectedPawn.setOrder(&order{orderType: order_unload})
				return
			}
		case "ESCAPE":
			return
		default:
			reRenderNeeded = false
		}
	}
}

func plr_giveOrderForMultiSelectWithMouse(selection *[]*pawn, f *faction) {
	log.appendMessage(fmt.Sprintf("%d units are awaiting orders.", len(*selection)))
	f.cursor.currentCursorMode = CURSOR_MOVE
	reRenderNeeded = true
	for {
		equivKey := "NONE"
		click := cw.GetMouseClickedButton()
		cx, cy := f.cursor.getCoords()

		if reRenderNeeded {
			r_renderScreenForFaction(f, CURRENT_MAP, selection, false)
		}
		equivKey = pcm_MouseOrderForMultiselectMenu(f, selection)
		if reRenderNeeded {
			cw.Flush_console()
		}

		keyPressed := cw.ReadKeyAsync()
		if plr_moveCameraOrCursorWithMouseIfNeeded(f) {
			continue
		}
		if click == "LEFT" {
			if equivKey == "NONE" && areGlobalCoordsOnScreen(cx, cy) {
				for _, p := range *selection {
					issueDefaultOrderToUnit(p, CURRENT_MAP, cx, cy)
				}
				reRenderNeeded = true
				return
			} else {
				keyPressed = equivKey
			}
		}
		if click == "RIGHT" {
			reRenderNeeded = true
			return
		}

		switch keyPressed {
		case "ENTER", "RETURN":
			for _, p := range *selection {
				issueDefaultOrderToUnit(p, CURRENT_MAP, cx, cy)
			}
			return
		case "a": // attack-move
			f.cursor.currentCursorMode = CURSOR_AMOVE
		case "m": // move
			f.cursor.currentCursorMode = CURSOR_MOVE
		case "ESCAPE":
			return
		default:
			reRenderNeeded = false
		}
	}
}

func plr_selectBuildingSiteWithMouse(p *pawn, b *pawn, m *gameMap) {
	log.appendMessage("Select construction site for " + b.name)
	reRenderNeeded = true
	for {
		f := p.faction
		cursor := f.cursor
		cx, cy := cursor.getCoords()
		click := cw.GetMouseClickedButton()
		cursor.currentCursorMode = CURSOR_BUILD

		cursor.buildingToConstruct = b
		if b.buildingInfo.allowsTightPlacement {
			cursor.w = b.buildingInfo.w
			cursor.h = b.buildingInfo.h
		} else {
			cursor.w = b.buildingInfo.w + 2
			cursor.h = b.buildingInfo.h + 2
		}

		cursor.radius = b.getMaxRadiusToFire()

		if reRenderNeeded { // TODO: move all that "if reRenderNeeded" to the renderer itself to keep code more clean.
			r_renderScreenForFaction(f, m, nil, true)
		}

		keyPressed := cw.ReadKeyAsync()

		if plr_moveCameraOrCursorWithMouseIfNeeded(f) {
			continue
		}

		if click == "LEFT" {
			if m.canBuildingBeBuiltAt(b, cx, cy) {
				b.x = cx - b.buildingInfo.w/2
				b.y = cy - b.buildingInfo.h/2
				p.setOrder(&order{orderType: order_build, x: cx, y: cy, buildingToConstruct: b})
				reRenderNeeded = true
				return
			} else {
				log.appendMessage("This building can't be placed here!")
			}
		}
		if click == "RIGHT" {
			reRenderNeeded = true
			log.appendMessage("Construction cancelled: " + b.name)
			return
		}

		switch keyPressed {
		case "ESCAPE":
			reRenderNeeded = true
			log.appendMessage("Construction cancelled: " + b.name)
			return
		default:
			reRenderNeeded = false
		}
	}
}

func plr_moveCameraOrCursorWithMouseIfNeeded(f *faction) bool { // returns true if something was actually moved.
	if cw.WasMouseMovedSinceLastEvent() {
		if cw.GetMouseHeldButton() == "RIGHT" {
			plr_moveCameraWithMouse(f)
			return true
		} else {
			plr_moveCursorWithMouse(f)
			return true
		}
	}
	return false
}

func plr_moveCameraWithMouse(f *faction) {
	vx, vy := cw.GetMouseMovementVector()
	if vx == 0 && vy == 0 {
		reRenderNeeded = false
		return
	}
	cx, cy := f.cursor.getCoords()
	if areCoordsValid(cx+vx, cy+vy) {
		f.cursor.cameraX += vx
		f.cursor.cameraY += vy
	}

	snapB := f.cursor.snappedPawn
	if snapB != nil { // unsnap cursor
		for snapB.isOccupyingCoords(f.cursor.x, f.cursor.y) {
			if areCoordsValid(f.cursor.x+vx, f.cursor.y+vy) {
				f.cursor.moveByVector(vx, vy)
			} else {
				break
			}
		}
		f.cursor.snappedPawn = nil
	}
	if f.cursor.currentCursorMode != CURSOR_BUILD {
		snapCursorToPawn(f)
	}
	reRenderNeeded = true
}

func plr_moveCursorWithMouse(f *faction) {
	cx, cy := cw.GetMouseCoords()
	camx, camy := f.cursor.getCameraCoords()

	reRenderNeeded = !(f.cursor.x == camx+cx && f.cursor.y == camy+cy) // rerender is needed if cursor was _actually_ moved

	if areCoordsValid(camx+cx, camy+cy) {
		f.cursor.x, f.cursor.y = camx+cx, camy+cy
		snapCursorToPawn(f)
	}
}
