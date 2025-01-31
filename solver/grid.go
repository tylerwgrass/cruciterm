package solver

import (
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tylerwgrass/cruciterm/puzzle"
)

type Direction int 

const (
	Horizontal Direction = iota
	Vertical
)

type gridModel struct {
	Grid [][]string
	solution string
	solved bool
	cursorX int
	cursorY int
}

func initGridModel(puz *puzzle.PuzzleDefinition) gridModel {
	grid := make([][]string, puz.NumRows)
	var initialX int
	var initialY int
	startFound := false
	solved := true
	for i := range puz.NumRows {
		grid[i] = make([]string, puz.NumCols)
		for j := range puz.NumCols {
			grid[i][j] = string(puz.CurrentState[i*puz.NumCols + j])
			if grid[i][j] != "." && !startFound {
				startFound = true
				initialX = j
				initialY = i
			}
			if grid[i][j] != string(puz.Answer[i*puz.NumCols+j]) {
				solved = false
			}
		}
	}

	return gridModel{
		Grid: grid,
		solved: solved,
		solution: puz.Answer,
		cursorX: initialX,
		cursorY: initialY,
	}
}

func (m gridModel) Init() tea.Cmd {
	return nil
}

func (m gridModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
    case tea.KeyMsg:
			if m.solved {
				break
			}

			if ok, _ := regexp.MatchString(`^[a-zA-Z0-9]$`, msg.String()); ok {
				m.Grid[m.cursorY][m.cursorX] = strings.ToUpper(string(msg.Runes[0]))
				m.advanceCursor(Horizontal, 1, true) 
				break
			}

			switch msg.String() {
			case "up":
					m.advanceCursor(Vertical, -1, true)
			case "down":
					m.advanceCursor(Vertical, 1, true)
			case "left":
					m.advanceCursor(Horizontal, -1, true)
			case "right":
					m.advanceCursor(Horizontal, 1, true)
			}
    }
		m.validateSolution()
    return m, nil
}

func (m gridModel) View() string {
	var sb strings.Builder
	for i, row := range m.Grid {
		sb.WriteString(" ")
		for j, cell := range row {
			if i == m.cursorY && j == m.cursorX {
				if m.solved {
					sb.WriteString(cell + " ")
				} else {
					sb.WriteString("> ")
				}
				continue
			}
			switch cell {
			case ".":
				sb.WriteString("■ ")
			case "-":
				sb.WriteString("  ")
			default:
				sb.WriteString(cell + " ")
			} 
		}
		if i < len(m.Grid)-1 {
			sb.WriteString("\n")
		}
	}

	return baseStyle.Render(sb.String())
}

func (m* gridModel) advanceCursor(dir Direction, delta int, wrap bool) {
	if dir == Horizontal {
		m.advanceHorizontal(delta, wrap)
	} else {
		m.advanceVertical(delta, wrap)
	}
}

func (m *gridModel) advanceHorizontal(delta int, wrap bool) (int, int) {
	row, col := m.cursorY, m.cursorX
	col += delta
	for row < len(m.Grid) {		
		for i := col; i >= 0 && i < len(m.Grid[0]); i += delta {
			if m.Grid[row][i] != "." {
				m.cursorX = i
				m.cursorY = row
				return m.cursorX, m.cursorY
			}
		}
		if !wrap {
			break
		}
		row++
		col = 0
	}
	return m.cursorX, m.cursorY
}

func (m *gridModel) advanceVertical(delta int, wrap bool) (int, int) {
	row, col := m.cursorY, m.cursorX
	row += delta
	for col < len(m.Grid[0]) {
		for i := row; i >= 0 && i < len(m.Grid); i += delta {
			if m.Grid[i][col] != "." {
				m.cursorX = col
				m.cursorY = i
				return m.cursorX, m.cursorY
			}
		}
		if !wrap {
			break
		}
		col++
		row = 0
	}
	return m.cursorX, m.cursorY
}

func (m *gridModel) validateSolution() {
	grid := m.Grid
	numRows := len(grid)
	numCols := len(grid[0])
	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			if (grid[i][j] != string(m.solution[(i*numCols)+j])) {
				m.solved = false
				return
			}
		}
	}
	m.solved = true
}
