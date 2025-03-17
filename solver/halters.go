package solver

import (
	"fmt"

	"github.com/tylerwgrass/cruciterm/logger"
)

type IHalter interface {
	Halt(*NavigationGrid, int, int) bool
	CheckInitialSquare() bool
}

type Halter struct {
	checkInitialSquare bool
}

type halterType int
const (
	ValidSquare halterType = iota
	EmptySquare
)

func makeHalter(hType halterType, checkInitialSquare bool) IHalter {
	switch hType {
	case ValidSquare:
		return ValidSquareHalter{checkInitialSquare} 
	case EmptySquare:
		return EmptySquareHalter{checkInitialSquare}
	}
	return nil
}

type ValidSquareHalter Halter
func (h ValidSquareHalter) Halt(g *NavigationGrid, i, j int) bool {
	logger.Debug(fmt.Sprintf("Checking valid square at [%d,%d]: %s", i, j, (*g)[i][j].content))
	return (*g)[i][j].content != "." 
}

func (h ValidSquareHalter) CheckInitialSquare() bool {
	return h.checkInitialSquare
}

type EmptySquareHalter Halter
func (h EmptySquareHalter) Halt(g *NavigationGrid, i, j int) bool {
	return (*g)[i][j].content == "-"
}

func (h EmptySquareHalter) CheckInitialSquare() bool {
	return h.checkInitialSquare
}