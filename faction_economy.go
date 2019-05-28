package main

type factionEconomy struct {
	minerals, vespene    int
	cursupply, maxsupply int
}

func (f *faction) canAffordSpend(m, v int) bool {
	return f.economy.minerals >= m && f.economy.vespene >= v
}

func (f *faction) spendResources(m, v int) {
	f.economy.minerals -= m
	f.economy.vespene -= v
}

func (f *faction) recalculateSupply(g *gameMap) { // move somewhere?
	f.economy.cursupply = 0
	f.economy.maxsupply = 0
	for _, p := range CURRENT_MAP.pawns {
		if p.faction == f && p.currentConstructionStatus == nil {
			f.economy.cursupply += p.takesSupply
			f.economy.maxsupply += p.givesSupply
		}
	}
}
