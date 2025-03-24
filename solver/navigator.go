package solver

import (
	"fmt"
	"strings"

	"github.com/tylerwgrass/cruciterm/logger"
	"github.com/tylerwgrass/cruciterm/puzzle"
)

type Cell struct {
	content string
	acrossClue *puzzle.Clue
	downClue *puzzle.Clue 
	nextAcross *puzzle.Clue  
	nextDown *puzzle.Clue  
	prevAcross *puzzle.Clue 
	prevDown *puzzle.Clue 
	isAcrossClueEnd bool
	isDownClueEnd bool
}

type IterationMode int
const (
	Clues IterationMode = iota
	Cardinal
)

type NavigationGrid [][]Cell

type Navigator struct {
	grid *NavigationGrid
	orientation Orientation
	direction Direction
	iterMode IterationMode
	halters []IHalter
}

type NavigationDeltas struct {
	dr int
	dc int
}

type NavigationState struct {
	startRow int
	startCol int
	row int
	col int
	didWrap bool
	didChangeClue bool
}

var defaultHalter = makeHalter(ValidSquare, false) 
func NewNavigator(puzzleGrid [][]string, puz *puzzle.PuzzleDefinition) *Navigator {
	navGrid := make(NavigationGrid, len(puzzleGrid))
	acrosses := puzzle.AcrossClues
	downs := puzzle.DownClues
	currentAcrossIndex := 0
	currentDownIndex := 0
	prevAcross := acrosses[len(acrosses) - 1]
	nextAcross := acrosses[(currentAcrossIndex + 1) % len(acrosses)]
	prevDown := downs[len(downs) - 1]
	nextDown := downs[(currentDownIndex + 1) % len(downs)]

	for row := range puz.NumRows {
		navGrid[row] = make([]Cell, len(puzzleGrid[0]))
		for col := range puz.NumCols {
			if row == nextAcross.StartRow && col == nextAcross.StartCol {
				prevAcross = acrosses[currentAcrossIndex]
				currentAcrossIndex = (currentAcrossIndex + 1) % len(acrosses)
				nextAcross = acrosses[(currentAcrossIndex + 1) % len(acrosses)]
			}
			if row == nextDown.StartRow && col == nextDown.StartCol {
				prevDown = downs[currentDownIndex]
				currentDownIndex = (currentDownIndex + 1) % len(downs)
				nextDown = downs[(currentDownIndex + 1) % len(downs)]
			}
			cell := Cell{
				content: puzzleGrid[row][col], 
			}
			if col == acrosses[currentAcrossIndex].EndCol {
				cell.isAcrossClueEnd = true
			}
			if row == downs[currentDownIndex].EndRow {
				cell.isDownClueEnd = true
			}
			if cell.content != "." {
				cell.prevAcross = prevAcross
				cell.nextAcross = nextAcross
				cell.prevDown = prevDown
				cell.nextDown = nextDown
				cell.acrossClue = acrosses[currentAcrossIndex]
				if row > 0 && navGrid[row - 1][col].content != "." {
					cell.downClue = navGrid[row - 1][col].downClue
					cell.nextDown = navGrid[row - 1][col].nextDown
					cell.prevDown = navGrid[row - 1][col].prevDown
				} else {
					cell.downClue = downs[currentDownIndex]
				}
			}
			navGrid[row][col] = cell
		}
	}
	return &Navigator{
		grid: &navGrid,
		orientation: Horizontal,
		direction: Forward,
		iterMode: Clues,
		halters: []IHalter{defaultHalter},
	}
}
func (n *Navigator) resetNavigatorOptions() {
	n.orientation = Horizontal
	n.direction = Forward
	n.iterMode = Clues
	n.halters = []IHalter{defaultHalter}
}

func (n *Navigator) withOrientation(o Orientation) *Navigator {
	n.orientation = o
	return n
}

func (n *Navigator) withDirection(d Direction) *Navigator {
	n.direction = d
	return n 
}

func (n *Navigator) withHalter(h IHalter) *Navigator {
	n.halters = []IHalter{h}
	return n
}

func (n *Navigator) withHalters(h []IHalter) *Navigator {
	n.halters = h
	return n
}

func (n *Navigator) withIterMode(i IterationMode) *Navigator {
	n.iterMode = i
	return n
}

func (navigator Navigator) advanceCursor(startCol, startRow int) []NavigationState {
	navStates := make([]NavigationState, 0, len(navigator.halters))
	row, col := startRow, startCol
	for _, halter := range(navigator.halters) {
		navState := NavigationState{row: row, col: col, startRow: row, startCol: col}
		if navigator.iterMode == Cardinal {
			navigator.iterateCardinal(&navState, halter)
		} else {
			navigator.iterateClues(&navState, halter)
		}
		navStates = append(navStates, navState)
	}
	return navStates
}

func (navigator Navigator) iterateCardinal(state *NavigationState, halter IHalter) {
	grid := *navigator.grid
	if halter.CheckInitialSquare() && halter.Halt(&navigator, state) {
		return 
	}
	deltas := navigator.getDeltas()
	for {
		if state.didWrap && state.row == state.startRow && state.col == state.startCol {
			return
		}
		nextRow, nextCol := state.row + deltas.dr, state.col + deltas.dc
		if grid.isVisitable(nextRow, nextCol) {
			state.row = nextRow
			state.col = nextCol
		} else {
			navigator.moveToNextValidCardinal(state)
		} 
		if halter.Halt(&navigator, state) {
			return 
		}
	}
}

func (navigator Navigator) iterateClues(state *NavigationState, halter IHalter) {
	grid := *navigator.grid
	var startClue *puzzle.Clue
	if navigator.orientation == Horizontal {
		startClue = (*navigator.grid)[state.startRow][state.startCol].acrossClue
	} else {
		startClue = (*navigator.grid)[state.startRow][state.startCol].downClue
	}
	if halter.CheckInitialSquare() && halter.Halt(&navigator, state) {
		return
	}
	deltas := navigator.getDeltas()
	for {
		if state.didWrap && state.row == state.startRow && state.col == state.startCol {
			break
		} 
		nextRow, nextCol := state.row + deltas.dr, state.col + deltas.dc
		if grid.isVisitable(nextRow, nextCol) {
			state.row = nextRow
			state.col = nextCol
		} else {
			_, ok := halter.(ClueChangeHalter) 
			moveToStartOfClue := ok || (!ok && navigator.direction == Forward)
			navigator.moveToNextClue(state, moveToStartOfClue)
			logger.Debug(fmt.Sprintf("Moved to cell [%d, %d]", state.col, state.row))
			if (navigator.orientation == Horizontal && startClue != (*navigator.grid)[state.row][state.col].acrossClue) ||
				(navigator.orientation == Vertical && startClue != (*navigator.grid)[state.row][state.col].downClue) {
				state.didChangeClue = true
			}
		} 
		if halter.Halt(&navigator, state) {
			return
		}
	}
} 

func (navigator Navigator) moveToNextValidCardinal(state *NavigationState) {
	deltas := navigator.getDeltas()
	grid := *navigator.grid
	for ok := true; ok; ok = !grid.isVisitable(state.row, state.col) {
		nextRow, nextCol := state.row + deltas.dr, state.col + deltas.dc
		if nextRow < 0 {
			nextRow = len(grid) - 1
		} else if nextRow == len(grid) {
			nextRow = 0
		}
		if nextCol < 0 {
			nextCol = len(grid[0]) - 1
		} else if nextCol == len(grid[0]) {
			nextCol = 0
		}

		if navigator.orientation == Horizontal {
			if navigator.direction == Forward {
				if nextCol < state.col {
					nextRow++
					if nextRow == len(grid) {
						nextRow = 0
						state.didWrap = true
					}
				}
			} else {
				if nextCol > state.col {
					nextRow--
					if nextRow == -1 {
						nextRow = len(grid) - 1
						state.didWrap = true
					}
				}
			}
		} else {
			if navigator.direction == Forward {
				if nextRow < state.row {
					nextCol++
					if nextCol == len(grid[0]) {
						nextCol = 0
						state.didWrap = true
					}
				}
			} else {
				if nextRow > state.row {
					nextCol--
					if nextCol == -1 {
						nextCol = len(grid[0]) - 1
						state.didWrap = true
					}
				}
			}
		}
		state.row, state.col = nextRow, nextCol
	}
	state.didChangeClue = true
}

func (navigator Navigator) moveToNextClue(state *NavigationState, moveToClueStart bool) {
	startRow, startCol := state.row, state.col
	grid := *navigator.grid
	currentClueCell := grid[startRow][startCol]
	var currentClue *puzzle.Clue
	var nextClue *puzzle.Clue
	if navigator.orientation == Horizontal {
		currentClue = currentClueCell.acrossClue
		if navigator.direction == Forward {
			nextClue = currentClueCell.nextAcross 
			if nextClue.Num < currentClue.Num {
				state.didWrap = true
			}
		} else {
			nextClue = currentClueCell.prevAcross
			if nextClue.Num > currentClue.Num {
				state.didWrap = true
			}
		}
	} else {
		currentClue = currentClueCell.downClue
		if navigator.direction == Forward {
			nextClue = currentClueCell.nextDown
			if currentClue.Num > nextClue.Num {
				state.didWrap = true
			}
		} else {
			nextClue = currentClueCell.prevDown
			if currentClue.Num < nextClue.Num {
				state.didWrap = true
			}
		}
	}
	if moveToClueStart {
		state.row, state.col = nextClue.StartRow, nextClue.StartCol
	} else {
		state.row, state.col = nextClue.EndRow, nextClue.EndCol
	}
	state.didChangeClue = true
}

func (grid NavigationGrid) isVisitable(row int, col int) bool {
	return row < len(grid) &&
		col < len(grid[0]) &&
		row >= 0 &&
		col >= 0 &&
		grid[row][col].content != "."
}

func (n Navigator) getDeltas() NavigationDeltas {
	var deltas NavigationDeltas
	if n.orientation == Horizontal {
		if n.direction == Forward {
			deltas.dc = 1
		} else {
			deltas.dc = -1
		}
	} else {
		if n.direction == Forward {
			deltas.dr = 1
		} else {
			deltas.dr = -1
		}
	}
	return deltas
}

func (n Navigator) String() string {
	return fmt.Sprintf("{o: %v, d: %v, iterMode: %v, halters: %v}", n.orientation, n.direction, n.iterMode, n.halters)
}

func (c Cell) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Cell{val:%s", c.content))
 	if c.acrossClue != nil {
		sb.WriteString(fmt.Sprintf(", across: %v", c.acrossClue))
 	}
 	if c.downClue != nil {
		sb.WriteString(fmt.Sprintf(", down: %v", c.downClue))
 	}
 	if c.nextAcross != nil {
		sb.WriteString(fmt.Sprintf(", nextAcross: %v, ", c.nextAcross))
 	}
 	if c.nextDown != nil {
		sb.WriteString(fmt.Sprintf(", nextDown: %v", c.nextDown))
 	}
 	if c.prevAcross != nil {
		sb.WriteString(fmt.Sprintf(", prevAcross: %v", c.prevAcross))
	}
 	if c.prevDown != nil {
		sb.WriteString(fmt.Sprintf(", prevDown: %v", c.prevDown))
 	}
 	sb.WriteString("}")
	return sb.String()
}