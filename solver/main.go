package solver

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tylerwgrass/cruciterm/logger"
	"github.com/tylerwgrass/cruciterm/puzzle"
	"golang.org/x/term"
)

type mainModel struct {
	title string
	author string
	copyright string
	clues tea.Model
	grid tea.Model
}

func initMainModel(puz *puzzle.PuzzleDefinition) mainModel {
	grid := initGridModel(puz)
	clues := initCluesModel(puz)

	return mainModel{
		title: puz.Title,
		author: puz.Author,
		copyright: puz.Copyright,
		grid: grid,
		clues: clues,
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
	m.clues, _ = m.clues.Update(msg)
	return m, nil
}

func (m mainModel) View() string {
	width, height, err := term.GetSize(0)
	if err != nil {
		logger.Debug(fmt.Sprintf("Failed to get terminal size"))
		width = 500
		height = 500
	}
	header := fmt.Sprintf("%s\n%s %s\n", m.title, m.author, m.copyright)
	if m.grid.(gridModel).solved {
		header += "Solved!\n"
	}
	footer := ("\nPress ctrl+c to quit.\n")
	mainContent := lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		lipgloss.NewStyle().AlignVertical(lipgloss.Center).Render(lipgloss.JoinHorizontal(lipgloss.Top, m.grid.View(), m.clues.View())),
		footer,
	)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, mainContent)
}

func Run(puz *puzzle.PuzzleDefinition) {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	os.Truncate("debug.log", 0)
	logger.SetLogFile(f)
	defer f.Close()

	p := tea.NewProgram(initMainModel(puz), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
	}
}