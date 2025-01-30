package solver

import (
	"fmt"
	"os"

	"github.com/tylerwgrass/cruciterm/puzzle"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type mainModel struct {
	title string
	author string
	copyright string
	acrossClues map[int]string
	downClues map[int]string
	clues cluesModel
	grid tea.Model
	solution string
	cursorX int
	cursorY int
}

func initMainModel(puz *puzzle.PuzzleDefinition) mainModel {
	grid := initGridModel(puz)
	clues := initCluesModel(puz)

	var initialX int
	var initialY int
	for i, char := range puz.CurrentState {
		if char != '.' {
			initialX = i % puz.NumCols
			initialY = i / puz.NumCols
			break
		}
	}

	return mainModel{
		title: puz.Title,
		author: puz.Author,
		copyright: puz.Copyright,
		acrossClues: puz.AcrossClues,
		downClues: puz.DownClues,
		grid: grid,
		clues: clues,
		solution: puz.Answer,
		cursorX: initialX,
		cursorY: initialY,
	}
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	m.grid, _ = m.grid.Update(msg)
	return m, nil
}

func (m mainModel) View() string {
	header := fmt.Sprintf("%s\n%s %s\n", m.title, m.author, m.copyright)
	if m.validateSolution() {
		header += "Solved!\n"
	}
	footer := ("\nPress ctrl+c to quit.\n")
	return lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		lipgloss.JoinHorizontal(lipgloss.Top, m.grid.View(), m.clues.View()),
		footer,
	)
}

	func (m mainModel) validateSolution() bool {
		grid := m.grid.(gridModel).Grid
		numRows := len(grid)
		numCols := len(grid[0])
		for i := 0; i < numRows; i++ {
			for j := 0; j < numCols; j++ {
				if (grid[i][j] != string(m.solution[(i*numCols)+j])) {
					return false	
				}
			}
		}
		return true 
	}

func Run(puz *puzzle.PuzzleDefinition) {
	p := tea.NewProgram(initMainModel(puz))
	if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
	}
}