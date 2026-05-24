package main

import (
	"fmt"
	"slices"

	"github.com/DNLSalazar/gosudoku/sudoku"
)

type Solvingcell struct {
	Cell             sudoku.Cell
	InitPosibilities []int
	Posibilities     []int
}

func NewSolvingCell(c sudoku.Cell) Solvingcell {
	return Solvingcell{
		Cell:             c,
		InitPosibilities: make([]int, 0),
		Posibilities:     make([]int, 0),
	}
}

type SolvingSudoku struct {
	sudoku       *sudoku.Sudoku
	Board        [][]Solvingcell
	Posibilities map[int][]*Solvingcell
}

func NewSolvingSudoku(s sudoku.Sudoku) SolvingSudoku {
	board := make([][]Solvingcell, 9)
	sudokuBoard := s.GetBoard()
	posibilities := make(map[int][]*Solvingcell)
	for i := range 9 {
		board[i] = make([]Solvingcell, 9)
		posibilities[i+1] = make([]*Solvingcell, 0)
		for j := range 9 {
			board[i][j] = NewSolvingCell(sudokuBoard[i][j])
		}
	}
	solving := SolvingSudoku{
		sudoku:       &s,
		Board:        board,
		Posibilities: posibilities,
	}

	solving.getInitialPosibilities()
	return solving
}

func (s *SolvingSudoku) getInitialPosibilities() {
	for i := range 9 {
		for j := range 9 {
			cell := &(s.Board[i][j])
			if cell.Cell.Static {
				continue
			}
			for v := 1; v < 10; v++ {
				result := s.sudoku.TestNewCell(cell.Cell.Coor, v)
				if !result {
					cell.Posibilities = append(cell.Posibilities, v)
				}
			}
			cell.InitPosibilities = slices.Clone(cell.Posibilities)
			s.Posibilities[len(cell.InitPosibilities)] = append(s.Posibilities[len(cell.InitPosibilities)], cell)
		}
	}
}

func (s *SolvingSudoku) getNewPosibilities() {
	for i := range 9 {
		s.Posibilities[i+1] = make([]*Solvingcell, 0)
	}
	for i := range 9 {
		for j := range 9 {
			cell := &(s.Board[i][j])
			if cell.Cell.Static || (cell.Cell.Value != 0 && !cell.Cell.HasErr) {
				continue
			}
			cell.Posibilities = make([]int, 0)
			for v := 1; v < 10; v++ {
				result := s.sudoku.TestNewCell(cell.Cell.Coor, v)
				// fmt.Println("ADDING POSIVILITY", v, cell.Cell.Coor)
				if !result {
					cell.Posibilities = append(cell.Posibilities, v)
				}
			}
			s.Posibilities[len(cell.Posibilities)] = append(s.Posibilities[len(cell.Posibilities)], cell)
		}
	}
}

func (s *SolvingSudoku) solveForLeastPosibilities() {
	for i := range 9 {
		n := i + 1
		p := s.Posibilities[n]
		if p == nil {
			continue
		}
		l := len(p)
		if l > 0 {
			for j := range l {
				cell := p[j]
				solvedCell := s.sudoku.ValidateNewCell(cell.Cell.Coor, cell.Posibilities[0])[cell.Cell.Coor.X][cell.Cell.Coor.Y]
				fmt.Println("VALIDATING on SOLVE", cell.Cell.Coor, cell.Posibilities[0])
				if solvedCell.HasErr {
					fmt.Printf("%v, %d", solvedCell, cell.Posibilities[0])
					panic("CELL CANNOT HAVE ERROR ON SOLVING FOR POSIBILITIES")
				}
				(*cell).Cell = solvedCell
			}
			break
		}
	}
}

func (s *SolvingSudoku) solveOnePosibleCell() {
	for i := range 9 {
		n := i + 1
		p := s.Posibilities[n]
		if p == nil {
			continue
		}
		l := len(p)
		if l > 0 {
			cell := p[0]
			solvedCell := s.sudoku.ValidateNewCell(cell.Cell.Coor, cell.Posibilities[0])[cell.Cell.Coor.X][cell.Cell.Coor.Y]
			fmt.Println("VALIDATING on SOLVE", cell.Cell.Coor, cell.Posibilities[0], len(cell.Posibilities))
			if solvedCell.HasErr {
				fmt.Printf("%v, %d", solvedCell, cell.Posibilities[0])
				panic("CELL CANNOT HAVE ERROR ON SOLVING FOR POSIBILITIES")
			}
			(*cell).Cell = solvedCell
			break
		}
	}
}

func (s *SolvingSudoku) Solve() {
	counter := 0
	for !s.sudoku.IsValidBoard() {
		s.solveOnePosibleCell()
		s.getNewPosibilities()
		counter++
	}
}

func Solve(s *sudoku.Sudoku) {
	s.PrintBoard()
	counter := 9
	for !s.IsValidBoard() && counter < 10 {
		for i := range 9 {
			for j := range 9 {
				cell := s.GetCell(sudoku.Coor{X: i, Y: j})
				if cell.Static {
					continue
				}

				for value := 1; value < 10; value++ {
					result := s.ValidateNewCell(sudoku.Coor{X: i, Y: j}, value)[i][j]
					if !result.HasErr {
						break
					}
					s.ValidateNewCell(sudoku.Coor{X: i, Y: j}, 0)
				}
			}
		}
		counter++
	}

	s.PrintBoard()
	fmt.Println(s.IsValidBoard())
}

func backtrackSolver(s *sudoku.Sudoku) {
	var dfs func(x, y int) bool
	iter := 0

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
				fmt.Print("\033[H\033[2J")

				s.PrintBoard()
				// time.Sleep(time.Duration(time.Millisecond * 10))
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

	solved := s.IsValidBoard()
	var result string = "Unsolved"

	if solved {
		result = "Solved"
	}

	fmt.Printf("\r\n\r\nThe board is %s. The number of excecutions for backtracking %d\r\n\r\n", result, iter)
}
