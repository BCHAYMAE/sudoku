package main

import (
	"fmt"
	"os"
	"sudoku/piscine"
)

func main() {
	if len(os.Args) == 1 {
		if err := runServer(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return
	}

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
