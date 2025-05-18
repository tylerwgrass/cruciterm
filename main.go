package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/tylerwgrass/cruciterm/loader"
	"github.com/tylerwgrass/cruciterm/logger"
	"github.com/tylerwgrass/cruciterm/preferences"
	"github.com/tylerwgrass/cruciterm/solver"
)

var TEST_FILE_PATH string = "./puzzles/test.puz"

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	os.Truncate("debug.log", 0)
	logger.SetLogFile(f)
	defer f.Close()
	puzFilePath := TEST_FILE_PATH
	if len(os.Args) == 2 {
		puzFilePath = os.Args[1]
	}

	puz, err := loader.LoadFile(puzFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	preferences.Init()
	solver.Run(&puz)
}
