package main

import (
	"fmt"
	"sudoku/game"
	"sudoku/sudoku"
)

func testBoardWithFiles() {
	boards := readBoardsFromFile()
	prevSudoku, currSudoku := CreateSudokuFromCells(boards)
	prevSudoku.PrintBoard()
	result := prevSudoku.ValidateNewCell(sudoku.Coor{X: 5, Y: 0}, 9)
	fmt.Println(result)
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
	// s := sudoku.CreateNewSudoku(17)
	// nums := make([][]int, 9)
	// for i := 0; i < 9; i++ {
	// 	nums[i] = make([]int, 9)
	// 	for j := 0; j < 9; j++ {
	// 		nums[i][j] = s.Board[i][j].Value
	// 	}
	// }
	// fmt.Println(s.IsValidBoard())
	// fmt.Println(nums)
	// s.PrintBoard()
	// Solve(&s)
	// s.PrintBoard()
	runGame()
}
