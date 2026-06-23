package piscine

import "fmt"

var Table [9][9]int // 9*9 to store the sudoku grid

func Error(args []string) bool {
	if len(args) != 10 { // we expect 10 grid one for the name and the rest for the grid if not return false
		fmt.Println("Error")
		return false
	}
	for i := 1; i <= 9; i++ { // loop through each row 
		r := args[i]
		if len(r) != 9 { //if the len is diff then 9 then false cause the grid should have 9 
			fmt.Println("Error")
			return false
		}
		for x := 0; x < 9; x++ { // now we loop through each char i the row if its diff then nbr or . then false
			c := r[x]
			if c != '.' && (c < '1' || c > '9') {
				fmt.Println("Error")
				return false
			}
			for y := x + 1; y < 9; y++ { // now we check if we have duplicate char in row if yes then false
				if c != '.' && c == r[y] {
					fmt.Println("Error")
					return false
				}
			}
		}
	}
	return true
}

func FillTable(a []string) { //convert each valid arg to int 
	Table = [9][9]int{}
	for i := 0; i < 9; i++ { //loop through rows 0 to 8 
		for j, c := range a[i+1] { //loop through each char
			if c == '.' {
				Table[i][j] = 0
			} else {
				Table[i][j] = int(c - '0')
			}
		}
	}
}

func PrintTable() {
	for _, r := range Table { //loop through each row
		for j, v := range r { //loop through each cell in the row
			fmt.Print(v) // print the cell value
			if j != 8 { //add space between each cell except the last one
				fmt.Print(" ")
			}
		}
		fmt.Println()//after finishing the row move to the next one 
	}
}
