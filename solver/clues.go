package solver

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tylerwgrass/cruciterm/puzzle"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type cluesModel struct {
	table table.Model
}

func initCluesModel(puz *puzzle.PuzzleDefinition) cluesModel {
  acrossClues, downClues := organizeClues(puz.AcrossClues, puz.DownClues)

	cols := []table.Column{
		{Title: "Across", Width: 40},
		{Title: "Down", Width: 40},
	}

	rows := buildRows(acrossClues, downClues)

	t := table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithHeight(10),
	)

	style := table.DefaultStyles()
	style.Header = style.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true)

	t.SetStyles(style)

	return cluesModel{
		table: t,
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
		case "down":
			m.table.MoveDown(1)
		case "up":
			m.table.MoveUp(1)
		}
	}
	return m, nil
}

func (m cluesModel) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func buildRows(acrosses []string, downs []string) []table.Row {
	numRows := max(len(acrosses), len(downs))
	rows := make([]table.Row, numRows)
	for i := range numRows {
		row := make([]string, 2)
		if i < len(acrosses) {
			row[0] = acrosses[i]
		} else {
			row[0] = ""
		}
		if i < len(downs) {
			row[1] = downs[i]
		} else {
			row[1] = ""
		}
		rows[i] = row
	}
	return rows
}

func organizeClues(acrosses, downs map[int]string) ([]string, []string) {
	acrossKeys := make([]int, 0, len(acrosses))
	downKeys := make([]int, 0, len(downs))

	for key := range acrosses {
		acrossKeys = append(acrossKeys, key)
	}

	for key := range downs {
		downKeys = append(downKeys, key)
	}

	sort.Ints(acrossKeys)
	sort.Ints(downKeys)

	acrossClues := make([]string, len(acrossKeys))
	downClues := make([]string, len(downKeys))

 	for i, key := range acrossKeys {
		acrossClues[i] = fmt.Sprintf("%d. %s", key, acrosses[key])
	}
	for i, key := range downKeys {
		downClues[i] = fmt.Sprintf("%d. %s", key, downs[key])
	}

	return acrossClues, downClues
}