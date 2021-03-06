package main

type pawnContainerInformation struct {
	maxSize             int
	pawnsInside         []*pawn
	allowFireFromInside bool
}

func (p *pawn) canContainPawns() bool {
	return p.containerInfo != nil
}

func (c *pawnContainerInformation) addPawnToContainer(p *pawn) {
	c.pawnsInside = append(c.pawnsInside, p)
}

func (c *pawnContainerInformation) canAddPawn(p *pawn) bool { // TODO: count units' individual size
	return c.maxSize - len(c.pawnsInside) > 0
}
