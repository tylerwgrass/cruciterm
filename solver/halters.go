package solver

type IHalter interface {
	Halt(*Navigator, *NavigationState) bool
	CheckInitialSquare() bool
}

type Halter struct {
	checkInitialSquare bool
}

type halterType int

const (
	ValidSquare halterType = iota
	EmptySquare
	ClueChange
)

func makeHalter(hType halterType, checkInitialSquare bool) IHalter {
	switch hType {
	case ValidSquare:
		return ValidSquareHalter{checkInitialSquare}
	case EmptySquare:
		return EmptySquareHalter{checkInitialSquare}
	case ClueChange:
		return ClueChangeHalter{checkInitialSquare}
	}
	return nil
}

type ValidSquareHalter Halter

func (h ValidSquareHalter) Halt(n *Navigator, state *NavigationState) bool {
	return (*n.grid)[state.row][state.col].content != "."
}

func (h ValidSquareHalter) CheckInitialSquare() bool {
	return h.checkInitialSquare
}

type EmptySquareHalter Halter

func (h EmptySquareHalter) Halt(n *Navigator, state *NavigationState) bool {
	return (*n.grid)[state.row][state.col].content == "-"
}

func (h EmptySquareHalter) CheckInitialSquare() bool {
	return h.checkInitialSquare
}

type ClueChangeHalter Halter

func (h ClueChangeHalter) Halt(n *Navigator, state *NavigationState) bool {
	return state.didChangeClue
}

func (h ClueChangeHalter) CheckInitialSquare() bool {
	return h.checkInitialSquare
}
