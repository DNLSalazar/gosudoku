package game

import (
	"fmt"
	"strconv"
	"sudoku/sudoku"

	tea "github.com/charmbracelet/bubbletea"
)

func (s *SudokuGame) HandleUpdateGameEnded(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlH:
			s.helpPage = !s.helpPage
			s.buildView()
			return s, nil
		case tea.KeyCtrlC:
			return s, tea.Quit
		}

		switch str := msg.String(); str {
		case "q", "Q":
			s.quiting = true
			s.buildView()
			return s, nil
		case "m", "M":
			s.inputEnable = !s.inputEnable
			s.buildView()
			return s, nil
		}
	}
	return s, nil
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
			if s.s.IsValidBoard() {
				return s, nil
			}
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

func (s *SudokuGame) updateCellByInput(value string) {
	if len(value) != 3 {
		return
	}

	valueToUpdate := string(value[2])
	s.UpdateCell(valueToUpdate)
	s.input.SetValue("")
}
