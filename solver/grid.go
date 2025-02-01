package solver

import (
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	prefs "github.com/tylerwgrass/cruciterm/preferences"
	"github.com/tylerwgrass/cruciterm/puzzle"
)

type Direction int 
type Orientation int
const (
	Forward Direction = 1
	Reverse Direction = -1
)

const (
	Horizontal Orientation = iota
	Vertical
)

type Grid [][]string

type NavHalter interface {
	Halt(Grid, int, int) bool
}

type gridModel struct {
	Grid Grid 
	solution string
	solved bool
	cursorX int
	cursorY int
	navOrientation Orientation
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
		navOrientation: Horizontal,
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
				m.advanceCursor(m.navOrientation, Forward) 
				break
			}

			switch msg.String() {
			case " ":
				m.changeNavOrientation()
			case "tab":
				m.advanceClue(m.navOrientation, Forward)
			case "up":
				m.handleCardinal(Vertical, Reverse)
			case "down":
				m.handleCardinal(Vertical, Forward)
			case "left":
				m.handleCardinal(Horizontal, Reverse)
			case "right":
				m.handleCardinal(Horizontal, Forward)
			}
    }
		m.validateSolution()
    return m, nil
}

func (m gridModel) View() string {
	var sb strings.Builder
	var cursor string
	if (m.navOrientation == Horizontal) {
		cursor = ">"
	} else {
		cursor = "v" 
	}
	cursor = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(cursor)
	for i, row := range m.Grid {
		sb.WriteString(" ")
		for j, cell := range row {
			if i == m.cursorY && j == m.cursorX {
				if m.solved {
					sb.WriteString(cell + " ")
				} else {
					sb.WriteString(string(cursor) + " ")
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

func (m *gridModel) changeNavOrientation() {
	if m.navOrientation == Horizontal {
		m.navOrientation = Vertical
	} else {
		m.navOrientation = Horizontal
	}
}

func (m *gridModel) handleCardinal(o Orientation, d Direction) {
	swapCursor := prefs.GetBool(prefs.SwapCursorOnDirectionChange)
	if swapCursor && m.navOrientation != Vertical {
		m.changeNavOrientation()
	} else {
		m.advanceCursor(o, d)
	}
}

func (m *gridModel) advanceCursor(or Orientation, dir Direction) {
	var h ValidSquareHalter
	cursorX, cursorY := m.advanceCursorWithNavigator(or, dir, h)
	m.cursorX = cursorX
	m.cursorY = cursorY
}

func (m* gridModel) advanceCursorWithNavigator(or Orientation, dir Direction, halter NavHalter) (int, int) {
	if or == Horizontal {
		return m.advanceHorizontal(int(dir), halter)
	} else {
		return m.advanceVertical(int(dir), halter)
	}
}

func (m gridModel) advanceHorizontal(delta int, halter NavHalter) (int, int) {
	shouldWrapGrid := prefs.GetBool(prefs.WrapAtEndOfGrid)
	hasWrappedGrid := false

	row, col := m.cursorY, m.cursorX
	col += delta
	for row < len(m.Grid) && row >= 0 {		
		for i := col; i >= 0 && i < len(m.Grid[0]); i += delta {
			if hasWrappedGrid && i == m.cursorX && row == m.cursorY {
				return m.cursorX, m.cursorY
			}
			if halter.Halt(m.Grid, row, i) {
				return i, row
			}
		}
		if delta == -1 {
			row--
			col = len(m.Grid[0]) - 1
			if shouldWrapGrid && row == -1 {
				hasWrappedGrid = true
				row = len(m.Grid) - 1
			}
		} else {
			row++
			col = 0
			if shouldWrapGrid && row == len(m.Grid) {
				hasWrappedGrid = true
				row = 0
			}
		}
	}
	return m.cursorX, m.cursorY
}

func (m gridModel) advanceVertical(delta int, halter NavHalter) (int, int) {
	shouldWrapGrid := prefs.GetBool(prefs.WrapAtEndOfGrid)
	hasWrappedGrid := false

	row, col := m.cursorY, m.cursorX
	row += delta
	for col < len(m.Grid[0]) && col >= 0 {
		for i := row; i >= 0 && i < len(m.Grid); i += delta {
			if hasWrappedGrid && i == m.cursorX && row == m.cursorY {
				return m.cursorX, m.cursorY
			}
			if halter.Halt(m.Grid, i, col) {
				return col, i
			}
		}
		if delta == -1 {
			col--
			row = len(m.Grid)-1
			if shouldWrapGrid && col == -1 {
				hasWrappedGrid = true
				col = len(m.Grid[0]) - 1
			}
		} else {
			col++
			row = 0
			if shouldWrapGrid && col == len(m.Grid[0]) {
				hasWrappedGrid = true
				col = 0
			}
		}
	}
	return m.cursorX, m.cursorY
}

func (m *gridModel) advanceClue(or Orientation, dir Direction) {
	var validSquareHalter ValidSquareHalter
	var blackSquareHalter BlackSquareHalter
	initX, initY := m.cursorX, m.cursorY
	m.cursorX, m.cursorY = m.advanceCursorWithNavigator(or, dir, blackSquareHalter)
	nextX, nextY := m.advanceCursorWithNavigator(or, dir, validSquareHalter)
	if m.cursorX == nextX && m.cursorY == nextY {
		m.cursorX, m.cursorY = initX, initY
	} else {
		m.cursorX, m.cursorY = nextX, nextY
	}
}

type ValidSquareHalter func(g Grid, i, j int) bool
func (h ValidSquareHalter) Halt(g Grid, i, j int) bool {
	return g[i][j] != "."
}

type BlackSquareHalter func(g Grid, i, j int) bool
func (h BlackSquareHalter) Halt(g Grid, i, j int) bool {
	return g[i][j] == "."
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
