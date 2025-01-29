package main

import (
	"fmt"
	"os"

	"github.com/tylerwgrass/cruciterm/loader"
	"github.com/tylerwgrass/cruciterm/puzzle"
)

var TEST_FILE_PATH string = "./puzzles/test.puz"

func main() {
	puzFilePath := TEST_FILE_PATH
	if len(os.Args) == 2 {
		puzFilePath = os.Args[1]
	}

	puz, err := loader.LoadFile(puzFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	puzzle.Format(&puz)
}