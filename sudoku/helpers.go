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
	newBoard := make([][]Cell, len(board))

	for i := range board {
		row := board[i]
		newRow := make([]Cell, len(row))

		for k := range row {
			cell := board[i][k]
			newCell := Cell{
				Value:  cell.Value,
				HasErr: cell.HasErr,
				Static: cell.Static,
				Coor: Coor{
					X: cell.Coor.X,
					Y: cell.Coor.Y,
				},
			}

			newRow[k] = newCell
		}
		newBoard[i] = newRow
	}
	return Sudoku{board: newBoard, cCoors: centralCoors}
}

func addToIntMap(m *map[int]int, value int) {
	if value != 0 {
		(*m)[value]++
	}
}

func validateSliceOfCells(cells *[]*Cell) {
	m := make(map[int][]*Cell)
	for i := 0; i < len(*cells); i++ {
		cell := (*cells)[i]
		if _, ok := m[cell.Value]; ok {
			m[cell.Value] = append(m[cell.Value], cell)
		} else {
			m[cell.Value] = []*Cell{cell}
		}
	}

	for k := range m {
		l := len(m[k])
		if l > 1 {
			for i := range l {
				m[k][i].HasErr = true
			}
		}
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
		{NC(1, 0, 0), NC(2, 0, 1), NC(3, 0, 2), NC(4, 0, 3), NC(5, 0, 4), NC(6, 0, 5), NC(7, 0, 6), NC(8, 0, 7), NC(9, 0, 8)},
		{NC(7, 1, 0), NC(8, 1, 1), NC(9, 1, 2), NC(1, 1, 3), NC(2, 1, 4), NC(3, 1, 5), NC(4, 1, 6), NC(5, 1, 7), NC(6, 1, 8)},
		{NC(4, 2, 0), NC(5, 2, 1), NC(6, 2, 2), NC(7, 2, 3), NC(8, 2, 4), NC(9, 2, 5), NC(1, 2, 6), NC(2, 2, 7), NC(3, 2, 8)},
		{NC(2, 3, 0), NC(3, 3, 1), NC(4, 3, 2), NC(5, 3, 3), NC(6, 3, 4), NC(7, 3, 5), NC(8, 3, 6), NC(9, 3, 7), NC(1, 3, 8)},
		{NC(8, 4, 0), NC(9, 4, 1), NC(1, 4, 2), NC(2, 4, 3), NC(3, 4, 4), NC(4, 4, 5), NC(5, 4, 6), NC(6, 4, 7), NC(7, 4, 8)},
		{NC(5, 5, 0), NC(6, 5, 1), NC(7, 5, 2), NC(8, 5, 3), NC(9, 5, 4), NC(1, 5, 5), NC(2, 5, 6), NC(3, 5, 7), NC(4, 5, 8)},
		{NC(3, 6, 0), NC(4, 6, 1), NC(5, 6, 2), NC(6, 6, 3), NC(7, 6, 4), NC(8, 6, 5), NC(9, 6, 6), NC(1, 6, 7), NC(2, 6, 8)},
		{NC(9, 7, 0), NC(1, 7, 1), NC(2, 7, 2), NC(3, 7, 3), NC(4, 7, 4), NC(5, 7, 5), NC(6, 7, 6), NC(7, 7, 7), NC(8, 7, 8)},
		{NC(6, 8, 0), NC(7, 8, 1), NC(8, 8, 2), NC(9, 8, 3), NC(1, 8, 4), NC(2, 8, 5), NC(3, 8, 6), NC(4, 8, 7), NC(5, 8, 8)},
	}
}

func addToMap(m *map[int]bool, value int) {
	if value != 0 {
		(*m)[value] = true
	}
}
