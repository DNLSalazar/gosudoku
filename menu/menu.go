package menu

import (
	"fmt"
	"math"

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
	*a.selectedIndex = int(math.Abs(float64((*a.selectedIndex + n) % len(items))))
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
	return "Help page"
}
