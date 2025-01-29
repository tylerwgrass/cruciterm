package solver

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tylerwgrass/cruciterm/puzzle"
)

type model struct {
	title string
	author string
	copyright string
	acrossClues map[int]string
	downClues map[int]string
	grid [][]string
	cursorX int
	cursorY int
}

func initialModel(puz *puzzle.PuzzleDefinition) model {
	grid := make([][]string, puz.NumRows)
	for i := range puz.NumRows {
		grid[i] = make([]string, puz.NumCols)
		for j := range puz.NumCols {
			grid[i][j] = string(puz.CurrentState[i*puz.NumCols + j])
		}
	}

	var initialX int
	var initialY int
	for i, char := range puz.CurrentState {
		if char != '.' {
			initialX = i % puz.NumCols
			initialY = i / puz.NumCols
			break
		}
	}

	return model{
		title: puz.Title,
		author: puz.Author,
		copyright: puz.Copyright,
		acrossClues: puz.AcrossClues,
		downClues: puz.DownClues,
		grid: grid,
		cursorX: initialX,
		cursorY: initialY,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
				case "ctrl+c":
					return m, tea.Quit

				case "up":
					if m.cursorY > 0 {
						m.getNextVerticalCell(-1)
					}

				case "down":
					if m.cursorY < len(m.grid)-1 {
						m.getNextVerticalCell(1)
					}
				case "left":
					if m.cursorX > 0 {
						m.getNextHorizontalCell(-1)
					}
				case "right":
					if m.cursorX < len(m.grid[0])-1 {
						m.getNextHorizontalCell(1)
					}
				}
    }

    return m, nil
}

func (m model) View() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\n%s\n", m.title))
	sb.WriteString(fmt.Sprintf("%s %s\n\n", m.author, m.copyright))

	for i, row := range m.grid {
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
		sb.WriteString("\n")
	}

	sb.WriteString("\nPress ctrl+c to quit.\n")
	return sb.String()
}

func (m *model) getNextHorizontalCell(dir int) {
	for i := m.cursorX + dir; i >= 0 && i < len(m.grid[0]); i += dir {
		if m.grid[m.cursorY][i] != "." {
			m.cursorX = i
			return
		}
	}
}

func (m *model) getNextVerticalCell(dir int) {
	for i := m.cursorY + dir; i >= 0 && i < len(m.grid); i += dir {
		if m.grid[i][m.cursorX] != "." {
			m.cursorY = i
			return
		}
	}
} 