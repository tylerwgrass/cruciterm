package loader

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/tylerwgrass/cruciterm/puzzle"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

var ErrFileParse = fmt.Errorf("could not parse file");
var ErrFileNotSupported = fmt.Errorf("file type not supported")

func LoadFile(path string) (puzzle.PuzzleDefinition, error) {
	var ext = filepath.Ext(path)
	if ext == ".puz" {
		return loadPuzFile(path)
	}
	return puzzle.PuzzleDefinition{}, ErrFileNotSupported
}

// .puz file definition: https://code.google.com/archive/p/puz/wikis/FileFormat.wiki
func loadPuzFile(path string) (puzzle.PuzzleDefinition, error) {
	file, err := os.Open(path) 
	if err != nil {
		return puzzle.PuzzleDefinition{}, err
	}
	defer file.Close()

	puz := puzzle.PuzzleDefinition{}
	if err := parseHeader(&puz, file); err != nil {
		return puzzle.PuzzleDefinition{}, ErrFileParse
	}
	if err := parseState(&puz, file); err != nil {
		return puzzle.PuzzleDefinition{}, ErrFileParse
	}
	if err = parseContent(&puz, file); err != nil {
		return puzzle.PuzzleDefinition{}, ErrFileParse
	}
	return puz, nil
}

func parseHeader(puz *puzzle.PuzzleDefinition, file *os.File) error {
	if _, err := file.Seek(0x18, io.SeekStart); err != nil {
		return err
	}

	ver := make([]byte, 0x4)
	if _, err := file.Read(ver); err != nil {
		return err
	}
	puz.Version = string(ver)

	if _, err := file.Seek(0x2C, io.SeekStart); err != nil {
		return err
	}

	dimensions := make([]byte, 0x4)
	if _, err := file.Read(dimensions); err != nil {
		return err
	}

	puz.NumCols = int(dimensions[0])
	puz.NumRows = int(dimensions[1])
	puz.NumClues = int(binary.LittleEndian.Uint16(dimensions[2:]))
	return nil
}

func parseState(puz *puzzle.PuzzleDefinition, file *os.File) error {
	if _, err := file.Seek(0x34, io.SeekStart); err != nil {
		return err
	}
	numCells := puz.NumCols * puz.NumRows
	puzzleState := make([]byte, numCells * 2)
	if _, err := file.Read(puzzleState); err != nil {
		return err
	}
	puz.Answer = string(puzzleState[:numCells])
	puz.CurrentState = string(puzzleState[numCells:])
	return nil
}

func parseContent(puz *puzzle.PuzzleDefinition, file *os.File) error {
	decoder := transform.NewReader(file, charmap.ISO8859_1.NewDecoder())
	reader := bufio.NewReader(decoder)
	delim := byte(0)
	content := make([]string, puz.NumClues + 4)
	index := 0
	for index < len(content) {
		contentBytes, err := reader.ReadBytes(delim)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}
		content[index] = string(contentBytes[:len(contentBytes) - 1])
		index++
	}

	puz.Title = content[0]
	puz.Author = content[1]
	puz.Copyright = content[2]
	puz.Notes = content[len(content)-1]
	assignClues(puz, content[3:len(content)-1])
	return nil
}

func assignClues(puz *puzzle.PuzzleDefinition, clues []string) {
	clueNum := 1
	clueIndex := 0
	acrossClues := make(map[int]string)
	downClues := make(map[int]string)
	for i := 0; i < len(puz.Answer); i++ {
		if string(puz.Answer[i]) == "." {
			continue
		}

		isStartOfClue := false
		row := i / puz.NumCols
		col := i % puz.NumCols

		if col == 0 || string(puz.Answer[(puz.NumCols * (row)) + col - 1]) == "." {
			acrossClues[clueNum] = clues[clueIndex]
			clueIndex++
			isStartOfClue = true
		}

		if row == 0 || string(puz.Answer[(puz.NumCols * (row - 1)) + col]) == "." {
			downClues[clueNum] = clues[clueIndex]
			clueIndex++
			isStartOfClue = true
		}

		if isStartOfClue {
			clueNum++
		}
	}
	
	puz.AcrossClues = acrossClues
	puz.DownClues = downClues
}