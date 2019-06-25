package main
import cw "github.com/sidav/golibrl/console"


func attackWithWeapon(wpn *pawnWeaponInformation, attacker, preferredTarget *pawn) bool { // true if was fired
	attackerCenterX, attackerCenterY := attacker.getCenter()
	if (wpn.canBeFiredOnMove && wpn.nextTurnToFire > CURRENT_TICK) || (!wpn.canBeFiredOnMove && !attacker.isTimeToAct()) {
		// log.appendMessage(fmt.Sprintf("Skipping fire: TtA:%b CBFoM:%b TRN: %b", attacker.isTimeToAct() ,wpn.canBeFiredOnMove, wpn.nextTurnToFire > CURRENT_TICK))
		return false
	}
	var target *pawn
	radius := wpn.attackRadius
	if preferredTarget != nil && attacker.isInDistanceFromPawn(preferredTarget, radius) {
		target = preferredTarget
	} else {
		potential_targets := CURRENT_MAP.getEnemyPawnsInRadiusFromPawn(attacker, radius, attacker.faction)
		for _, potentialTarget := range potential_targets {
			ptx, pty := potentialTarget.getCoords()
			if attacker.faction.areCoordsInSight(ptx, pty) || attacker.faction.areCoordsInRadarRadius(ptx, pty) {
				target = potentialTarget
			}
		}
	}
	if target != nil {
		// draw the pew pew laser TODO: move this crap somewhere already
		if areGlobalCoordsOnScreenForFaction(attackerCenterX, attackerCenterY, CURRENT_FACTION_SEEING_THE_SCREEN) || areGlobalCoordsOnScreenForFaction(target.x, target.y, CURRENT_FACTION_SEEING_THE_SCREEN) {
			cw.SetFgColor(cw.RED)
			cx, cy := target.getCenter()
			camx, camy := CURRENT_FACTION_SEEING_THE_SCREEN.cursor.getCameraCoords()
			renderLine(attackerCenterX, attackerCenterY, cx, cy, false, camx, camy)
			FIRE_WAS_OPENED_ON_SCREEN_THIS_TURN = true
		}
		dealDamageToTarget(attacker, wpn, target)
		return true
	}
	return false
}
