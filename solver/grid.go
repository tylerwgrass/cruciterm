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
	navigator *Navigator	
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
	navigator := NewNavigator(basicGrid, puz)
	grid := navigator.grid
	currentAcrossClue = (*grid)[initialY][initialX].acrossClue
	currentDownClue = (*grid)[initialY][initialX].downClue
	return gridModel{
		navigator: navigator,
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

			newCursorX, newCursorY := m.cursorX, m.cursorY
			halters := make([]IHalter, 0, 1)
			defaultHalter := makeHalter(ValidSquare, false)
			halters = append(halters, defaultHalter)
			if prefs.GetBool(prefs.JumpToEmptySquare) {
				halters = append(halters, makeHalter(EmptySquare, true))
			}

			var didWrap bool
			if ok, _ := regexp.MatchString(`^[a-zA-Z0-9]$`, msg.String()); ok {
				(*m.navigator.grid)[m.cursorY][m.cursorX].content = strings.ToUpper(string(msg.Runes[0]))
				m.cursorY, m.cursorX, didWrap = m.navigator.
					withOrientation(m.navOrientation).
					withHalters(halters).
					advanceCursor(m.cursorX, m.cursorY)
				if didWrap && prefs.GetBool(prefs.SwapCursorOnGridWrap) {
					m.changeNavOrientation()	
				} 
				break
			}

			switch msg.String() {
			// TODO: moves to start of prev clue instead of end
			case "backspace":
				(*m.navigator.grid)[m.cursorY][m.cursorX].content = "-"
				newCursorY, newCursorX, didWrap = m.navigator.
					withOrientation(m.navOrientation).
					withDirection(Reverse).
					advanceCursor(m.cursorX, m.cursorY)
			case " ":
				m.changeNavOrientation()
			case "shift+tab":
				newCursorY, newCursorX, didWrap = m.navigator.
					withOrientation(m.navOrientation).
					withDirection(Reverse).
					advanceClue(m.cursorX, m.cursorY)
			case "tab":
				newCursorY, newCursorX, didWrap = m.navigator.
					withOrientation(m.navOrientation).
					advanceClue(m.cursorX, m.cursorY)
			case "up":
				newCursorY, newCursorX, didWrap = m.navigator.
					withOrientation(Vertical).	
					withDirection(Reverse).
					withIterMode(Cardinal).
					advanceCursor(m.cursorX, m.cursorY)
			case "down":
				newCursorY, newCursorX, didWrap = m.navigator.
					withOrientation(Vertical).	
					withIterMode(Cardinal).
					advanceCursor(m.cursorX, m.cursorY)
			case "left":
				newCursorY, newCursorX, didWrap = m.navigator.
					withDirection(Reverse).
					withIterMode(Cardinal).
					advanceCursor(m.cursorX, m.cursorY)
			case "right":
				newCursorY, newCursorX, didWrap = m.navigator.
					withIterMode(Cardinal).
					advanceCursor(m.cursorX, m.cursorY)
			}
			m.cursorX, m.cursorY = newCursorX, newCursorY
			if didWrap && prefs.GetBool(prefs.SwapCursorOnGridWrap) {
				m.changeNavOrientation()	
			} 
    }
		m.navigator.resetNavigatorOptions()
		m.validateSolution()
		currentAcrossClue = (*m.navigator.grid)[m.cursorY][m.cursorX].acrossClue
		currentDownClue = (*m.navigator.grid)[m.cursorY][m.cursorX].downClue
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
	for i, row := range *m.navigator.grid {
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
		if i < len(*m.navigator.grid) - 1 {
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
	grid := *m.navigator.grid
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
