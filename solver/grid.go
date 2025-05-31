package solver

import (
	"regexp"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	prefs "github.com/tylerwgrass/cruciterm/preferences"
	"github.com/tylerwgrass/cruciterm/puzzle"
	"github.com/tylerwgrass/cruciterm/theme"
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
	navigator      *Navigator
	solution       string
	solved         bool
	cursorX        int
	cursorY        int
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
			basicGrid[i][j] = string(puz.CurrentState[i*puz.NumCols+j])
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
		navigator:      navigator,
		solved:         solved,
		solution:       puz.Answer,
		cursorX:        initialX,
		cursorY:        initialY,
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

		navStates := make([]NavigationState, 0)
		navStates = append(navStates, NavigationState{row: m.cursorY, col: m.cursorX, startRow: m.cursorY, startCol: m.cursorX})
		var didWrap bool
		halters := make([]IHalter, 0, 1)
		defaultHalter := makeHalter(ValidSquare, false)
		halters = append(halters, defaultHalter)
		if prefs.GetBool(prefs.JumpToEmptySquare) {
			halters = append(halters, makeHalter(EmptySquare, true))
		}

		if ok, _ := regexp.MatchString(`^[a-zA-Z0-9]$`, msg.String()); ok {
			(*m.navigator.grid)[m.cursorY][m.cursorX].content = strings.ToUpper(string(msg.String()[0]))
			navStates = m.navigator.
				withOrientation(m.navOrientation).
				withHalters(halters).
				advanceCursor(m.cursorX, m.cursorY)
			endNavState := navStates[len(navStates)-1]
			m.cursorX, m.cursorY = endNavState.col, endNavState.row
			didWrap = slices.ContainsFunc(navStates, func(ns NavigationState) bool {
				return ns.didWrap
			})
			if didWrap && prefs.GetBool(prefs.SwapCursorOnGridWrap) && endNavState.haltedOnMatch {
				m.changeNavOrientation()
			}
			break
		}

		switch {
		case key.Matches(msg, keys.Delete):
			(*m.navigator.grid)[m.cursorY][m.cursorX].content = "-"
			navStates = m.navigator.
				withOrientation(m.navOrientation).
				withMoveDirection(Reverse).
				withJumpDirection(Reverse).
				withJumpLocation(ClueEnd).
				advanceCursor(m.cursorX, m.cursorY)
		case key.Matches(msg, keys.ToggleDirection):
			m.changeNavOrientation()
		case key.Matches(msg, keys.PrevClue):
			halters = make([]IHalter, 0, 2)
			halters = append(halters, makeHalter(ClueChange, false))
			if prefs.GetBool(prefs.JumpToEmptySquare) {
				halters = append(halters, makeHalter(EmptySquare, true))
			}
			navStates = m.navigator.
				withOrientation(m.navOrientation).
				withJumpDirection(Reverse).
				withHalters(halters).
				advanceCursor(m.cursorX, m.cursorY)
		case key.Matches(msg, keys.NextClue):
			halters = make([]IHalter, 0, 2)
			halters = append(halters, makeHalter(ClueChange, false))
			if prefs.GetBool(prefs.JumpToEmptySquare) {
				halters = append(halters, makeHalter(EmptySquare, true))
			}
			navStates = m.navigator.
				withOrientation(m.navOrientation).
				withHalters(halters).
				advanceCursor(m.cursorX, m.cursorY)
		case key.Matches(msg, keys.Up):
			if prefs.GetBool(prefs.SwapCursorOnDirectionChange) && m.navOrientation != Vertical {
				m.changeNavOrientation()
				break
			}
			navStates = m.navigator.
				withOrientation(Vertical).
				withMoveDirection(Reverse).
				withIterMode(Cardinal).
				advanceCursor(m.cursorX, m.cursorY)
		case key.Matches(msg, keys.Down):
			if prefs.GetBool(prefs.SwapCursorOnDirectionChange) && m.navOrientation != Vertical {
				m.changeNavOrientation()
				break
			}
			navStates = m.navigator.
				withOrientation(Vertical).
				withIterMode(Cardinal).
				advanceCursor(m.cursorX, m.cursorY)
		case key.Matches(msg, keys.Left):
			if prefs.GetBool(prefs.SwapCursorOnDirectionChange) && m.navOrientation != Horizontal {
				m.changeNavOrientation()
				break
			}
			navStates = m.navigator.
				withMoveDirection(Reverse).
				withIterMode(Cardinal).
				advanceCursor(m.cursorX, m.cursorY)
		case key.Matches(msg, keys.Right):
			if prefs.GetBool(prefs.SwapCursorOnDirectionChange) && m.navOrientation != Horizontal {
				m.changeNavOrientation()
				break
			}
			navStates = m.navigator.
				withIterMode(Cardinal).
				advanceCursor(m.cursorX, m.cursorY)
		}
		endNavState := navStates[len(navStates)-1]
		m.cursorX, m.cursorY = endNavState.col, endNavState.row
		didWrap = slices.ContainsFunc(navStates, func(ns NavigationState) bool {
			return ns.didWrap
		})
		if didWrap && prefs.GetBool(prefs.SwapCursorOnGridWrap) && endNavState.haltedOnMatch {
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
	activeClueStyle := theme.Get().Foreground(theme.Primary())
	sb := theme.NewThemedStringBuilder(theme.Get())
	var cursor string
	if m.navOrientation == Horizontal {
		cursor = ">"
	} else {
		cursor = "v"
	}

	for i, row := range *m.navigator.grid {
		sb.WriteString(" ")
		for j, cell := range row {
			if i == m.cursorY && j == m.cursorX && !m.solved {
				sb.WriteStyledString(cursor+" ", activeClueStyle)
				continue
			}
			switch cell.content {
			case ".":
				sb.WriteString("■ ")
			case "-":
				if m.isCellInActiveClue(i, j) {
					sb.WriteStyledString("_ ", activeClueStyle)
				} else {
					sb.WriteString("  ")
				}
			default:
				if m.isCellInActiveClue(i, j) {
					sb.WriteStyledString(cell.content+" ", activeClueStyle)
				} else {
					sb.WriteString(cell.content + " ")
				}
			}
		}
		if i < len(*m.navigator.grid)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func (m gridModel) isCellInActiveClue(row, col int) bool {
	return (m.navOrientation == Horizontal &&
		col >= currentAcrossClue.StartCol &&
		col <= currentAcrossClue.EndCol &&
		row == m.cursorY) ||
		(m.navOrientation == Vertical &&
			row >= currentDownClue.StartRow &&
			row <= currentDownClue.EndRow &&
			col == m.cursorX)
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
			if grid[i][j].content != string(m.solution[(i*numCols)+j]) {
				m.solved = false
				return
			}
		}
	}
	m.solved = true
}
