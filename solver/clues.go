package solver

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/lipgloss/v2/list"
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
	activeClueOrientation Orientation
}

func initCluesModel(puz *puzzle.PuzzleDefinition) cluesModel {
  acrossClues, downClues := organizeClues(puz)
	return cluesModel{
		acrossClues: acrossClues,
		downClues: downClues,
		activeClueOrientation: Horizontal,
	}
}

func (m cluesModel) Init() tea.Cmd {
	return nil
}

func (m cluesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (m cluesModel) View() string {
	CONTAINER_WIDTH := 80 
	COLUMN_WIDTH := 36
	clueContainerStyle := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Width(CONTAINER_WIDTH).
			Padding(0, 2)
	return clueContainerStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(COLUMN_WIDTH).Border(lipgloss.HiddenBorder()).Render(
			lipgloss.JoinVertical( lipgloss.Left,
				lipgloss.PlaceHorizontal(COLUMN_WIDTH, lipgloss.Center, "~~~ ACROSS ~~~"),
				m.acrossClues.String(),
			)), 
		lipgloss.NewStyle().Width(COLUMN_WIDTH).Border(lipgloss.HiddenBorder()).Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.PlaceHorizontal(COLUMN_WIDTH, lipgloss.Center, "~~~ DOWN ~~~"),
			m.downClues.String(),
		)),
	))
}

func organizeClues(puz *puzzle.PuzzleDefinition) (*list.List, *list.List) {
	activeClueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#EE6FF8"))
	crossClueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#EFC1F3"))
	acrossClues := list.New().
		Enumerator(func(_ list.Items, i int) string {
			return fmt.Sprintf("%d. ",puz.AcrossClues[i].Num)
		}).
		ItemStyleFunc(func(items list.Items, i int) lipgloss.Style {
			if currentAcrossClue.Num == puz.AcrossClues[i].Num {
				if solvingOrientation == Horizontal {
					return activeClueStyle
				} else {
					return crossClueStyle
				}
			}
				return lipgloss.NewStyle()
		})
	downClues := list.New().
		Enumerator(func(_ list.Items, i int) string {
			return fmt.Sprintf("%d. ",puz.DownClues[i].Num)
		}).
		ItemStyleFunc(func(_ list.Items, i int) lipgloss.Style {
			if currentDownClue.Num == puz.DownClues[i].Num {
				if solvingOrientation == Vertical {
					return activeClueStyle
				} else {
					return crossClueStyle
				}
			}
				return lipgloss.NewStyle()
		})
	
 	for _, clue := range puz.AcrossClues {
		acrossClues.Item(clue.Clue)
	}

	for _, clue := range puz.DownClues {
		downClues.Item(clue.Clue)
	}

	return acrossClues, downClues
}