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
	grid *NavigationGrid
	solution string
	solved bool
	cursorX int
	cursorY int
	navOrientation Orientation
}

func initGridModel(puz *puzzle.PuzzleDefinition) gridModel {
	basicGrid := make([][]string, puz.NumRows)
	var initialX int
	var initialY int
	startFound := false
	solved := true
	for i := range puz.NumRows {
		basicGrid[i] = make([]string, puz.NumCols)
		for j := range puz.NumCols {
			basicGrid[i][j] = string(puz.CurrentState[i*puz.NumCols + j])
			if basicGrid[i][j] != "." && !startFound {
				startFound = true
				initialX = j
				initialY = i
			}
			if basicGrid[i][j] != string(puz.Answer[i*puz.NumCols+j]) {
				solved = false
			}
		}
	}
	grid := NewNavigationGrid(basicGrid, puz)
	currentAcrossClue = (*grid)[initialY][initialX].acrossClue
	currentDownClue = (*grid)[initialY][initialX].downClue
	return gridModel{
		grid: grid,
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
				(*m.grid)[m.cursorY][m.cursorX].content = strings.ToUpper(string(msg.Runes[0]))
				var h ValidSquareHalter
				m.cursorY, m.cursorX, didWrap = m.grid.advanceCursor(m.cursorX, m.cursorY, m.navOrientation, Forward, h, Clues); 
				if (*m.grid)[m.cursorY][m.cursorX].content != "-" && prefs.GetBool(prefs.JumpToEmptySquare) {
					var e EmptySquareHalter
					m.cursorY, m.cursorX, didWrap = m.grid.advanceCursor(m.cursorX, m.cursorY, m.navOrientation, Forward, e, Clues)
				}
				if didWrap && prefs.GetBool(prefs.SwapCursorOnGridWrap) {
					m.changeNavOrientation()	
				} 
				break
			}

			switch msg.String() {
			// TODO: moves to start of prev clue instead of end
			case "backspace":
				(*m.grid)[m.cursorY][m.cursorX].content = "-"
				var h ValidSquareHalter
				m.cursorY, m.cursorX, _ = m.grid.advanceCursor(m.cursorX, m.cursorY, m.navOrientation, Reverse, h, Clues)
			case " ":
				m.changeNavOrientation()
			case "shift+tab":
				m.cursorX, m.cursorY, didWrap = m.grid.advanceClue(m.cursorX, m.cursorY, m.navOrientation, Reverse)
				if didWrap && prefs.GetBool(prefs.SwapCursorOnGridWrap) {
					m.changeNavOrientation()
				}
			case "tab":
				m.cursorX, m.cursorY, didWrap = m.grid.advanceClue(m.cursorX, m.cursorY, m.navOrientation, Forward)
				if didWrap && prefs.GetBool(prefs.SwapCursorOnGridWrap) {
					m.changeNavOrientation()
				}
			case "up":
				var h ValidSquareHalter
				m.cursorY, m.cursorX, _ = m.grid.advanceCursor(m.cursorX, m.cursorY, Vertical, Reverse, h, Cardinal)
			case "down":
				var h ValidSquareHalter
				m.cursorY, m.cursorX, _ = m.grid.advanceCursor(m.cursorX, m.cursorY, Vertical, Forward, h, Cardinal)
			case "left":
				var h ValidSquareHalter
				m.cursorY, m.cursorX, _ = m.grid.advanceCursor(m.cursorX, m.cursorY, Horizontal, Reverse, h, Cardinal)
			case "right":
				var h ValidSquareHalter
				m.cursorY, m.cursorX, _ = m.grid.advanceCursor(m.cursorX, m.cursorY, Horizontal, Forward, h, Cardinal)
			}
    }
		m.validateSolution()
		currentAcrossClue = (*m.grid)[m.cursorY][m.cursorX].acrossClue
		currentDownClue = (*m.grid)[m.cursorY][m.cursorX].downClue
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
	for i, row := range *m.grid {
		sb.WriteString(" ")
		for j, cell := range row {
			if i == m.cursorY && j == m.cursorX {
				if m.solved {
					sb.WriteString(cell.content + " ")
				} else {
					sb.WriteString(string(cursor) + " ")
				}
				continue
			}
			switch cell.content {
			case ".":
				sb.WriteString("■ ")
			case "-":
				sb.WriteString("  ")
			default:
				sb.WriteString(cell.content + " ")
			} 
		}
		if i < len(*m.grid) - 1 {
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

func (m *gridModel) validateSolution() {
	grid := *m.grid
	numRows := len(grid)
	numCols := len(grid[0])
	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			if (grid[i][j].content != string(m.solution[(i * numCols) + j])) {
				m.solved = false
				return
			}
		}
	}
	m.solved = true
}
