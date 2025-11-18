package game

import (
	"github.com/DNLSalazar/gosudoku/sudoku"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const dartGray = lipgloss.Color("#767676")
const blue = lipgloss.Color("#35a8ed")
const green = lipgloss.Color("#46c759")
const red = lipgloss.Color("#e61723")
const yellow = lipgloss.Color("#f8c30b")
const black = lipgloss.Color("#000000")

type CursorDirection int

const (
	dirUp CursorDirection = iota + 1
	dirDown
	dirLeft
	dirRight
)

var inputStyle = lipgloss.NewStyle().Foreground(dartGray)
var inStTextBlue = lipgloss.NewStyle().Foreground(blue)
var inStTextGreen = lipgloss.NewStyle().Foreground(green)
var inStTextRed = lipgloss.NewStyle().Foreground(red)
var inStBgYello = lipgloss.NewStyle().Background(yellow).Foreground(black)

type SudokuGame struct {
	s            *sudoku.Sudoku
	input        textinput.Model
	highlightRow int
	highlightCol int
	helpPage     bool
	inputEnable  bool
	content      string
	board        [][]sudoku.Cell
	quiting      bool
	fastMode     bool
}

func (s SudokuGame) Init() tea.Cmd {
	return nil
}

func (s SudokuGame) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if s.quiting {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "y", "Y":
				return s, tea.Quit
			case "n", "N":
				s.quiting = false
				s.buildView()
				return s, nil
			}
		}
	} else {
		if s.s.IsValidBoard() {
			return s.HandleUpdateGameEnded(msg)
		} else {
			if s.inputEnable {
				return s.HandleUpdateInputMode(msg)
			} else {
				return s.HandleUpdateMoveMode(msg)
			}
		}
	}
	return s, nil
}

func (s SudokuGame) View() string {
	return s.content
}
