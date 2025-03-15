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

type NavigationDeltas struct {
	dr int
	dc int
}

type NavigationGrid [][]Cell
var NavGrid NavigationGrid

type NavHalter interface {
	Halt(*NavigationGrid, int, int) bool
}

func NewNavigationGrid(puzzleGrid [][]string, puz *puzzle.PuzzleDefinition) *NavigationGrid {
	NavGrid = make([][]Cell, len(puzzleGrid))
	acrosses := puzzle.AcrossClues
	downs := puzzle.DownClues
	currentAcrossIndex := 0
	currentDownIndex := 0
	prevAcross := acrosses[len(acrosses) - 1]
	nextAcross := acrosses[(currentAcrossIndex + 1) % len(acrosses)]
	prevDown := downs[len(downs) - 1]
	nextDown := downs[(currentDownIndex + 1) % len(downs)]

	for row := range puz.NumRows {
		NavGrid[row] = make([]Cell, len(puzzleGrid[0]))
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
				if row > 0 && NavGrid[row - 1][col].content != "." {
					cell.downClue = NavGrid[row - 1][col].downClue
					cell.nextDown = NavGrid[row - 1][col].nextDown
					cell.prevDown = NavGrid[row - 1][col].prevDown
				} else {
					cell.downClue = downs[currentDownIndex].Num
				}
			}
			NavGrid[row][col] = cell
		}
	}
	return &NavGrid
}

type ValidSquareHalter func(g *NavigationGrid, i, j int) bool
func (h ValidSquareHalter) Halt(g *NavigationGrid, i, j int) bool {
	return (*g)[i][j].content != "."
}

type BlackSquareHalter func(g *NavigationGrid, i, j int) bool
func (h BlackSquareHalter) Halt(g *NavigationGrid, i, j int) bool {
	return (*g)[i][j].content == "."
}

type EmptySquareHalter func(g *NavigationGrid, i, j int) bool
func (h EmptySquareHalter) Halt(g *NavigationGrid, i, j int) bool {
	return (*g)[i][j].content == "-"
}

func (grid NavigationGrid) advanceHorizontal(initX, initY, delta int, halter NavHalter) (int, int, bool) {
	shouldWrapGrid := prefs.GetBool(prefs.WrapAtEndOfGrid)
	hasWrappedGrid := false

	row, col := initY, initX
	col += delta
	for row < len(grid) && row >= 0 {		
		for i := col; i >= 0 && i < len(grid[0]); i += delta {
			if hasWrappedGrid && i == initX && row == initY {
				return initX, initY, true
			}
			if halter.Halt(&grid, row, i) {
				return i, row, hasWrappedGrid
			}
		}
		if delta == -1 {
			row--
			col = len(grid[0]) - 1
			if shouldWrapGrid && row == -1 {
				hasWrappedGrid = true
				row = len(grid) - 1
			}
		} else {
			row++
			col = 0
			if shouldWrapGrid && row == len(grid) {
				hasWrappedGrid = true
				row = 0
			}
		}
	}
	return initX, initY, false
}

func (grid NavigationGrid) advanceVertical(initX, initY, delta int, halter NavHalter) (int, int, bool) {
	shouldWrapGrid := prefs.GetBool(prefs.WrapAtEndOfGrid)
	hasWrappedGrid := false

	row, col := initY, initX
	row += delta
	for col < len(grid[0]) && col >= 0 {
		for i := row; i >= 0 && i < len(grid); i += delta {
			if hasWrappedGrid && i == initX && row == initY {
				return initX, initY, true
			}
			if halter.Halt(&grid, i, col) {
				return col, i, hasWrappedGrid
			}
		}
		if delta == -1 {
			col--
			row = len(grid)-1
			if shouldWrapGrid && col == -1 {
				hasWrappedGrid = true
				col = len(grid[0]) - 1
			}
		} else {
			col++
			row = 0
			if shouldWrapGrid && col == len(grid[0]) {
				hasWrappedGrid = true
				col = 0
			}
		}
	}
	return initX, initY, false
}

func (grid NavigationGrid) advanceCursor(startCol, startRow int, or Orientation, dir Direction, h NavHalter, iterMode IterationMode) (int, int, bool) {
	if iterMode == Cardinal {
		return grid.advanceCursorWithNavigator(startCol, startRow, or, dir, h)
	} else {
		return grid.iterateClues(startRow, startCol, or, dir, h)
	}
}

func (grid NavigationGrid) advanceCursorWithNavigator(startX, startY int, or Orientation, dir Direction, halter NavHalter) (int, int, bool) {
	if or == Horizontal {
		return grid.advanceHorizontal(startX, startY, int(dir), halter)
	} else {
		return grid.advanceVertical(startX, startY, int(dir), halter)
	}
}

func (grid NavigationGrid) advanceClue(startX, startY int, or Orientation, dir Direction) (int, int, bool) {
	currentCell := grid[startY][startX]
	var nextClueNum int
	didWrap := false
	swapOnWrap := prefs.GetBool(prefs.SwapCursorOnGridWrap)
	if or == Horizontal {
		if dir == Forward {
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
		if dir == Forward {
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
	return nextClue.StartX, nextClue.StartY, didWrap
}

func (grid NavigationGrid) iterateClues(startRow int, startCol int, or Orientation, dir Direction, halter NavHalter) (row int, col int, didWrap bool) {
	currentRow, currentCol := startRow, startCol
	didWrap = false
	deltas := getDeltas(or, dir)
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
			currentRow, currentCol, didWrapGrid = grid.getNextClueLocation(currentRow, currentCol, or, dir)
			didWrap = didWrap || didWrapGrid
		} 
		if halter.Halt(&grid, currentRow, currentCol) {
			return currentRow, currentCol, didWrap
		}
	}

	return startRow, startCol, false
} 

func (grid NavigationGrid) getNextClueLocation(startRow int, startCol int, or Orientation, dir Direction) (row int, col int, didWrap bool) {
	currentClue := grid[startRow][startCol]
	var nextClueNum int
	if or == Horizontal {
		if dir == Forward {
			nextClueNum = currentClue.nextAcross 
			didWrap = nextClueNum < currentClue.acrossClue
		} else {
			nextClueNum = currentClue.prevAcross
			didWrap = nextClueNum > currentClue.acrossClue
		}
	} else {
		if dir == Forward {
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

func getDeltas(or Orientation, dir Direction) NavigationDeltas {
	var deltas NavigationDeltas
	if or == Horizontal {
		if dir == Forward {
			deltas.dc = 1
		} else {
			deltas.dc = -1
		}
	} else {
		if dir == Forward {
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