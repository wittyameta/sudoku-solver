package main

import (
	"fmt"
	"strconv"

	"github.com/wittyameta/sudoku-solver/datatypes"
)

func main() {
	grid := *datatypes.InitGrid()
	for i := 0; i < 9; i++ {
		readRow(&grid, i)
	}
	grid.Print(0) // remove
	solve(&grid)
	grid.Print(0)
}

func readRow(grid *datatypes.Grid, rownum int) {
	var row [9]string
	format := ""
	for i := 0; i < 9; i++ {
		format += "%s"
	}
	format += "\n"
	n, err := fmt.Scanf(format, &row[0], &row[1], &row[2], &row[3], &row[4], &row[5], &row[6], &row[7], &row[8])
	if n != 9 || err != nil {
		panic(err)
	}
	for i, elem := range row {
		val := &grid[rownum][i].IterationValues
		input := verifyElement(elem)
		if input > 0 {
			if (*val)[0].Possible[input] {
				(*val)[0] = *datatypes.SetValue(input)
			}
		}
	}
}

func verifyElement(elem string) int {
	if "_" == elem {
		return 0
	}
	n, err := strconv.Atoi(elem)
	if err != nil {
		panic(err)
	}
	if n < 1 || n > 9 {
		panic("number should be from 1 to 9.")
	} else {
		return n
	}
}

func solve(grid *datatypes.Grid) {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			val := *grid[i][j].IterationValues[0].Val
			if val > 0 {
				eliminatePossibilities(grid, 0, i, j, val)
			}
		}
	}
}

func eliminatePossibilities(grid *datatypes.Grid, iteration int, row int, column int, val int) {
	for i := 0; i < 9; i++ {
		if i == row {
			continue
		}
		// TODO currently for row, check for col and block as well
		cell := &grid[i][column]
		// TODO lock on cell
		setValue, updated, backtrack := updateCell(cell, iteration, val)
		// TODO unlock cell
		if setValue > 0 {
			eliminatePossibilities(grid, iteration, i, column, setValue)
		}
		if updated {
			uniquePositions := checkIfUnique(grid, iteration, val)
			for _, pos := range uniquePositions {
				// TODO
				// find cell
				//lock cell
				// update cell - set value
				// unlock cell
			}
		}
		if backtrack {
			if iteration == 0 {
				panic("no solution possible")
			}
			// TODO
		}
	}
}

// returns setValue,updated,backtrack (int,bool,bool)
func updateCell(cell *datatypes.Cell, iteration int, valToDelete int) (int, bool, bool) {
	existingValue := cell.IterationValues[iteration]
	if *existingValue.Val == valToDelete {
		return 0, false, true
	}
	updated := false
	setValue := 0
	if *existingValue.Val == 0 && existingValue.Possible[valToDelete] {
		updated = true
		delete(existingValue.Possible, valToDelete)
		if len(existingValue.Possible) == 1 {
			for key := range existingValue.Possible {
				setValue = key
				*existingValue.Val = key
			}
		}
	}
	return setValue, updated, false
}

// TODO add support for row/col/block
func checkIfUnique(grid *datatypes.Grid, iteration int, valDeleted int) []datatypes.Position {
	// TODO
	var uniquePositions []datatypes.Position
	return uniquePositions
}
