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

var NUM_SHOWN_CLUES int = 9
var acrossClues []*puzzle.Clue
var downClues []*puzzle.Clue
var currentAcrossClue *puzzle.Clue
var currentDownClue *puzzle.Clue

type cluesModel struct {
	acrossClues []*puzzle.Clue
	downClues []*puzzle.Clue
	activeClueOrientation Orientation
}

func initCluesModel(puz *puzzle.PuzzleDefinition) cluesModel {
	acrossClues, downClues = puz.AcrossClues, puz.DownClues
	return cluesModel{
		acrossClues: puz.AcrossClues,
		downClues: puz.DownClues,
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
	renderedAcrossClues := getClueRendering(currentAcrossClue, acrossClues, Horizontal)
	renderedDownClues := getClueRendering(currentDownClue, downClues, Vertical)
	return clueContainerStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(COLUMN_WIDTH).Border(lipgloss.HiddenBorder()).Render(
			lipgloss.JoinVertical( lipgloss.Left,
				lipgloss.PlaceHorizontal(COLUMN_WIDTH, lipgloss.Center, "~~~ ACROSS ~~~"),
				renderedAcrossClues,
			)), 
		lipgloss.NewStyle().Width(COLUMN_WIDTH).Border(lipgloss.HiddenBorder()).Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.PlaceHorizontal(COLUMN_WIDTH, lipgloss.Center, "~~~ DOWN ~~~"),
				renderedDownClues,
		)),
	))
}

func getClueRendering(currentClue *puzzle.Clue, clues []*puzzle.Clue, orientation Orientation) string {
	var currentClueIndex int
	for i, clue := range(clues) {
		if clue == currentClue {
			currentClueIndex = i
			break
		}
	}
	rangeStart, rangeEnd := currentClueIndex, currentClueIndex
	for rangeEnd - rangeStart + 1 < NUM_SHOWN_CLUES && rangeStart - rangeEnd + 1 != len(clues) {
		if rangeStart != 0 {
			rangeStart--
		}

		if rangeEnd != len(clues) - 1 {
			rangeEnd++
		}
	}

	activeClueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#EE6FF8"))
	crossClueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#EFC1F3"))
	clueList := list.New().
		Enumerator(func(_ list.Items, i int) string {
			return fmt.Sprintf("%d. ", clues[i + rangeStart].Num)
		}).
		ItemStyleFunc(func(_ list.Items, i int) lipgloss.Style {
			if currentClue.Num == clues[i + rangeStart].Num {
				if solvingOrientation == orientation {
					return activeClueStyle
				} else {
					return crossClueStyle
				}
			}
				return lipgloss.NewStyle()
		})

	for i := rangeStart; i <= rangeEnd; i++ {
		clueList.Item(clues[i].Clue)	
	}

	rendered := clueList.String()
	if rangeStart != 0 {
		rendered  = lipgloss.JoinVertical(lipgloss.Left, "...", rendered)
	}
	
	if rangeEnd != len(clues) - 1 {
		rendered = lipgloss.JoinVertical(lipgloss.Left, rendered, "...")
	}

	return rendered
}
