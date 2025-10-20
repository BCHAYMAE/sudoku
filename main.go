package main

import (
	"fmt"
	"os"
	"sudoku/piscine"
)

func main() {
	if !piscine.Error(os.Args) {
		return
	}

	piscine.FillTable(os.Args)

	if piscine.SolveSudoku() {
		piscine.PrintTable()
	} else {
		fmt.Println("No solution found")
	}
}
