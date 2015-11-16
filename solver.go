package main

import (
	"fmt"
	"strconv"

	"github.com/wittyameta/sudoku-solver/datatypes"
)

func main() {
	val := *datatypes.SetValue(1)
	fmt.Println("done", val)
	grid := *datatypes.InitGrid()
	for i := 0; i < 9; i++ {
		readRow(&grid, i)
		fmt.Println(grid[i]) // remove
		grid.Print(0) // remove
	}
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
		initValue := grid[rownum][i].IterationValues[0]
		input := verifyElement(elem)
		if input > 0 {
			if initValue.Possible[input] {
				grid[rownum][i].IterationValues[0] = *datatypes.SetValue(input)
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

func eliminatePossibilities(grid *datatypes.Grid, iteration int, row int, column int) {

}
