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

type gridModel struct {
	Grid [][]string 
	navGrid *NavigationGrid
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
	navGrid := NewNavigationGrid(grid, puz)
	currentAcrossClue = (*navGrid)[initialY][initialX].acrossClue
	currentDownClue = (*navGrid)[initialY][initialX].downClue
	return gridModel{
		Grid: grid,
		navGrid: navGrid,
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

			var didWrap bool
			if ok, _ := regexp.MatchString(`^[a-zA-Z0-9]$`, msg.String()); ok {
				m.Grid[m.cursorY][m.cursorX] = strings.ToUpper(string(msg.Runes[0]))
				(*m.navGrid)[m.cursorY][m.cursorX].content = strings.ToUpper(string(msg.Runes[0]))
				m.cursorX, m.cursorY, didWrap = m.navGrid.advanceCursor(m.cursorX, m.cursorY, m.navOrientation, Forward); 
				if m.Grid[m.cursorY][m.cursorX] != "-" && prefs.GetBool(prefs.JumpToEmptySquare) {
					var h EmptySquareHalter
					m.cursorX, m.cursorY, didWrap = m.navGrid.advanceCursorWithNavigator(m.cursorX, m.cursorY, m.navOrientation, Forward, h)
				}
				if didWrap && prefs.GetBool(prefs.SwapCursorOnGridWrap) {
					m.changeNavOrientation()	
				} 
				break
			}

			switch msg.String() {
			case "backspace":
				m.Grid[m.cursorY][m.cursorX] = "-"
				m.navGrid.advanceCursor(m.cursorX, m.cursorY, m.navOrientation, Reverse)
			case " ":
				m.changeNavOrientation()
			case "shift+tab":
				m.cursorX, m.cursorY, didWrap = m.navGrid.advanceClue(m.cursorX, m.cursorY, m.navOrientation, Reverse)
				if didWrap && prefs.GetBool(prefs.SwapCursorOnGridWrap) {
					m.changeNavOrientation()
				}
			case "tab":
				m.cursorX, m.cursorY, didWrap = m.navGrid.advanceClue(m.cursorX, m.cursorY, m.navOrientation, Forward)
				if didWrap && prefs.GetBool(prefs.SwapCursorOnGridWrap) {
					m.changeNavOrientation()
				}
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
		currentAcrossClue = (*m.navGrid)[m.cursorY][m.cursorX].acrossClue
		currentDownClue = (*m.navGrid)[m.cursorY][m.cursorX].downClue
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
	if swapCursor && m.navOrientation != o {
		m.changeNavOrientation()
	} else {
		m.cursorX, m.cursorY, _ = m.navGrid.advanceCursor(m.cursorX, m.cursorY, o, d)
	}
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
