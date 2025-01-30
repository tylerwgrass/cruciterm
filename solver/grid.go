package solver

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tylerwgrass/cruciterm/puzzle"
)

type gridModel struct {
	Grid [][]string
	cursorX int
	cursorY int
}

func initGridModel(puz *puzzle.PuzzleDefinition) gridModel {
	grid := make([][]string, puz.NumRows)
	var initialX int
	var initialY int
	startFound := false

	for i := range puz.NumRows {
		grid[i] = make([]string, puz.NumCols)
		for j := range puz.NumCols {
			grid[i][j] = string(puz.CurrentState[i*puz.NumCols + j])
			if grid[i][j] != "." && !startFound {
				startFound = true
				initialX = j
				initialY = i
			}
		}
	}

	return gridModel{
		Grid: grid,
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
        switch msg.String() {
				case "up":
					if m.cursorY > 0 {
						m.navigateVertical(-1)
					}

				case "down":
					if m.cursorY < len(m.Grid)-1 {
						m.navigateVertical(1)
					}

				case "left":
					if m.cursorX > 0 {
						m.navigateHorizontal(-1)
					}

				case "right":
					if m.cursorX < len(m.Grid[0])-1 {
						m.navigateHorizontal(1)
					}
				//TODO: handle ctrl chars
				default:
					m.Grid[m.cursorY][m.cursorX] = strings.ToUpper(string(msg.Runes[0]))
					m.navigateHorizontal(1) 
				}
    }

    return m, nil
}

func (m gridModel) View() string {
	var sb strings.Builder

	for i, row := range m.Grid {
		for j, cell := range row {
			if i == m.cursorY && j == m.cursorX {
				sb.WriteString("_ ")
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

func (m *gridModel) navigateHorizontal(dir int) {
	for i := m.cursorX + dir; i >= 0 && i < len(m.Grid[0]); i += dir {
		if m.Grid[m.cursorY][i] != "." {
			m.cursorX = i
			return
		}
	}
}

func (m *gridModel) navigateVertical(dir int) {
	for i := m.cursorY + dir; i >= 0 && i < len(m.Grid); i += dir {
		if m.Grid[i][m.cursorX] != "." {
			m.cursorY = i
			return
		}
	}
} 