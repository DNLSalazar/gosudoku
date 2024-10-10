package sudoku

var centralCoors = []Coor{
	{1, 1},
	{1, 4},
	{1, 7},
	{4, 1},
	{4, 4},
	{4, 7},
	{7, 1},
	{7, 4},
	{7, 7},
}

func CreateNewSudoku(staticCells int) Sudoku {
	if staticCells > 81 {
		panic("Invalid static ceels")
	}
	board := make([][]Cell, 9)

	for i := range board {
		board[i] = make([]Cell, 9)
	}

	fillBoard(&board)

	s := Sudoku{
		board:       board,
		cCoors:      centralCoors,
		staticCells: staticCells,
	}

	s.emptyBoard()
	return s
}

func CreateSudokuFromCells(board [][]Cell) Sudoku {
	return Sudoku{board: board, cCoors: centralCoors}
}

func addToIntMap(m *map[int]int, value int) {
	if value != 0 {
		(*m)[value]++
	}
}

func validateIntMap(m *map[int]int) bool {
	for _, v := range *m {
		if v > 1 {
			return false
		}
	}
	return true
}

func validateIntMapForValue(m *map[int]int, value int) bool {
	if (*m)[value] > 1 {
		return false
	}
	return true
}

func fillBoard(board *[][]Cell) {
	*board = [][]Cell{
		{NC(1), NC(2), NC(3), NC(4), NC(5), NC(6), NC(7), NC(8), NC(9)},
		{NC(7), NC(8), NC(9), NC(1), NC(2), NC(3), NC(4), NC(5), NC(6)},
		{NC(4), NC(5), NC(6), NC(7), NC(8), NC(9), NC(1), NC(2), NC(3)},
		{NC(2), NC(3), NC(4), NC(5), NC(6), NC(7), NC(8), NC(9), NC(1)},
		{NC(8), NC(9), NC(1), NC(2), NC(3), NC(4), NC(5), NC(6), NC(7)},
		{NC(5), NC(6), NC(7), NC(8), NC(9), NC(1), NC(2), NC(3), NC(4)},
		{NC(3), NC(4), NC(5), NC(6), NC(7), NC(8), NC(9), NC(1), NC(2)},
		{NC(9), NC(1), NC(2), NC(3), NC(4), NC(5), NC(6), NC(7), NC(8)},
		{NC(6), NC(7), NC(8), NC(9), NC(1), NC(2), NC(3), NC(4), NC(5)},
	}
}

func addToMap(m *map[int]bool, value int) {
	if value != 0 {
		(*m)[value] = true
	}
}
