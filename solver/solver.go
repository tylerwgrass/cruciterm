package solver

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/v2/help"
	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/stopwatch"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/tylerwgrass/cruciterm/logger"
	"github.com/tylerwgrass/cruciterm/puzzle"
	"golang.org/x/term"
)

type mainModel struct {
	title string
	author string
	copyright string
	clues cluesModel
	grid gridModel
	preferences preferencesModel
	stopwatch stopwatch.Model
	help help.Model
	activeView ActiveView
}

type ActiveView int
const (
	GridAndClues ActiveView = iota
	Preferences
)

var solvingOrientation Orientation = Horizontal

func initMainModel(puz *puzzle.PuzzleDefinition) mainModel {
	grid := initGridModel(puz)
	clues := initCluesModel(puz)
	preferences := initPreferencesModel()
	stopwatch := stopwatch.New()
	help := help.New()
	help.ShowAll = true
	return mainModel{
		stopwatch: stopwatch,
		title: puz.Title,
		author: puz.Author,
		copyright: puz.Copyright,
		grid: grid,
		clues: clues,
		help: help,
		activeView: GridAndClues,
		preferences: preferences,
	}
}

func (m mainModel) Init() tea.Cmd {
	return m.stopwatch.Init()
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.ViewPreferences):
			if m.activeView == Preferences {
				m.activeView = GridAndClues
			} else {
				m.activeView = Preferences
			}
		}
	}

	if m.activeView == Preferences {
		preferences, _ := m.preferences.Update(msg)
		m.preferences = preferences.(preferencesModel)
		return m, nil 
	}

	grid, _ := m.grid.Update(msg)
	m.grid = grid.(gridModel)
	solvingOrientation = m.grid.navOrientation
	var cmd tea.Cmd
	if m.grid.solved {
		cmd = m.stopwatch.Stop()
	} else {
		m.stopwatch, cmd = m.stopwatch.Update(msg)
	}
	return m, cmd
}

func (m mainModel) View() string {
	if m.activeView == Preferences {
		return m.preferences.View()
	}
	
	return m.getSolverView()
}

func (m mainModel) getSolverView() string {
	width, height, err := term.GetSize(0)
	if err != nil {
		logger.Debug("Failed to get terminal size")
		width = 500
		height = 500
	}
	header := lipgloss.NewStyle().PaddingTop(height / 20).Render(fmt.Sprintf("%s\n%s %s\n", m.title, m.author, m.copyright))
	if m.grid.solved {
		header += "Solved!\n"
	}
	footer := m.help.View(keys) 
	mainContent := lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		lipgloss.NewStyle().AlignVertical(lipgloss.Center).Render(
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				lipgloss.JoinVertical(lipgloss.Left, m.grid.View(), m.stopwatch.View()),
				m.clues.View(),
			)),
		footer,
	)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Top, mainContent)
}

func Run(puz *puzzle.PuzzleDefinition) {
	p := tea.NewProgram(initMainModel(puz), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
	}
}