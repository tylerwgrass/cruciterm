package solver

import (
	"fmt"
	"os"

	"github.com/tylerwgrass/cruciterm/puzzle"

	tea "github.com/charmbracelet/bubbletea"
)


func Run(puz *puzzle.PuzzleDefinition) {
	p := tea.NewProgram(initialModel(puz))
	if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
	}
}