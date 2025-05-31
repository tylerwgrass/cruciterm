package puzzle

import (
	"fmt"
	"strings"
)

type PuzzleDefinition struct {
	Title        string
	Author       string
	Copyright    string
	Version      string
	Notes        string
	NumRows      int
	NumCols      int
	NumClues     int
	AcrossClues  []*Clue
	DownClues    []*Clue
	Answer       string
	CurrentState string
}

type Clue struct {
	Num      int
	StartRow int
	StartCol int
	EndRow   int
	EndCol   int
	Clue     string
	Answer   string
}

var Clues map[int]Clue
var AcrossClues []*Clue
var DownClues []*Clue

func (puz *PuzzleDefinition) AssignClues(clues []string) {
	Clues = make(map[int]Clue)
	clueNum := 1
	clueIndex := 0
	for i := 0; i < len(puz.Answer); i++ {
		if string(puz.Answer[i]) == "." {
			continue
		}

		row := i / puz.NumCols
		col := i % puz.NumCols

		isAcrossClueStart := col == 0 || string(puz.Answer[(puz.NumCols*(row))+col-1]) == "."
		isDownClueStart := row == 0 || string(puz.Answer[(puz.NumCols*(row-1))+col]) == "."

		if !(isAcrossClueStart || isDownClueStart) {
			continue
		}

		if isAcrossClueStart {
			clue := puz.parseClue(clueNum, row, col, true)
			clue.Clue = clues[clueIndex]
			AcrossClues = append(AcrossClues, clue)
			clueIndex++
		}

		if isDownClueStart {
			clue := puz.parseClue(clueNum, row, col, false)
			clue.Clue = clues[clueIndex]
			DownClues = append(DownClues, clue)
			clueIndex++
		}
		clueNum++
	}
	puz.AcrossClues = AcrossClues
	puz.DownClues = DownClues
}

func (p PuzzleDefinition) parseClue(clueNumber, startRow, startCol int, isAcrossClue bool) *Clue {
	clue := Clue{
		Num:      clueNumber,
		StartRow: startRow,
		StartCol: startCol,
	}
	row, col := startRow, startCol
	var sb strings.Builder
	var dr, dc int
	if isAcrossClue {
		dr, dc = 0, 1
	} else {
		dr, dc = 1, 0
	}
	for ok := true; ok; ok = p.isVisitable(row, col) {
		cellLocation := (row * p.NumCols) + col
		sb.WriteString(string(p.Answer[cellLocation]))
		row, col = row+dr, col+dc
	}

	if isAcrossClue {
		clue.EndRow = startRow
		clue.EndCol = col - 1
	} else {
		clue.EndRow = row - 1
		clue.EndCol = startCol
	}
	clue.Answer = sb.String()
	return &clue
}

func (p PuzzleDefinition) isVisitable(row, col int) bool {
	return row < p.NumRows &&
		col < p.NumCols &&
		row > -1 &&
		col > -1 &&
		string(p.Answer[(row*p.NumCols)+col]) != "."
}

func (p PuzzleDefinition) String() string {
	output := "~~~Across Clues~~~\n"
	for _, clue := range p.AcrossClues {
		output += fmt.Sprintf("%v\n", clue)
	}
	output += "~~~Down Clues~~~\n"
	for _, clue := range p.DownClues {
		output += fmt.Sprintf("%v\n", clue)
	}
	return output
}

func (c Clue) String() string {
	return fmt.Sprintf("Clue{num: %d, r: %d, c:%d, endr: %d, endc: %d, Clue: %s, Answer: %s}",
		c.Num,
		c.StartRow,
		c.StartCol,
		c.EndRow,
		c.EndCol,
		c.Clue,
		c.Answer,
	)
}
