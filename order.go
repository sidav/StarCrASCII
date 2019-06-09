package main

type ORDER_TYPE_ENUM uint8

const (
	order_hold ORDER_TYPE_ENUM = iota
	order_move
	order_attack
	order_attack_move
	order_build  // maybe merge build and repair?
	order_construct
	order_gather_minerals
	order_return_resources
	order_enter_container
)

type order struct {
	orderType                    ORDER_TYPE_ENUM
	x, xSecondary, y, ySecondary int
	targetPawn                   *pawn

	buildingHasBeenPlaced bool // for build orders
	buildingToConstruct   *pawn
	constructingQueue     []*pawn // for units
}

func (clone *order) cloneFrom(o *order) {
	clone.orderType = o.orderType
	clone.x, clone.y = o.x, o.y
	clone.xSecondary, clone.ySecondary = o.xSecondary, o.ySecondary
	clone.targetPawn = o.targetPawn
	clone.buildingHasBeenPlaced = o.buildingHasBeenPlaced
	clone.buildingToConstruct = o.buildingToConstruct
}

func (o *order) canBeDrawnAsLine() bool {
	return o.orderType != order_hold && o.orderType != order_construct
}
