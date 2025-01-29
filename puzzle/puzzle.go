package puzzle

import (
	"fmt"
	"sort"
)

type PuzzleDefinition struct {
	Title string
	Author string
	Copyright string
	Version string
	Notes string
	NumRows int
	NumCols int
	AcrossClues map[int]string
	DownClues map[int]string
	NumClues int
	Answer string
	CurrentState string
}

func (puz PuzzleDefinition) Format() {
	fmt.Println(puz.Title, puz.Author, puz.Copyright)
	acrossKeys := make([]int, 0, len(puz.AcrossClues))
	downKeys := make([]int, 0, len(puz.DownClues))
	for key := range puz.AcrossClues {
		acrossKeys = append(acrossKeys, key)
	}
	for key := range puz.DownClues {
		downKeys = append(downKeys, key)
	}
	sort.Ints(acrossKeys)
	sort.Ints(downKeys)
	fmt.Println("*** ACROSS ***")
 	for _, key := range acrossKeys {
		fmt.Printf("%dA. %s\n", key, puz.AcrossClues[key])
	}
	fmt.Println("*** DOWN ***")
	for _, key := range downKeys {
		fmt.Printf("%dD. %s\n", key, puz.DownClues[key])
	}
}