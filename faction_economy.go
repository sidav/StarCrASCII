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
