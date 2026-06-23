# Sudoku

A small Go Sudoku project with:

- a backtracking solver
- a simple browser interface rendered by Go
- the command-line solver mode

## Run The Interface

Start the app:

```bash
go run .
```

Then open:

```text
http://localhost:8080/
```

## Interface Features

- sample Sudoku loaded by default
- original clue cells are locked
- editable cells accept numbers `1` to `9`
- `Check` highlights wrong cells after the board is full
- `Solve` fills the board using the Go solver
- `Clear` resets the editable cells

## Run In CLI Mode

You can also solve a Sudoku directly from the terminal by passing 9 rows:

```bash
go run . "..9748..." "7........" ".2.1.9..." "..7...24." ".64.1.59." ".98...3.." "...8.3.2." "........6" "...2759.."
```

Use `.` for empty cells.

## Project Files

- `main.go`: starts the web interface or CLI solver
- `server.go`: simple HTML form interface and game checks
- `piscine/backtrack.go`: backtracking solver
- `piscine/check.go`: Sudoku validation rules
- `piscine/table.go`: board loading and printing
