package solver

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/tylerwgrass/cruciterm/puzzle"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var currentAcrossClue *puzzle.Clue
var currentDownClue *puzzle.Clue

type cluesModel struct {
	acrossClues *list.List
	downClues *list.List
}

func clueEnumerator(items list.Items, i int) string {
	return ""
}

func initCluesModel(puz *puzzle.PuzzleDefinition) cluesModel {
  acrossClues, downClues := organizeClues(puz)
	return cluesModel{
		acrossClues: acrossClues,
		downClues: downClues,
	}
}

func (m cluesModel) Init() tea.Cmd {
	return nil
}

func (m cluesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m cluesModel) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top, m.acrossClues.String(), m.downClues.String())
}

func organizeClues(puz *puzzle.PuzzleDefinition) (*list.List, *list.List) {
	acrossClues := list.New()
	downClues := list.New()

 	for _, clue := range puz.AcrossClues {
		acrossClues.Item(fmt.Sprintf("%d. %s", clue.Num, clue.Clue)).
			Enumerator(clueEnumerator).
			ItemStyleFunc(func(items list.Items, i int) lipgloss.Style {
				if currentAcrossClue.Num == puz.AcrossClues[i].Num {
					return lipgloss.NewStyle().Foreground(lipgloss.Color("#EE6FF8"))
				}
					return lipgloss.NewStyle()
			})
	}
	for _, clue := range puz.DownClues {
		downClues.Item(fmt.Sprintf("%d. %s", clue.Num, clue.Clue)).
			Enumerator(clueEnumerator).
			ItemStyleFunc(func(_ list.Items, i int) lipgloss.Style {
					if currentDownClue.Num == puz.DownClues[i].Num {
						return lipgloss.NewStyle().Foreground(lipgloss.Color("#EE6FF8"))
					}
					return lipgloss.NewStyle()
			})
	}
	return acrossClues, downClues
}