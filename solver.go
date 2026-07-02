package main

import (
	"fmt"
	"time"

	"github.com/DNLSalazar/gosudoku/sudoku"
)

func gamePrinter() (*chan bool, *chan sudoku.Sudoku) {
	t := time.NewTicker(time.Duration(time.Millisecond * 25))
	s := make(chan bool)
	gameChan := make(chan sudoku.Sudoku)
	go func() {
		var game sudoku.Sudoku
		for {
			select {
			case <-t.C:
				fmt.Print("\033[H\033[2J")
				game.PrintBoard()
			case <-s:
				t.Stop()
			case s := <-gameChan:
				game = s
			}
		}
	}()

	return &s, &gameChan
}

func backtrackSolver(s *sudoku.Sudoku, speed int) {
	if speed < 1 {
		speed = 1
	}

	var timeSleep int64

	switch speed {
	default:
		timeSleep = 0
	case 1:
		timeSleep = 30
	case 2:
		timeSleep = 15
	}
	var dfs func(x, y int) bool
	iter := 0

	stop, gameChan := gamePrinter()

	dfs = func(x, y int) bool {
		iter++
		if x == 9 {
			if s.IsValidBoard() {
				return true
			}
			return false
		}

		c := s.GetCell(sudoku.Coor{
			X: x,
			Y: y,
		})

		var nx, ny int

		if y == 8 {
			nx = x + 1
			ny = 0
		} else {
			ny = y + 1
			nx = x
		}

		if c.Static {
			return dfs(nx, ny)
		} else {
			for i := 1; i <= 9; i++ {
				s.ValidateNewCell(c.Coor, i)

				*gameChan <- *s
				if timeSleep != 0 {
					time.Sleep(time.Duration(time.Millisecond * time.Duration(timeSleep)))
				}
				if !c.HasErr {
					res := dfs(nx, ny)
					if res {
						return true
					}
				}
				s.ValidateNewCell(c.Coor, 0)
			}
		}
		return false
	}

	dfs(0, 0)
	*stop <- true

	solved := s.IsValidBoard()
	var result string = "Unsolved"

	if solved {
		result = "Solved"
	}

	fmt.Print("\033[H\033[2J")
	s.PrintBoard()
	fmt.Printf("\r\n\r\nThe board is %s. Number of excecutions for backtracking: %d\r\n\r\n", result, iter)
}
