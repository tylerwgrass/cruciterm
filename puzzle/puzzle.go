package puzzle

import "fmt"

type PuzzleDefinition struct {
	Title string
	Author string
	Copyright string
	Version string
	Notes string
	NumRows int
	NumCols int
	Clues map[int]Clue
	NumClues int
	Answer string
	CurrentState string
}

type Clue struct {
	Num int
	StartY int
	StartX int
	AcrossClue string
	DownClue string
}

var Clues map[int]Clue
var AcrossClues []*Clue
var DownClues []*Clue

func (puz * PuzzleDefinition) AssignClues(clues []string) {
	Clues = make(map[int]Clue)
	clueNum := 1
	clueIndex := 0
	for i := 0; i < len(puz.Answer); i++ {
		if string(puz.Answer[i]) == "." {
			continue
		}

		row := i / puz.NumCols
		col := i % puz.NumCols

		isAcrossClueStart := col == 0 || string(puz.Answer[(puz.NumCols * (row)) + col - 1]) == "." 
		isDownClueStart := row == 0 || string(puz.Answer[(puz.NumCols * (row - 1)) + col]) == "." 
		if (!(isAcrossClueStart || isDownClueStart)) { 
			continue
		}
		
		clue := Clue{
			Num: clueNum,
			StartY: row,
			StartX: col,
		}

		if isAcrossClueStart {
			AcrossClues = append(AcrossClues, &clue)
			clue.AcrossClue = clues[clueIndex]
			clueIndex++
		}

		if isDownClueStart {
			DownClues = append(DownClues, &clue)
			clue.DownClue = clues[clueIndex]
			clueIndex++
		}
		Clues[clueNum] = clue

		clueNum++
	}
	puz.Clues = Clues
}

func (p PuzzleDefinition) ToString() string {
	output := "Clues:\n"
	for i := range(p.NumClues) {
		output += fmt.Sprintf("%s\n", p.Clues[i].ToString())
	}

	return output
}

func (c Clue) ToString() string {
	return fmt.Sprintf("Clue{num: %d, r: %d, c:%d, acrossClue: %s, downClue: %s}", 
		c.Num,
		c.StartY,
		c.StartX,
		c.AcrossClue,
		c.DownClue,	
	)
}