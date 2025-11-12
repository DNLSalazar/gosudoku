package menu

import (
	"fmt"
	"sudoku/server"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	ServeGame = iota
	PlayOnTerminal
)

var items = []string{
	"Serve game and play on any browser.",
	"Play on Terminal, as fun as it sounds.",
}

type App struct {
	selectedIndex *int
	helpPage      bool
}

func (a App) Init() tea.Cmd {
	return nil
}

func NewMenuApp(selection *int) *tea.Program {
	app := App{
		selectedIndex: selection,
	}
	return tea.NewProgram(app)
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if a.helpPage {
		return a.handleHelpPage(msg)
	}
	return a.handleSelectPage(msg)
}

func (a App) handleHelpPage(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlH:
			return a.toggleHelp()
		}
		switch msg.String() {
		case "q", "Q":
			return a.toggleHelp()
		}
	}
	return a, nil
}

func (a App) handleSelectPage(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlH:
			return a.toggleHelp()
		case tea.KeyDown:
			a.updateSelected(1)
		case tea.KeyUp:
			a.updateSelected(-1)
		case tea.KeyEnter:
			tea.Quit()
			return a, tea.Quit
		case tea.KeyCtrlC:
			tea.Quit()
			return a, tea.Quit
		}
		switch msg.String() {
		case "q", "Q":
			tea.Quit()
			return a, tea.Quit
		case "j":
			a.updateSelected(1)
		case "k":
			a.updateSelected(-1)
		}
	}
	return a, nil
}

func (a App) updateSelected(n int) {
	newValue := *a.selectedIndex + n
	if newValue < 0 {
		*a.selectedIndex = len(items) - 1
		return
	}
	if newValue >= len(items) {
		*a.selectedIndex = 0
		return
	}
	*a.selectedIndex = newValue
}

func (a App) toggleHelp() (tea.Model, tea.Cmd) {
	a.helpPage = !a.helpPage
	return a, nil
}

func (a App) View() string {
	if a.helpPage {
		return a.buildHelpPage()
	}
	return a.buildMainPage()
}

func (a App) buildMainPage() string {
	var s string
	s += "Press Ctrl+H to toggle Help\r\n"
	s += "Use Up/k and Down/j to move across the options, press enter to confirm you selection\r\n\r\n"
	s += "Please select an option\r\n"

	for i, v := range items {
		s += "["
		if i == *a.selectedIndex {
			s += "x"
		} else {
			s += " "
		}
		s += "]"

		s += fmt.Sprintf("  %s\r\n", v)
	}

	return s
}

func (a App) buildHelpPage() string {
	str := "Press Ctrl+H to toggle Help\r\n"

	str += "------------ Serve Game ------------\r\n\r\n"
	str += fmt.Sprintf("Serve a sudoku game on port %s available to play in any browser\r\n", server.PORT)
	str += "with multiplayer support\r\n\r\n"

	str += "------------ Terminal Game ------------\r\n\r\n"
	str += "Play a game on the terminal by yourself, no multiplayer"
	return str
}
