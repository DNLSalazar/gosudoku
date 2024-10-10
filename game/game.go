package game

import (
	"fmt"
	"math"
	"strconv"
	"sudoku/sudoku"

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
}

func NewSudokuGame(s *sudoku.Sudoku) *tea.Program {
	var input textinput.Model
	input = textinput.New()
	input.Placeholder = "RCV"
	input.Focus()
	input.CharLimit = 3
	input.Width = 5
	input.Prompt = ""

	game := SudokuGame{
		s:            s,
		highlightRow: 0,
		highlightCol: 0,
		helpPage:     false,
		input:        input,
	}

	game.buildView()
	return tea.NewProgram(game)
}

func (s SudokuGame) Init() tea.Cmd {
	return nil
}

func (s SudokuGame) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return s, nil
		case tea.KeyCtrlC:
			return s, tea.Quit
		}

		switch str := msg.String(); str {
		case "q", "Q":
			s.buildView()
			return s, tea.Quit
		case "j":
			s.MoveCursor(dirDown)
			s.buildView()
			return s, nil
		case "k":
			s.MoveCursor(dirUp)
			s.buildView()
			return s, nil
		case "h":
			s.MoveCursor(dirLeft)
			s.buildView()
			return s, nil
		case "l":
			s.MoveCursor(dirRight)
			s.buildView()
			return s, nil
		case "m":
			s.inputEnable = !s.inputEnable
			s.buildView()
			return s, nil
		default:
			s.UpdateCell(str)
			s.buildView()
			return s, nil
		}
	}
	return s, nil
}

func strToInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		return -1
	}
	return num
}

func (s *SudokuGame) UpdateCell(str string) {
	num := strToInt(str)
	if num == -1 {
		return
	}
	coor := sudoku.Coor{
		X: s.highlightRow,
		Y: s.highlightCol,
	}
	s.s.ValidateNewCell(coor, num)
}

func (s *SudokuGame) MoveCursor(dir CursorDirection) {
	switch dir {
	case dirUp:
		addToCursor(&(s.highlightRow), -1)
	case dirDown:
		addToCursor(&(s.highlightRow), +1)
	case dirLeft:
		addToCursor(&(s.highlightCol), -1)
	case dirRight:
		addToCursor(&(s.highlightCol), +1)
	}
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

func (s *SudokuGame) buildView() {
	validBoard := s.s.IsValidBoard()
	if validBoard {
		s.content = "YOU WON!!!!!!"
		return
	}
	str := "     "
	lastRow := "\r\n\r\n     "
	for i := 0; i < 9; i++ {
		strToAdd := ""
		if i == s.highlightCol {
			strToAdd += "|" + inStBgYello.Render(fmt.Sprintf(" %d ", i+1))
		} else {
			strToAdd += fmt.Sprintf("| %d ", i+1)
		}
		str += strToAdd
		lastRow += strToAdd
	}

	str += fmt.Sprintf("\r\n\r\n")
	lastRow += fmt.Sprintf("\r\n\r\n")
	str += fmt.Sprintf("     ------------------------------------- \r\n")

	for i := 0; i < 9; i++ {
		strInner := ""
		lastCol := ""
		strToAdd := ""
		if i == s.highlightRow {
			strToAdd += inStBgYello.Render(fmt.Sprintf(" %d ", i+1))
			strInner += " " + strToAdd + " "
			lastCol += " " + strToAdd
		} else {
			strToAdd += fmt.Sprintf("  %d  ", i+1)
			strInner += strToAdd
			lastCol += strToAdd
		}

		colorNumRow := (int(math.Floor(float64(i)/3))%3)%2 == 0
		for j := 0; j < 9; j++ {
			colorNumCol := (int(math.Floor(float64(j)/3))%3)%2 == 0

			var color lipgloss.Style
			if colorNumRow == colorNumCol {
				color = inStTextBlue
			} else {
				color = inStTextGreen
			}

			if j == 0 {
				strInner += color.Render("|")
			}
			cell := s.s.GetCell(sudoku.Coor{X: i, Y: j})
			var value string
			if cell.Value == 0 {
				value = "   "
			} else {
				value = " " + fmt.Sprint(cell.Value) + " "
			}
			if i == s.highlightRow && j == s.highlightCol {
				value = inStBgYello.Render(value)
			}
			strInner += color.Render(fmt.Sprintf("%s|", value))
		}
		str += fmt.Sprintf("%s   %s\r\n", strInner, lastCol)
		str += fmt.Sprintf("     ------------------------------------- \r\n")
	}
	str += lastRow
	s.content = str
	s.printError()
	s.addInput()
}

func (s SudokuGame) boardView() string {
	return s.content
}

func (s *SudokuGame) printError() {
	if s.s.HasErr {
		s.content += inStTextRed.Render("Error on board")
	}
}

func (s *SudokuGame) addInput() {
	if s.inputEnable {
		s.content += inputStyle.Width(5).Render("Move") + s.input.View()
	}
}

func (s SudokuGame) View() string {
	return s.content
}
