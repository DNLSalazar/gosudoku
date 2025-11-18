package game

import (
	"strconv"

	"github.com/DNLSalazar/gosudoku/sudoku"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func strToInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		return -1
	}
	return num
}

func addToCursor(cursorAxis *int, value int) {
	if *cursorAxis+value < 0 {
		*cursorAxis = 8
		return
	}
	if *cursorAxis+value > 8 {
		*cursorAxis = 0
		return
	}
	*cursorAxis += value
}

func NewSudokuGame(s *sudoku.Sudoku) *tea.Program {
	var input textinput.Model
	input = textinput.New()
	input.Placeholder = "RCV"
	input.Focus()
	input.CharLimit = 3
	input.Width = 5
	input.Prompt = ""
	board := s.GetBoard()

	game := SudokuGame{
		s:            s,
		highlightRow: 0,
		highlightCol: 0,
		helpPage:     false,
		input:        input,
		board:        board,
	}

	game.buildView()
	return tea.NewProgram(game)
}
