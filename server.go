package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"sync"

	"sudoku/piscine"
)

var solverMu sync.Mutex

var sampleRows = []string{
	"53..7....",
	"6..195...",
	".98....6.",
	"8...6...3",
	"4..8.3..1",
	"7...2...6",
	".6....28.",
	"...419..5",
	"....8..79",
}

var pageTemplate = template.Must(template.New("sudoku").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Sudoku</title>
	<style>
		* { box-sizing: border-box; }
		body {
			margin: 0;
			min-height: 100vh;
			display: grid;
			place-items: center;
			padding: 24px;
			font-family: Arial, sans-serif;
			background: #f4f1eb;
			color: #1f2933;
		}
		.panel {
			width: min(100%, 700px);
			background: white;
			border-radius: 18px;
			padding: 24px;
			box-shadow: 0 18px 40px rgba(0, 0, 0, 0.08);
		}
		h1 {
			margin: 0 0 8px;
		}
		p {
			margin: 0 0 18px;
			color: #52606d;
		}
		.board {
			display: grid;
			grid-template-columns: repeat(9, 1fr);
			border: 3px solid #243b53;
			margin-bottom: 16px;
		}
		.cell {
			width: 100%;
			aspect-ratio: 1;
			border: 1px solid #9fb3c8;
			text-align: center;
			font-size: 1.2rem;
			font-weight: 700;
		}
		.cell[readonly] {
			background: #d9e2ec;
			color: #243b53;
		}
		.wrong {
			background: #fde8e8;
			color: #b42318;
		}
		.top { border-top: 3px solid #243b53; }
		.left { border-left: 3px solid #243b53; }
		.actions {
			display: grid;
			grid-template-columns: repeat(4, 1fr);
			gap: 10px;
		}
		button {
			border: none;
			border-radius: 10px;
			padding: 12px;
			font: inherit;
			font-weight: 600;
			background: #243b53;
			color: white;
			cursor: pointer;
		}
		.status {
			margin-top: 16px;
			min-height: 24px;
			color: #52606d;
		}
		.error {
			color: #b42318;
		}
		@media (max-width: 620px) {
			.panel {
				padding: 16px;
			}
			.actions {
				grid-template-columns: 1fr;
			}
		}
	</style>
</head>
<body>
	<main class="panel">
		<h1>Sudoku</h1>
		<p>A simple interface for your Go solver.</p>

		<form method="post" action="/">
			{{range $index, $row := .BaseRows}}
				<input type="hidden" name="base{{$index}}" value="{{$row}}">
			{{end}}
			<div class="board">
				{{range .Cells}}
					{{range .}}
						<input
							class="cell {{.Class}}"
							type="number"
							inputmode="numeric"
							min="1"
							max="9"
							step="1"
							maxlength="1"
							name="{{.Name}}"
							value="{{.Value}}"
							{{if .ReadOnly}}readonly{{end}}>
					{{end}}
				{{end}}
			</div>

			<div class="actions">
				<button type="submit" name="action" value="sample">Load Sample</button>
				<button type="submit" name="action" value="clear">Clear</button>
				<button type="submit" name="action" value="check">Check</button>
				<button type="submit" name="action" value="solve">Solve</button>
			</div>
		</form>

		<div class="status {{if .Error}}error{{end}}">{{.Message}}</div>
	</main>
</body>
</html>
`))

type cellData struct {
	Name     string
	Value    string
	Class    string
	ReadOnly bool
}

type pageData struct {
	Cells    [][]cellData
	BaseRows []string
	Message  string
	Error    bool
}

func runServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleHome)

	fmt.Println("Open http://localhost:8080")
	return http.ListenAndServe(":8080", mux)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderPage(w, buildPage(sampleRows, sampleRows, "Fill the empty cells, then click Check.", false, nil))
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		renderPage(w, buildPage(sampleRows, sampleRows, "Could not read the form.", true, nil))
		return
	}

	action := r.FormValue("action")
	baseRows := baseRowsFromForm(r)

	switch action {
	case "sample":
		renderPage(w, buildPage(sampleRows, sampleRows, "Sample puzzle loaded.", false, nil))
	case "clear":
		renderPage(w, buildPage(baseRows, baseRows, "Board cleared.", false, nil))
	default:
		rows, invalidInput := rowsFromForm(r, baseRows)
		if invalidInput {
			renderPage(w, buildPage(baseRows, rows, "Use only numbers 1 to 9 in the board.", true, nil))
			return
		}

		switch action {
		case "check":
			checkAndRender(w, baseRows, rows)
		case "solve":
			solveAndRender(w, baseRows)
		default:
			renderPage(w, buildPage(baseRows, rows, "Unknown action.", true, nil))
		}
	}
}

func checkAndRender(w http.ResponseWriter, baseRows, rows []string) {
	if !boardIsFull(rows) {
		renderPage(w, buildPage(baseRows, rows, "Fill every empty cell before checking.", true, nil))
		return
	}

	solvedRows, ok := solveRows(baseRows)
	if !ok {
		renderPage(w, buildPage(baseRows, rows, "This puzzle could not be solved.", true, nil))
		return
	}

	wrongCells := wrongPositions(rows, solvedRows)
	if len(wrongCells) == 0 {
		renderPage(w, buildPage(baseRows, rows, "Perfect. You solved it.", false, nil))
		return
	}

	renderPage(w, buildPage(baseRows, rows, fmt.Sprintf("%d wrong cell(s) found.", len(wrongCells)), true, wrongCells))
}

func solveAndRender(w http.ResponseWriter, baseRows []string) {
	solvedRows, ok := solveRows(baseRows)
	if !ok {
		renderPage(w, buildPage(baseRows, baseRows, "No solution found.", true, nil))
		return
	}

	renderPage(w, buildPage(baseRows, solvedRows, "Solved.", false, nil))
}

func solveRows(rows []string) ([]string, bool) {
	args := append([]string{"sudoku"}, rows...)

	solverMu.Lock()
	defer solverMu.Unlock()

	if !piscine.Error(args) {
		return nil, false
	}

	piscine.FillTable(args)
	if !piscine.SolveSudoku() {
		return nil, false
	}

	return rowsFromTable(), true
}

func renderPage(w http.ResponseWriter, data pageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = pageTemplate.Execute(w, data)
}

func buildPage(baseRows, rows []string, message string, isError bool, wrongCells map[string]bool) pageData {
	cells := make([][]cellData, 9)

	for row := 0; row < 9; row++ {
		cells[row] = make([]cellData, 9)
		for col := 0; col < 9; col++ {
			value := ""
			if rows[row][col] != '.' {
				value = string(rows[row][col])
			}

			className := ""
			if row == 0 || row == 3 || row == 6 {
				className += "top "
			}
			if col == 0 || col == 3 || col == 6 {
				className += "left"
			}
			if wrongCells[fieldName(row, col)] {
				className += " wrong"
			}

			cells[row][col] = cellData{
				Name:     fieldName(row, col),
				Value:    value,
				Class:    strings.TrimSpace(className),
				ReadOnly: baseRows[row][col] != '.',
			}
		}
	}

	return pageData{
		Cells:    cells,
		BaseRows: append([]string(nil), baseRows...),
		Message:  message,
		Error:    isError,
	}
}

func rowsFromForm(r *http.Request, baseRows []string) ([]string, bool) {
	rows := make([]string, 9)
	hasInvalidInput := false

	for row := 0; row < 9; row++ {
		var builder strings.Builder
		for col := 0; col < 9; col++ {
			if baseRows[row][col] != '.' {
				builder.WriteByte(baseRows[row][col])
				continue
			}

			value := r.FormValue(fieldName(row, col))
			switch {
			case value == "":
				builder.WriteByte('.')
			case len(value) == 1 && value[0] >= '1' && value[0] <= '9':
				builder.WriteByte(value[0])
			default:
				builder.WriteByte('.')
				hasInvalidInput = true
			}
		}
		rows[row] = builder.String()
	}

	return rows, hasInvalidInput
}

func rowsFromTable() []string {
	rows := make([]string, 9)

	for row := 0; row < 9; row++ {
		var builder strings.Builder
		for col := 0; col < 9; col++ {
			builder.WriteByte(byte(piscine.Table[row][col] + '0'))
		}
		rows[row] = builder.String()
	}

	return rows
}

func emptyRows() []string {
	rows := make([]string, 9)
	for i := range rows {
		rows[i] = "........."
	}
	return rows
}

func baseRowsFromForm(r *http.Request) []string {
	rows := make([]string, 9)

	for row := 0; row < 9; row++ {
		value := r.FormValue(fmt.Sprintf("base%d", row))
		if len(value) == 9 {
			rows[row] = value
			continue
		}
		rows[row] = sampleRows[row]
	}

	return rows
}

func boardIsFull(rows []string) bool {
	for _, row := range rows {
		if strings.ContainsRune(row, '.') {
			return false
		}
	}
	return true
}

func wrongPositions(rows, solvedRows []string) map[string]bool {
	wrongCells := make(map[string]bool)

	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			if rows[row][col] != solvedRows[row][col] {
				wrongCells[fieldName(row, col)] = true
			}
		}
	}

	return wrongCells
}

func fieldName(row, col int) string {
	return fmt.Sprintf("r%dc%d", row, col)
}
