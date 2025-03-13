package solver

import (
	"fmt"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/tylerwgrass/cruciterm/puzzle"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var acrossKeys []int
var downKeys []int
var currentAcrossClue int
var currentDownClue int

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
	acrossKeys = make([]int, 0)
	downKeys = make([]int, 0)

	for key, clue := range puz.Clues {
		if clue.AcrossClue != "" {
			acrossKeys = append(acrossKeys, key)
		}
		if clue.DownClue != "" {
			downKeys = append(downKeys, key)
		}
	}

	sort.Ints(acrossKeys)
	sort.Ints(downKeys)

	acrossClues := list.New()
	downClues := list.New()

 	for _, key := range acrossKeys {
		acrossClues.Item(fmt.Sprintf("%d. %s", key, puz.Clues[key].AcrossClue)).
			Enumerator(clueEnumerator).
			ItemStyleFunc(func(_ list.Items, i int) lipgloss.Style {
					if acrossKeys[i] == currentAcrossClue {
						return lipgloss.NewStyle().Foreground(lipgloss.Color("#EE6FF8"))
					}
					return lipgloss.NewStyle()
			})
	}
	for _, key := range downKeys {
		downClues.Item(fmt.Sprintf("%d. %s", key, puz.Clues[key].DownClue)).
			Enumerator(clueEnumerator).
			ItemStyleFunc(func(_ list.Items, i int) lipgloss.Style {
					if downKeys[i] == currentDownClue {
						return lipgloss.NewStyle().Foreground(lipgloss.Color("#EE6FF8"))
					}
					return lipgloss.NewStyle()
			})
	}
	return acrossClues, downClues
}