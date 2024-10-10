package main

import (
	"fmt"
	"sudoku/sudoku"
)

func Solve(s *sudoku.Sudoku) {
	counter := 9
	for !s.IsValidBoard() && counter < 10 {
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				cell := s.GetCell(sudoku.Coor{X: i, Y: j})
				if cell.Static {
					continue
				}

				for value := 1; value < 10; value++ {
					result := s.ValidateNewCell(sudoku.Coor{X: i, Y: j}, value)
					if result {
						break
					}
				}
			}
		}
		counter++
	}

	fmt.Println(s.IsValidBoard())
}
