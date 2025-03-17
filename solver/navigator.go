package solver

import (
	"fmt"

	prefs "github.com/tylerwgrass/cruciterm/preferences"
	"github.com/tylerwgrass/cruciterm/puzzle"
)

type Cell struct {
	content string
	acrossClue int
	downClue int
	nextAcross int
	nextDown int
	prevAcross int
	prevDown int
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
			if row == nextAcross.StartY && col == nextAcross.StartX {
				prevAcross = acrosses[currentAcrossIndex]
				currentAcrossIndex = (currentAcrossIndex + 1) % len(acrosses)
				nextAcross = acrosses[(currentAcrossIndex + 1) % len(acrosses)]
			}
			if row == nextDown.StartY && col == nextDown.StartX {
				prevDown = downs[currentDownIndex]
				currentDownIndex = (currentDownIndex + 1) % len(downs)
				nextDown = downs[(currentDownIndex + 1) % len(downs)]
			}
			cell := Cell{
				content: puzzleGrid[row][col], 
			}
			if cell.content != "." {
				cell.prevAcross = prevAcross.Num
				cell.nextAcross = nextAcross.Num
				cell.prevDown = prevDown.Num
				cell.nextDown = nextDown.Num
				cell.acrossClue = acrosses[currentAcrossIndex].Num
				if row > 0 && navGrid[row - 1][col].content != "." {
					cell.downClue = navGrid[row - 1][col].downClue
					cell.nextDown = navGrid[row - 1][col].nextDown
					cell.prevDown = navGrid[row - 1][col].prevDown
				} else {
					cell.downClue = downs[currentDownIndex].Num
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

func (navigator Navigator) advanceCursor(startCol, startRow int) (int, int, bool) {
	didWrap := false
	row, col := startRow, startCol
	for _, halter := range(navigator.halters) {
		if navigator.iterMode == Cardinal {
			iterRow, iterCol, didCurrentWrap := navigator.iterateCardinal(row, col, halter)
			row, col, didWrap = iterRow, iterCol, didWrap || didCurrentWrap
		} else {
			iterRow, iterCol, didCurrentWrap := navigator.iterateClues(row, col, halter)
			row, col, didWrap = iterRow, iterCol, didWrap || didCurrentWrap
		}
	}
	return row, col, didWrap
}

func (navigator Navigator) advanceClue(startX, startY int) (int, int, bool) {
	grid := *navigator.grid
	currentCell := grid[startY][startX]
	var nextClueNum int
	didWrap := false
	swapOnWrap := prefs.GetBool(prefs.SwapCursorOnGridWrap)
	if navigator.orientation == Horizontal {
		if navigator.direction == Forward {
			nextClueNum = currentCell.nextAcross
			didWrap = nextClueNum < currentCell.acrossClue
		} else {
			nextClueNum = currentCell.prevAcross
			didWrap = nextClueNum > currentCell.acrossClue
			if didWrap && swapOnWrap {
				nextClueNum = currentCell.prevDown
			}
		}
	} else {
		if navigator.direction == Forward {
			nextClueNum = currentCell.nextDown
			didWrap = nextClueNum < currentCell.downClue
		} else {
			nextClueNum = currentCell.prevDown
			didWrap = nextClueNum > currentCell.downClue
			if didWrap && swapOnWrap {
				nextClueNum = currentCell.prevAcross
			}
		}
	}
	nextClue := puzzle.Clues[nextClueNum]
	return nextClue.StartY, nextClue.StartX, didWrap
}

func (navigator Navigator) iterateCardinal(startRow, startCol int, halter IHalter) (row, col int, didWrap bool) {
	grid := *navigator.grid
	if halter.CheckInitialSquare() && halter.Halt(&grid, startRow, startCol) {
		return startRow, startCol, false
	}
	currentRow, currentCol := startRow, startCol
	didWrap = false
	deltas := navigator.getDeltas()
	for {
		if didWrap && currentRow == startRow && currentCol == startCol {
			break
		}
		nextRow, nextCol := currentRow + deltas.dr, currentCol + deltas.dc
		if grid.isVisitable(nextRow, nextCol) {
			currentRow = nextRow
			currentCol = nextCol
		} else {
			var didWrapGrid bool
			currentRow, currentCol, didWrapGrid = navigator.getNextCardinalCell(currentRow, currentCol)
			didWrap = didWrap || didWrapGrid
		} 
		if halter.Halt(&grid, currentRow, currentCol) {
			return currentRow, currentCol, didWrap
		}
	}
	return startRow, startCol, false
}

func (navigator Navigator) iterateClues(startRow int, startCol int, halter IHalter) (row int, col int, didWrap bool) {
	grid := *navigator.grid
	if halter.CheckInitialSquare() && halter.Halt(&grid, startRow, startCol) {
		return startRow, startCol, false
	}
	currentRow, currentCol := startRow, startCol
	didWrap = false
	deltas := navigator.getDeltas()
	for {
		if didWrap && currentRow == startRow && currentCol == startCol {
			break
		} 
		nextRow, nextCol := currentRow + deltas.dr, currentCol + deltas.dc
		if grid.isVisitable(nextRow, nextCol) {
			currentRow = nextRow
			currentCol = nextCol
		} else {
			var didWrapGrid bool
			currentRow, currentCol, didWrapGrid = navigator.getNextClueLocation(currentRow, currentCol)
			didWrap = didWrap || didWrapGrid
		} 
		if halter.Halt(&grid, currentRow, currentCol) {
			return currentRow, currentCol, didWrap
		}
	}

	return startRow, startCol, false
} 

func (navigator Navigator) getNextCardinalCell(startRow, startCol int) (row, col int, didWrap bool) {
	row, col, didWrap = startRow, startCol, false
	deltas := navigator.getDeltas()
	grid := *navigator.grid
	for ok := true; ok; ok = !grid.isVisitable(row, col) {
		nextRow, nextCol := row + deltas.dr, col + deltas.dc
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
				if nextCol < col {
					nextRow++
					if nextRow == len(grid) {
						nextRow = 0
						didWrap = true
					}
				}
			} else {
				if nextCol > col {
					nextRow--
					if nextRow == -1 {
						nextRow = len(grid) - 1
						didWrap = true
					}
				}
			}
		} else {
			if navigator.direction == Forward {
				if nextRow < row {
					nextCol++
					if nextCol == len(grid[0]) {
						nextCol = 0
						didWrap = true
					}
				}
			} else {
				if nextRow > row {
					nextCol--
					if nextCol == -1 {
						nextCol = len(grid[0]) - 1
						didWrap = true
					}
				}
			}
		}
		row, col = nextRow, nextCol
	}
	return row, col, didWrap
}

func (navigator Navigator) getNextClueLocation(startRow int, startCol int) (row int, col int, didWrap bool) {
	grid := *navigator.grid
	currentClue := grid[startRow][startCol]
	var nextClueNum int
	if navigator.orientation == Horizontal {
		if navigator.direction == Forward {
			nextClueNum = currentClue.nextAcross 
			didWrap = nextClueNum < currentClue.acrossClue
		} else {
			nextClueNum = currentClue.prevAcross
			didWrap = nextClueNum > currentClue.acrossClue
		}
	} else {
		if navigator.direction == Forward {
			nextClueNum = currentClue.nextDown
			didWrap = nextClueNum < currentClue.downClue
		} else {
			nextClueNum = currentClue.prevDown
			didWrap = nextClueNum > currentClue.downClue
		}
	}
	nextClue := puzzle.Clues[nextClueNum]
	row = nextClue.StartY
	col = nextClue.StartX
	return 
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

func (c Cell) ToString() string {
	return fmt.Sprintf("{val:%s, acrossClueNum:%d, downClueNum:%d, nextAcrossNum:%d, nextDownNum:%d, prevAcrossNum:%d, prevDownNum:%d}\n",
		c.content,
		c.acrossClue,
		c.downClue,
		c.nextAcross,
		c.nextDown,
		c.prevAcross,
		c.prevDown,
	)
}