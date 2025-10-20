package piscine

// now will check whether n can be placed in the col or row without breaking the rules 


func IsValid(row, col, n int) bool {
	//check the row and col 
	for i := 0; i < 9; i++ { // loop through all cells in the same row or col
		if Table[row][i] == n || Table[i][col] == n { //if the same n already exist in the same row or col
			return false 
		}
	}
	//find the starting of coordinayes
	startRow := (row / 3) * 3
	startCol := (col / 3) * 3

	for i := 0; i < 3; i++ { //loop over 3 row in 3*3
		for j := 0; j < 3; j++ { //loop over 3 col in 3*3
			if Table[startRow+i][startCol+j] == n { //if n exist anywher inside 3*3
				return false
			}
		}
	}
	return true
}
