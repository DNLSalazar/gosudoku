package main

import (
	"fmt"
	"sudoku/game"
	"sudoku/server"
	"sudoku/sudoku"
)

func testBoardWithFiles() {
	boards := readBoardsFromFile()
	prevSudoku, currSudoku := CreateSudokuFromCells(boards)
	prevSudoku.PrintBoard()
	prevSudoku.ValidateNewCell(sudoku.Coor{X: 5, Y: 0}, 9)
	currSudoku.PrintBoard()
}

func runGame() {
	fmt.Print("\033[H\033[2J")
	s := sudoku.CreateNewSudoku(17)
	game := game.NewSudokuGame(&s)
	if _, err := game.Run(); err != nil {
		panic("Error running app")
	}
}

func main() {
	s := sudoku.CreateNewSudoku(20)
	s.PrintBoard()
	// solving := NewSolvingSudoku(s)
	// solving.Solve()
	// solving.sudoku.PrintBoard()
	// Solve(&s)
	// runGame()
	server.Server()
}
