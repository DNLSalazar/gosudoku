package sudoku

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"strings"

	"github.com/fatih/color"
)

var colorBlue = color.New(color.FgBlue)
var colorGreen = color.New(color.FgCyan)

type Coor struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Cell struct {
	Value  int  `json:"value"`
	Static bool `json:"static"`
	Coor   Coor `json:"coor"`
	HasErr bool `json:"hasError"`
}

type DumbCell struct {
	Coor     Coor `json:"coor"`
	HasErr   bool `json:"hasError"`
	HasValue bool `json:"hasValue"`
}

func NC(n int, x, y int) Cell {
	return Cell{
		Value:  n,
		Static: true,
		Coor: Coor{
			x,
			y,
		},
	}
}

type Sudoku struct {
	board       [][]Cell
	cCoors      []Coor
	staticCells int
	HasErr      bool
	errorCells  []*Cell
}

func (s *Sudoku) GetCell(c Coor) *Cell {
	return &s.board[c.X][c.Y]
}

func (s *Sudoku) GetBoard() [][]Cell {
	return s.board
}

func (s *Sudoku) PrintStateOfBoard() {
	marshallResult, err := json.Marshal(s.board)
	if err != nil {
		panic("Error parsing sudoku")
	}
	fmt.Println(string(marshallResult))
}

func (s *Sudoku) ValidateCells() {
	s.HasErr = false
	for i := range 9 {
		for k := range 9 {
			s.board[i][k].HasErr = false
		}
	}
	for i := range 9 {
		rowCells := make([]*Cell, 0)
		colCells := make([]*Cell, 0)
		for j := range 9 {
			rc := &(s.board[i][j])
			if rc.Value != 0 {
				rowCells = append(rowCells, rc)
			}

			cc := &(s.board[j][i])
			if cc.Value != 0 {
				colCells = append(colCells, cc)
			}
		}
		validateSliceOfCells(&rowCells)
		validateSliceOfCells(&colCells)
	}

	s.validateCuadrantsCells()
}

func (s *Sudoku) validateCuadrantsCells() {
	for _, v := range s.cCoors {
		initX := v.X - 1
		initY := v.Y - 1

		cuadrant := make([]*Cell, 0)
		for i := initX; i <= v.X+1; i++ {
			for j := initY; j <= v.Y+1; j++ {
				cell := &(s.board[i][j])
				if cell.Value != 0 {
					cuadrant = append(cuadrant, cell)
				}
			}
		}
		validateSliceOfCells(&cuadrant)
	}
}

func (s *Sudoku) IsValidBoard() bool {
	for i := range 9 {
		cell := s.board[i][i]
		if cell.Value == 0 {
			return false
		}
		col := make(map[int]bool)
		col[cell.Value] = true
		row := make(map[int]bool)
		row[cell.Value] = true
		for j := range 9 {
			addToMap(&col, s.board[i][j].Value)
			addToMap(&row, s.board[j][i].Value)
		}
		if len(col) != 9 || len(row) != 9 {
			return false
		}
	}
	return s.validateAllCuadrants()
}

func (s *Sudoku) validateAllCuadrants() bool {
	for _, v := range s.cCoors {
		initX := v.X - 1
		initY := v.Y - 1

		cuadrant := make(map[int]bool)
		for i := initX; i <= v.X+1; i++ {
			for j := initY; j <= v.Y+1; j++ {
				addToMap(&cuadrant, s.board[i][j].Value)
			}
		}

		if len(cuadrant) != 9 {
			return false
		}
	}
	return true
}

func (s *Sudoku) emptyBoard() {
	cells := make([]*Cell, 0)
	for i := 0; i < len(s.board); i++ {
		for j := 0; j < len(s.board[i]); j++ {
			cells = append(cells, &(s.board[i][j]))
		}
	}

	for i := range cells {
		j := rand.IntN(i + 1)
		cells[i], cells[j] = cells[j], cells[i]
	}

	for len(cells) > s.staticCells {
		for i := range cells {
			num := rand.IntN(100)
			if num < 50 {
				cells[i].Value = 0
				cells[i].Static = false
				cells = slices.Delete(cells, i, i+1)
				break
			}
			if len(cells) == s.staticCells {
				break
			}
		}
	}
}

func (s *Sudoku) PrintBoard() {
	fmt.Printf("     | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | \r\n\r\n")
	fmt.Printf("     ------------------------------------- \r\n")
	for i := range 9 {
		colorNumRow := (int(math.Floor(float64(i)/3))%3)%2 == 0
		var str strings.Builder
		fmt.Fprintf(&str, "  %d  ", i+1)
		for j := range 9 {
			colorNumCol := (int(math.Floor(float64(j)/3))%3)%2 == 0

			var color *color.Color
			if colorNumRow == colorNumCol {
				color = colorBlue
			} else {
				color = colorGreen
			}

			if j == 0 {
				str.WriteString(color.Sprintf("|"))
			}
			cell := s.board[i][j]
			var value string
			if cell.Value == 0 {
				value = " "
			} else {
				value = fmt.Sprint(cell.Value)
			}
			str.WriteString(color.Sprintf(" %s |", value))
		}
		fmt.Printf("%s   %d\r\n", str.String(), i+1)
		fmt.Printf("     ------------------------------------- \r\n")
	}
	fmt.Printf("\r\n     | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | \r\n")
}

func (s *Sudoku) ValidateCuadrantGame(center Coor) bool {
	initX := center.X - 1
	initY := center.Y - 1

	cuadrant := make(map[int]int)
	for i := initX; i <= center.X+1; i++ {
		for j := initY; j <= center.Y+1; j++ {
			addToIntMap(&cuadrant, s.board[i][j].Value)
		}
	}
	return validateIntMap(&cuadrant)
}

// NOTE: Unused
func (s *Sudoku) cellHasError(cell *Cell) bool {
	col := make(map[int]int)
	row := make(map[int]int)
	value := cell.Value
	c := cell.Coor

	for i := range 9 {
		addToIntMap(&col, s.board[c.X][i].Value)
		addToIntMap(&row, s.board[i][c.Y].Value)
	}

	validRow, validCol := validateIntMapForValue(&row, value), validateIntMapForValue(&col, value)
	center := s.GetCenterOfCuadrantForCoor(c)
	validCuadrant := s.ValidateCuadrantGame(center)
	if !validRow || !validCol || !validCuadrant {
		cell.HasErr = true
		s.errorCells = append(s.errorCells, cell)
		return true
	}
	cell.HasErr = false
	return false
}

func (s *Sudoku) ValidateNewCell(c Coor, value int) [][]Cell {
	cell := &(s.board[c.X][c.Y])
	if (*cell).Static {
		return s.board
	}

	cell.Value = value

	s.ValidateCells()
	return s.board
}

func (s *Sudoku) TestNewCell(c Coor, value int) bool {
	cell := s.GetCell(c)
	return s.cellHasError(cell)
}

func (s *Sudoku) GetCenterOfCuadrantForCoor(c Coor) Coor {
	for _, v := range s.cCoors {
		if (c.X == v.X-1 || c.X == v.X || c.X == v.X+1) &&
			(c.Y == v.Y-1 || c.Y == v.Y || c.Y == v.Y+1) {
			return v
		}
	}
	panic("NO CENTER FOUND")
}

func (s *Sudoku) GetDumbBoard() [][]DumbCell {
	newBoard := make([][]DumbCell, len(s.board))

	for i, row := range s.board {
		newRow := make([]DumbCell, len(row))

		for k, cell := range row {
			dc := DumbCell{
				Coor:     cell.Coor,
				HasErr:   cell.HasErr,
				HasValue: cell.Value != 0,
			}
			newRow[k] = dc
		}

		newBoard[i] = newRow
	}

	return newBoard
}
