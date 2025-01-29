package main

import (
	"fmt"

	"github.com/tylerwgrass/cruciterm/loader"
	"github.com/tylerwgrass/cruciterm/puzzle"
)

func main() {
	puzFile := "./puzzles/test.puz"
	puz, err := loader.LoadFile(puzFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	puzzle.Format(&puz)
}