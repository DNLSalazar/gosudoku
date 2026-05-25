package main

import (
	"fmt"

	flagargs "github.com/DNLSalazar/gosudoku/flagArgs"
	"github.com/DNLSalazar/gosudoku/game"
	"github.com/DNLSalazar/gosudoku/menu"
	"github.com/DNLSalazar/gosudoku/server"
	"github.com/DNLSalazar/gosudoku/sudoku"
)

func runGame(initialCells int) {
	fmt.Print("\033[H\033[2J")
	s := sudoku.CreateNewSudoku(initialCells)
	game := game.NewSudokuGame(&s)
	if _, err := game.Run(); err != nil {
		panic("Error running app")
	}
}

func startMenuApp() (int, bool) {
	fmt.Print("\033[H\033[2J")

	var selection int
	var quit bool
	a := menu.NewMenuApp(&selection, &quit)
	if _, err := a.Run(); err != nil {
		fmt.Println("Error running menu", err)
		panic("Error running menu")
	}

	return selection, quit
}

func main() {
	args := flagargs.GetArgs()
	var result int
	if args.Mode != 0 {
		result = args.Mode
	} else {
		resultApp, quit := startMenuApp()
		result = resultApp

		if quit {
			return
		}
	}

	switch result {
	case menu.ServeGame:
		fmt.Print("\033[H\033[2J")
		fmt.Println("Starting server game")
		server.Server(args.InitialCells)
	case menu.PlayOnTerminal:
		fmt.Println("Playing on terminal...")
		runGame(args.InitialCells)
	case menu.Solver:
		game := sudoku.CreateNewSudoku(args.InitialCells)
		backtrackSolver(&game, args.Speed)
	}
}
