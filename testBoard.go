package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/DNLSalazar/gosudoku/sudoku"
)

func readBoardsFromFile() [][][]sudoku.Cell {
	var boards [][][]sudoku.Cell

	content, err := os.ReadFile("./testBoards.txt")
	if err != nil {
		log.Println("Error reading file")
		panic(err)
	}

	err = json.Unmarshal(content, &boards)
	if err != nil {
		log.Println("Error on unmarshal")
		panic(err)
	}

	return boards
}

func CreateSudokuFromCells(boards [][][]sudoku.Cell) (sudoku.Sudoku, sudoku.Sudoku) {
	return sudoku.CreateSudokuFromCells(boards[0]), sudoku.CreateSudokuFromCells(boards[1])
}
