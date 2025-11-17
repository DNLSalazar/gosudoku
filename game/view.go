package game

import (
	"fmt"
	"math"

	"github.com/charmbracelet/lipgloss"
)

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
	str += "Press Ctrl+H for help\r\n\r\n"
	if validBoard {
		str += "YOU WON!!!!!!\r\n\r\n"
	}

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
