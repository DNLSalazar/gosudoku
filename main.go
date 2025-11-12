package main

import (
	"fmt"
	"sudoku/game"
	"sudoku/menu"
	"sudoku/server"
	"sudoku/sudoku"
)

func runGame() {
	fmt.Print("\033[H\033[2J")
	s := sudoku.CreateNewSudoku(17)
	game := game.NewSudokuGame(&s)
	if _, err := game.Run(); err != nil {
		panic("Error running app")
	}
}

func startMenuApp() int {
	fmt.Print("\033[H\033[2J")

	var selection int
	a := menu.NewMenuApp(&selection)
	if _, err := a.Run(); err != nil {
		fmt.Println("Error running menu", err)
		panic("Error running menu")
	}

	return selection
}

func main() {
	result := startMenuApp()
	switch result {
	case menu.ServeGame:
		fmt.Println("Starting server game")
		server.Server()
	case menu.PlayOnTerminal:
		fmt.Println("Playing on terminal...")
		runGame()
	}
}
