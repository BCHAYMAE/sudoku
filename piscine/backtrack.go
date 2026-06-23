package piscine

func SolveSudoku() bool {
	//1:loop over all cells
	for i := 0; i < 9; i++ {//loop over rows
		for j := 0; j < 9; j++ {//loop over cols
			//2:find an empty cell
			if Table[i][j] == 0 {  //0 means the cell is empty
				//3: try palcing n from 1 to 9
				for num := 1; num <= 9; num++ {
					//4:check if the n can be placed by calling isvalid()
					if IsValid(i, j, num) {
						//we place it temporarily to see if it works
						// but we may remove it if it doesn’t lead to a solution
						Table[i][j] = num
						//6:recurse to solve the rest 
						//if a number works in an empty cell
						// the function keeps checking the rest recursively
						// and once a solution is found
						// it tells all the previous steps it worked
						if SolveSudoku() {
							return true
						}
						// backtrack : undo the placement
						Table[i][j] = 0
					}
				}
				//no valid n found
				return false
			}
		}
	}
	return true
}

//Backtracking is trying a number
//checking if it works, and undoing it if it doesn’t
//repeating this process until the Sudoku is solved.
