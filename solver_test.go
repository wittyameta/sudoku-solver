// author: Jayant Ameta
// https://github.com/wittyameta

package main

import (
	"github.com/wittyameta/sudoku-solver/datatypes"
	"testing"
)

// TestVerifyElement verifies the input integer.
func TestVerifyElement(t *testing.T) {
	v := verifyElement("3")
	if v != 3 {
		t.Error("Expected 3, got ", v)
	}
}

// TestSolve initializes and solves the grid.
func TestSolve(t *testing.T) {
	grid := *datatypes.InitGrid()
	count := setInput(&grid)
	positions := solve(&grid, count)
	if len(positions) < 37 {
		t.Error("Expected at least 37, got ", len(positions))
	}
	solveByGuessing(&grid, positions, 0)
}

func BenchmarkSolve(b *testing.B) {
	for n := 0; n < b.N; n++ {
		grid := *datatypes.InitGrid()
		count := setInput(&grid)
		positions := solve(&grid, count)
		if len(positions) != 37 {
			b.Error("Expected 37, got ", len(positions))
		}
		solveByGuessing(&grid, positions, 0)
	}
}

// setInput creates the initial grid for testing.
func setInput(grid *datatypes.Grid) int {
	setValue(grid, 0, 4, 4)
	setValue(grid, 0, 5, 5)

	setValue(grid, 1, 0, 8)
	setValue(grid, 1, 6, 2)
	setValue(grid, 1, 8, 7)

	setValue(grid, 2, 2, 2)
	setValue(grid, 2, 8, 4)

	setValue(grid, 3, 2, 6)
	setValue(grid, 3, 6, 3)
	setValue(grid, 3, 8, 2)

	setValue(grid, 4, 3, 1)

	setValue(grid, 5, 0, 2)
	setValue(grid, 5, 2, 7)
	setValue(grid, 5, 3, 4)
	setValue(grid, 5, 6, 6)

	setValue(grid, 6, 0, 6)
	setValue(grid, 6, 1, 4)
	setValue(grid, 6, 4, 9)
	setValue(grid, 6, 5, 8)

	setValue(grid, 7, 0, 7)
	setValue(grid, 7, 1, 9)
	setValue(grid, 7, 5, 4)

	setValue(grid, 8, 7, 3)

	grid.Print()
	return 23
}

func setValue(grid *datatypes.Grid, row int, column int, val int) {
	*grid[row][column].Val = val
	grid[row][column].IterationValues[0] = *datatypes.SetValue(val)
}
