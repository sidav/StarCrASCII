package main

type pawnContainerInformation struct {
	maxSize             int
	pawnsInside         []*pawn
	allowFireFromInside bool
}

func (p *pawn) canContainPawns() bool {
	return p.containerInfo != nil
}
