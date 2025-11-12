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
	board        [][]sudoku.Cell
	quiting      bool
	fastMode     bool
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

func (s SudokuGame) Init() tea.Cmd {
	return nil
}

func (s *SudokuGame) updateCellByInput(value string) {
	if len(value) != 3 {
		return
	}

	valueToUpdate := string(value[2])
	s.UpdateCell(valueToUpdate)
	s.input.SetValue("")
}

func (s *SudokuGame) HandleUpdateInputMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlH:
			s.helpPage = !s.helpPage
			s.buildView()
			return s, nil
		case tea.KeyCtrlC:
			return s, tea.Quit
		case tea.KeyUp:
			return s, tea.Quit
		case tea.KeyEnter:
			value := s.input.Value()
			s.updateCellByInput(value)
			s.buildView()
			return s, nil
		}

		switch str := msg.String(); str {
		case "q", "Q":
			s.quiting = true
			s.buildView()
			return s, nil
		case "u", "U":
			if s.highlightRow != -1 && s.highlightCol != -1 {
				s.input.SetValue(fmt.Sprintf("%d%d", s.highlightRow+1, s.highlightCol+1))
			}
			s.buildView()
			return s, nil
		case "m", "M":
			s.inputEnable = !s.inputEnable
			s.buildView()
			return s, nil
		case "f", "F":
			s.fastMode = !s.fastMode
			s.buildView()
			return s, nil
		default:
			_, err := strconv.Atoi(str)
			if err != nil && str != "backspace" {
				s.buildView()
				return s, nil
			}

			cmd := s.updateInput(msg)
			s.buildView()
			return s, cmd
		}
	}

	return s, nil
}

func (s *SudokuGame) updateInput(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	s.input, cmd = s.input.Update(msg)
	s.input.Focus()
	value := s.input.Value()
	l := len(value)
	if l == 1 {
		result := strToInt(string(value[0]))
		if result != -1 {
			s.highlightRow = result - 1
		}
	}

	if l == 2 {
		result := strToInt(string(value[1]))
		if result != -1 {
			s.highlightCol = result - 1
		}
	}

	if s.fastMode && l == 3 {
		s.UpdateCell(string(value[2]))
		s.input.SetValue("")
	}
	return cmd
}

func (s *SudokuGame) HandleUpdateMoveMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlH:
			s.helpPage = !s.helpPage
			s.buildView()
			return s, nil
		case tea.KeyCtrlC:
			return s, tea.Quit
		case tea.KeyDown:
			s.MoveCursor(dirDown)
			s.buildView()
			return s, nil
		case tea.KeyUp:
			s.MoveCursor(dirUp)
			s.buildView()
			return s, nil
		case tea.KeyLeft:
			s.MoveCursor(dirLeft)
			s.buildView()
			return s, nil
		case tea.KeyRight:
			s.MoveCursor(dirRight)
			s.buildView()
			return s, nil
		}

		switch str := msg.String(); str {
		case "q", "Q":
			s.quiting = true
			s.buildView()
			return s, nil
		case "y", "Y":
			s.buildView()
			return s, tea.Quit
		case "j", "s", "J", "S":
			s.MoveCursor(dirDown)
			s.buildView()
			return s, nil
		case "k", "w", "K", "W":
			s.MoveCursor(dirUp)
			s.buildView()
			return s, nil
		case "h", "a", "H", "A":
			s.MoveCursor(dirLeft)
			s.buildView()
			return s, nil
		case "l", "d", "L", "D":
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
		if s.inputEnable {
			return s.HandleUpdateInputMode(msg)
		} else {
			return s.HandleUpdateMoveMode(msg)
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
	if s.highlightRow != -1 && s.highlightCol != -1 {
		coor := sudoku.Coor{
			X: s.highlightRow,
			Y: s.highlightCol,
		}
		s.board = s.s.ValidateNewCell(coor, num)
	}
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

func (s *SudokuGame) printHelpPage() {
	str := "Press Ctrl+H to toggle Help\r\n"
	str += "Press m to change between input mode or movement mode\r\n\r\n"

	str += "------------ Input Mode ------------\r\n\r\n"
	str += "  - Use F to toggle fast mode\r\n"
	str += "  - Use the input to make your move and change the cell\r\n"
	str += "  - First value is the row, second one is column and third is the value the cell will hold\r\n"
	str += "  - Press enter to make the move\r\n"
	str += "  - If fast mode is enabled, you don't need to press enter to make the mov\r\n"
	str += "  - Press U after you made your move to fill the input with the previous column and row\r\n\r\n"

	str += "------------ Movement Mode ------------\r\n\r\n"
	str += "  - Use Up/K/W, Right/L/D, Left/H/A, Down/J/S to move the cursor\r\n"
	str += "  - Press any number between 0-9 to update a cell\r\n"
	s.content = str
}

func (s *SudokuGame) buildView() {
	if s.helpPage {
		s.printHelpPage()
		return
	}
	if s.quiting {
		s.content = "Are you sure you want to quit? Y(es)/N(o)"
		return
	}

	str := ""
	validBoard := s.s.IsValidBoard()
	if validBoard {
		str += "YOU WON!!!!!!"
	}
	str += "Press Ctrl+H for help\r\n\r\n"

	str += "     "
	lastRow := "\r\n     "
	for i := range 9 {
		strToAdd := ""
		if i == s.highlightCol {
			strToAdd += "|" + inStBgYello.Render(fmt.Sprintf(" %d ", i+1))
		} else {
			strToAdd += fmt.Sprintf("| %d ", i+1)
		}
		str += strToAdd
		lastRow += strToAdd
	}

	str += "\r\n\r\n"
	lastRow += "\r\n\r\n"
	str += "     ------------------------------------- \r\n"

	for i := range 9 {
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
		for j := range 9 {
			colorNumCol := (int(math.Floor(float64(j)/3))%3)%2 == 0
			cell := s.board[i][j]

			var color lipgloss.Style
			if cell.HasErr {
				color = inStTextRed
			} else {
				if colorNumRow == colorNumCol {
					color = inStTextBlue
				} else {
					color = inStTextGreen
				}
				if cell.Static {
					color = inputStyle
				}
			}

			if j == 0 {
				strInner += color.Render("|")
			}
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
		str += fmt.Sprintf("%s%s\r\n", strInner, lastCol)
		str += "     ------------------------------------- \r\n"
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
		if s.fastMode {
			s.content += "Fast Mode: ON"
		} else {
			s.content += "Fast Mode: OFF"
		}
	}
}

func (s SudokuGame) View() string {
	return s.content
}
