package main

import (
	"fmt"
	cmenu "github.com/sidav/golibrl/console_menu"
	"github.com/sidav/golibrl/geometry"
	primitives "github.com/sidav/golibrl/graphic_primitives"
	vectorMath "github.com/sidav/golibrl/math"
	"sort"
)
import cw "github.com/sidav/golibrl/console"

func renderFactionStats(f *faction) {
	eco := f.economy
	statsx := VIEWPORT_W + 1

	// fr, fg, fb := getFactionRGB(f.factionNumber)
	// cw.SetFgColorRGB(fr, fg, fb)
	if IS_PAUSED {
		cw.SetFgColor(f.getFactionColor())
		cw.PutString(f.name+": ", statsx, 0)
		cw.SetFgColor(cw.YELLOW)
		cw.PutString(fmt.Sprintf("turn %d (PAUSED)", getCurrentTurn()), statsx+len(f.name)+2, 0)
	} else {
		cw.SetFgColor(f.getFactionColor())
		cw.PutString(fmt.Sprintf("%s: turn %d", f.name, getCurrentTurn()), statsx, 0)
	}

	cw.SetFgColor(cw.DARK_CYAN)
	cw.PutString(fmt.Sprintf("MINERALS %d", eco.minerals), statsx, 2)
	cw.SetFgColor(cw.DARK_GREEN)
	cw.PutString(fmt.Sprintf("VESPENE %d", eco.vespene), statsx, 3)
	cw.SetFgColor(cw.DARK_BLUE)
	renderStatusbar(fmt.Sprintf("SUPPLY %d/%d", eco.cursupply, eco.maxsupply), eco.cursupply, eco.maxsupply,
		statsx, 4, SIDEBAR_W-3, cw.DARK_BLUE)
	// cw.PutString(fmt.Sprintf("SUPPLY %d/%d", eco.cursupply, eco.maxsupply), statsx, 4)

}

func renderInfoOnCursor(f *faction, g *gameMap) {

	title := "Unidentified Object"
	color := 2
	details := make([]string, 0)
	// var res *pawnResourceInformation
	sp := f.cursor.snappedPawn

	if sp != nil {

		if sp.faction == f {
			renderOrderLine(sp)
		}
		color = sp.faction.getFactionColor()
		if CURRENT_FACTION_SEEING_THE_SCREEN.areCoordsInSight(sp.x, sp.y) {
			title = sp.name
			if sp.faction != f {
				if sp.isBuilding() {
					details = append(details, "(Enemy building)")
				} else {
					details = append(details, "(Enemy unit)")
				}
			} else {
				details = append(details, sp.getCurrentOrderDescription())
				if sp.res != nil && sp.currentConstructionStatus == nil {
					// res = sp.res
				}
			}
			r_renderAttackRadius(sp)
		} else {
			details = append(details, "(Enemy on radar)")
		}
	}

	if len(details) > 0 {
		if CURRENT_FACTION_SEEING_THE_SCREEN.areCoordsInSight(sp.x, sp.y) {
			details = append(details, sp.getArmorDescriptionString())
			if sp.hasWeapons() {
				for _, wpn := range sp.weapons {
					details = append(details, wpn.getDescriptionString())
				}
			}
			//if res != nil {
			//	economyInfo := fmt.Sprintf("METAL: (+%d / -%d) ENERGY: (+%d / -%d)",
			//		res.metalIncome, res.metalSpending, res.energyIncome, res.energySpending+res.energyDrain)
			//	details = append(details, economyInfo)
			//}
		}
		cmenu.DrawSidebarInfoMenu(title, color, SIDEBAR_X, SIDEBAR_FLOOR_2, SIDEBAR_W, details)
	}
}

func r_renderPossibleOrdersForPawn(p *pawn) {
	orders := make([]string, 0)
	if p.canConstructBuildings() {
		orders = append(orders, "(B)uild")
	}
	if p.canConstructUnits() {
		orders = append(orders, "(C)onstruct units")
		if p.repeatConstructionQueue {
			orders = append(orders, "(R)epeat queue: ENABLED")
		} else {
			orders = append(orders, "(R)epeat queue: DISABLED")
		}
		if p.faction.cursor.currentCursorMode == CURSOR_AMOVE {
			orders = append(orders, "Default order: attack-move")
			orders = append(orders, "(m): set to move")
		} else {
			orders = append(orders, "Default order: move")
			orders = append(orders, "(a): set to attack-move")
		}
	}
	if p.faction.cursor.currentCursorMode == CURSOR_AMOVE && p.canMove() {
		orders = append(orders, "(M)ove")
	} else {
		if p.hasWeapons() {
			orders = append(orders, "(A)ttack-move")
		}
	}
	cmenu.DrawSidebarInfoMenu("Orders for: "+p.name, p.faction.getFactionColor(),
		SIDEBAR_X, SIDEBAR_FLOOR_3, SIDEBAR_W, orders)
}

func r_renderPossibleOrdersForMultiselection(f *faction, selection *[]*pawn) {
	orders := make([]string, 0)
	selectedUnitsCounter := make(map[string]int)
	for _, p := range *selection {
		selectedUnitsCounter[p.name]++
	}
	// sort the map because of dumbass Go developers thinking that they know your needs better than you do
	keys := make([]string, 0, len(selectedUnitsCounter))
	for k := range selectedUnitsCounter {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		orders = append(orders, fmt.Sprintf("%dx %s", selectedUnitsCounter[key], key))
	}
	if f.cursor.currentCursorMode == CURSOR_AMOVE {
		orders = append(orders, "(M)ove")
	} else {
		orders = append(orders, "(A)ttack-move")
	}
	cmenu.DrawSidebarInfoMenu(fmt.Sprintf("ORDERS FOR %d UNITS", len(*selection)), f.getFactionColor(),
		SIDEBAR_X, SIDEBAR_FLOOR_3, SIDEBAR_W, orders)
}

func renderStatusbar(name string, curvalue, maxvalue, x, y, width, barColor int) {
	barTitle := name
	cw.PutString(barTitle, x, y)
	barWidth := width - len(name)
	var filledCells int
	if maxvalue > 0 {
		filledCells = barWidth * curvalue / maxvalue
	} else {
		filledCells = 0
	}
	barStartX := x + len(barTitle) + 1
	for i := 0; i < barWidth; i++ {
		if i < filledCells {
			cw.SetFgColor(barColor)
			cw.PutChar('=', i+barStartX, y)
		} else {
			cw.SetFgColor(cw.DARK_BLUE)
			cw.PutChar('-', i+barStartX, y)
		}
	}
}

func renderLog(flush bool) {
	cw.SetFgColor(cw.WHITE)
	for i := 0; i < LOG_HEIGHT; i++ {
		cw.PutString(log.last_msgs[i].getText(), 0, VIEWPORT_H+i+1)
	}
	if flush {
		flushView()
	}
}

func r_renderAttackRadius(p *pawn) {
	if len(p.weapons) > 0 {
		// (p.x, p.y, p.weapons[0].attackRadius, false, CURRENT_FACTION_SEEING_THE_SCREEN.cursor.x-VIEWPORT_W/2, CURRENT_FACTION_SEEING_THE_SCREEN.cursor.y-VIEWPORT_H/2)
		vx, vy := CURRENT_FACTION_SEEING_THE_SCREEN.cursor.getCameraCoords()
		var line *[]primitives.Point
		if p.isBuilding() {
			line = primitives.GetApproxCircleAroundRect(p.x, p.y, p.buildingInfo.w, p.buildingInfo.h, p.weapons[0].attackRadius)
		} else {
			px, py := p.getCenter()
			line = primitives.GetCircle(px, py, p.weapons[0].attackRadius)
		}
		for _, point := range *line {
			x, y := point.X, point.Y
			cw.SetFgColor(cw.BLACK)
			if geometry.AreCoordsInRect(x-vx, y-vy, 0, 0, VIEWPORT_W, VIEWPORT_H) && areCoordsValid(x, y) {
				tileApp := CURRENT_MAP.tileMap[x][y].appearance
				cw.SetBgColor(tileApp.color)
				cw.PutChar(tileApp.char, x-vx, y-vy)
			}
		}
		cw.SetBgColor(cw.BLACK)
	}
}

func renderOrderLine(p *pawn) {
	var ordr *order
	if p.order != nil {
		if p.order.canBeDrawnAsLine() {
			ordr = p.order
		}
	}
	if ordr == nil {
		if p.canConstructUnits() {
			ordr = p.productionInfo.defaultOrderForUnitBuilt
		}
	}
	if ordr != nil {
		if ordr.orderType == order_attack {
			cw.SetFgColor(cw.RED)
		} else {
			cw.SetFgColor(cw.YELLOW)
		}
		f := p.faction
		cx, cy := p.getCenter()
		camx, camy := f.cursor.getCameraCoords()
		renderLine(cx, cy, ordr.x, ordr.y, false, camx, camy)
	}
}

func renderLine(fromx, fromy, tox, toy int, flush bool, vx, vy int) {
	line := primitives.GetLine(fromx, fromy, tox, toy)
	char := '?'
	if len(line) > 1 {
		dirVector := vectorMath.CreateVectorByStartAndEndInt(fromx, fromy, tox, toy)
		dirVector.TransformIntoUnitVector()
		dirx, diry := dirVector.GetRoundedCoords()
		char = getTargetingChar(dirx, diry)
	}
	//if fromx == tox && fromy == toy {
	//	renderPawn(d.player, true)
	//}
	for i := 1; i < len(line); i++ {
		// x, y := line[i].X, line[i].Y
		//if d.isPawnPresent(x, y) {
		//	renderPawn(d.getPawnAt(x, y), true)
		//} else {
		// cw.SetFgColor(cw.YELLOW)
		if i == len(line)-1 {
			char = 'X'
		}
		viewx, viewy := line[i].X-vx, line[i].Y-vy
		if geometry.AreCoordsInRect(viewx, viewy, 0, 0, VIEWPORT_W, VIEWPORT_H) {
			cw.PutChar(char, viewx, viewy)
		}
		// }
	}
	if flush {
		cw.Flush_console()
	}
}

func renderPawnInfo(p *pawn) {
	var name, desc string
	if p.isUnit() {
		name, desc = getUnitNameAndDescription(p.codename)
	} else if p.isBuilding() {
		name, desc = getBuildingNameAndDescription(p.codename)
	}
	cmenu.ShowSimpleInfoWindow(name, desc, 60, 15, p.faction.getFactionColor())
}

func renderCircle(fromx, fromy, radius int, char rune, flush bool) {
	if radius == 0 {
		return
	} else {
		line := primitives.GetCircle(fromx, fromy, radius)
		for _, point := range *line {
			x, y := point.X, point.Y
			renderCharByGlobalCoords(char, x, y)
		}
	}
	if flush {
		cw.Flush_console()
	}
}

func renderApproxCircleAroundRect(fromx, fromy, width, height, radius int, char rune, flush bool) {
	if radius == 0 {
		return
	} else {
		line := primitives.GetApproxCircleAroundRect(fromx, fromy, width, height, radius)
		for _, point := range *line {
			x, y := point.X, point.Y
			renderCharByGlobalCoords(char, x, y)
		}
	}
	if flush {
		cw.Flush_console()
	}
}

func getTargetingChar(x, y int) rune {
	if abs(x) > 1 {
		x /= abs(x)
	}
	if abs(y) > 1 {
		y /= abs(y)
	}
	if x == 0 {
		return '|'
	}
	if y == 0 {
		return '-'
	}
	if x*y == 1 {
		return '\\'
	}
	if x*y == -1 {
		return '/'
	}
	return '?'
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
