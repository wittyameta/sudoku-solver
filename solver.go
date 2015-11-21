package main

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/wittyameta/sudoku-solver/datatypes"
)

// TODO for all loops check if current cell needs to be excluded
const max int = 9

func main() {
	grid := *datatypes.InitGrid()
	count := 0
	for i := 0; i < max; i++ {
		count += readRow(&grid, i)
	}
	solve(&grid, count)
	grid.Print(0)
}

func readRow(grid *datatypes.Grid, rownum int) (count int) {
	var row [max]string
	format := ""
	for i := 0; i < max; i++ {
		format += "%s"
	}
	format += "\n"
	n, err := fmt.Scanf(format, &row[0], &row[1], &row[2], &row[3], &row[4], &row[5], &row[6], &row[7], &row[8])
	if n != max || err != nil {
		panic(err)
	}
	for i, elem := range row {
		val := &grid[rownum][i].IterationValues
		input := verifyElement(elem)
		if input > 0 {
			if (*val)[0].Possible[input] {
				(*val)[0] = *datatypes.SetValue(input)
				count++
			} else {
				panic("no solution possible")
			}
		}
	}
	return
}

func verifyElement(elem string) int {
	if "_" == elem {
		return 0
	}
	n, err := strconv.Atoi(elem)
	if err != nil {
		panic(err)
	}
	if n < 1 || n > max {
		panic("number should be from 1 to " + strconv.Itoa(max) + ".")
	} else {
		return n
	}
}

func solve(grid *datatypes.Grid, count int) {
	wg := sync.WaitGroup{}
	wg.Add(count)
	verificationCount := 0
	for i := 0; i < max; i++ {
		for j := 0; j < max; j++ {
			val := *grid[i][j].IterationValues[0].Val
			if val > 0 {
				if verificationCount < count {
					verificationCount++
				} else {
					wg.Add(1)
				}
				eliminateUsingGivenValues(grid, 0, i, j, val, &wg) // TODO goroutine
			}
		}
	}
	potentialCountDiff := count - verificationCount
	if potentialCountDiff > 0 {
		wg.Add(potentialCountDiff)
	}
	wg.Wait()
}

func eliminateUsingGivenValues(grid *datatypes.Grid, iteration int, row int, column int, val int, wg *sync.WaitGroup) {
	defer wg.Done()
	eliminatePossibilities(grid, 0, row, column, val)
}

// Called when a number is set in a cell.
// Eliminates the possibilities from the grid given a number at a row,col
// TODO eliminate only for row/col/block required
func eliminatePossibilities(grid *datatypes.Grid, iteration int, row int, column int, val int) {
	for i := 0; i < max; i++ {
		if i == row {
			continue
		}
		eliminatePossibilitiesForPosition(grid, iteration, i, column, val)
	}
	for j := 0; j < max; j++ {
		if j == column {
			continue
		}
		eliminatePossibilitiesForPosition(grid, iteration, row, j, val)
	}
	rowMin, columnMin := getBlockTopLeft(row, column)
	for i := rowMin; i < rowMin+3; i++ {
		for j := columnMin; j < columnMin+3; j++ {
			if i == row && j == column {
				continue
			}
			eliminatePossibilitiesForPosition(grid, iteration, i, j, val)
		}
	}
}

// Called when a number is set in a cell.
// Eliminates the possibilities from the grid given a number at a i,j
func eliminatePossibilitiesForPosition(grid *datatypes.Grid, iteration int, i int, j int, val int) {
	cell := &grid[i][j]
	cell.Mutex.Lock()
	setValue, updated, backtrack := updateCell(cell, iteration, val)
	cell.Mutex.Unlock()
	if setValue > 0 {
		eliminatePossibilities(grid, iteration, i, j, setValue)
	}
	if updated {
		fmt.Println("for", i, j, val)
		fmt.Println(cell.IterationValues[0].Possible, *cell.IterationValues[0].Val)
		fmt.Println(grid[i][j].IterationValues[0].Possible, *grid[i][j].IterationValues[0].Val)
		uniquePositions, conflict := checkIfUnique(grid, iteration, val, datatypes.Position{X: i, Y: j})
		if conflict {
			fmt.Println("backtrack logic") //TODO remove
			// TODO backtrack
		}
		for _, pos := range uniquePositions {
			setCell := &grid[pos.X][pos.Y]
			setCell.Mutex.Lock()
			eliminatedValues, isValueSet := setValueForCell(setCell, iteration, val)
			setCell.Mutex.Unlock()
			if isValueSet {
				fmt.Println("value is set", pos.X, pos.Y, val)
				fmt.Println(grid[pos.X][pos.Y].IterationValues[0].Possible, *grid[pos.X][pos.Y].IterationValues[0].Val)
				for _, eliminatedVal := range eliminatedValues {
					// check unique
					// TODO
					if eliminatedVal == 0 {
						panic("0")
					} // TODO remove
					checkIfUniqueAndEliminate(grid, iteration, pos.X, pos.Y, eliminatedVal)
					fmt.Println(eliminatedVal)
				}
				eliminatePossibilities(grid, iteration, pos.X, pos.Y, val)
			} else {
				fmt.Println("backtrack while set logic") //TODO remove
				// TODO backtrack
			}
		}
	}
	if backtrack {
		if iteration == 0 {
			panic("no solution possible")
		}
		// TODO backtrack
	}
}

// return valuesEliminated,isValueSet
func setValueForCell(cell *datatypes.Cell, iteration int, setValue int) (eliminatedValues []int, isValueSet bool) {
	existingValue := cell.IterationValues[iteration]
	if existingValue.Possible[setValue] && (*existingValue.Val == 0 || *existingValue.Val == setValue) {
		*existingValue.Val = setValue
		for key := range existingValue.Possible {
			if key != setValue {
				delete(existingValue.Possible, key)
				eliminatedValues = append(eliminatedValues, key)
			}
		}
		isValueSet = true
		return
	}
	isValueSet = false
	return
}

// Called as part of eliminatePossibilities
// Removes the value from the possibilities in the cell
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

func checkIfUniqueAndEliminate(grid *datatypes.Grid, iteration int, i int, j int, val int) {
	uniquePositions, conflict := checkIfUnique(grid, iteration, val, datatypes.Position{X: i, Y: j})
	if conflict {
		fmt.Println("backtrack logic", i, j, val) //TODO remove
		// TODO backtrack
		panic("yo") // remove
	}
	for _, pos := range uniquePositions {
		setCell := &grid[pos.X][pos.Y]
		setCell.Mutex.Lock()
		eliminatedValues, isValueSet := setValueForCell(setCell, iteration, val)
		setCell.Mutex.Unlock()
		if isValueSet {
			fmt.Println("value is set", pos.X, pos.Y, val)
			fmt.Println(grid[pos.X][pos.Y].IterationValues[0].Possible, *grid[pos.X][pos.Y].IterationValues[0].Val)
			for _, eliminatedVal := range eliminatedValues {
				// check unique
				// TODO
				checkIfUniqueAndEliminate(grid, iteration, pos.X, pos.Y, eliminatedVal)
				fmt.Println(eliminatedVal)
			}
			eliminatePossibilities(grid, iteration, pos.X, pos.Y, val)
		} else {
			fmt.Println("backtrack while set logic") //TODO remove
			// TODO backtrack
		}
	}
}

// TODO add support for row/col/block
// Called when a cell is updated
// If the value deleted now exists once in the row/block/col, then the cell is returned
// returns uniquePositions array, conflict
func checkIfUnique(grid *datatypes.Grid, iteration int, valDeleted int, pos datatypes.Position) ([]datatypes.Position, bool) {
	var uniquePositions []datatypes.Position
	var uniquePos datatypes.Position
	var foundUnique, atLeastOnce bool

	uniquePos, foundUnique, atLeastOnce = checkIfUniqueInRow(grid, iteration, valDeleted, pos)
	if foundUnique {
		uniquePositions = append(uniquePositions, uniquePos)
	} else if !atLeastOnce {
		return uniquePositions, true
	}
	uniquePos, foundUnique, atLeastOnce = checkIfUniqueInColumn(grid, iteration, valDeleted, pos)
	if foundUnique {
		uniquePositions = append(uniquePositions, uniquePos)
	} else if !atLeastOnce {
		return uniquePositions, true
	}
	uniquePos, foundUnique, atLeastOnce = checkIfUniqueInBlock(grid, iteration, valDeleted, pos)
	if foundUnique {
		uniquePositions = append(uniquePositions, uniquePos)
	} else if !atLeastOnce {
		return uniquePositions, true
	}
	return uniquePositions, false
}

// returns pos, foundUnique, atLeastOnce
func checkIfUniqueInRow(grid *datatypes.Grid, iteration int, valDeleted int, pos datatypes.Position) (datatypes.Position, bool, bool) {
	row := pos.X
	found := false
	column := max
	for j := 0; j < max; j++ {
		val := grid[row][j].IterationValues[iteration]
		if val.Possible[valDeleted] {
			if found || *val.Val > 0 {
				return pos, false, true
			}
			found = true
			column = j
		}
	}
	if found {
		return datatypes.Position{X: row, Y: column}, true, true
	}
	return pos, false, false
}

// returns pos, foundUnique, atLeastOnce
func checkIfUniqueInColumn(grid *datatypes.Grid, iteration int, valDeleted int, pos datatypes.Position) (datatypes.Position, bool, bool) {
	column := pos.Y
	found := false
	row := max
	for i := 0; i < max; i++ {
		val := grid[i][column].IterationValues[iteration]
		if val.Possible[valDeleted] {
			if found || *val.Val > 0 {
				return pos, false, true
			}
			found = true
			row = i
		}
	}
	if found {
		return datatypes.Position{X: row, Y: column}, true, true
	}
	return pos, false, false
}

// returns pos, foundUnique, atLeastOnce
func checkIfUniqueInBlock(grid *datatypes.Grid, iteration int, valDeleted int, pos datatypes.Position) (datatypes.Position, bool, bool) {
	row := pos.X
	column := pos.Y
	rowMin, columnMin := getBlockTopLeftFromPosition(pos)
	found := false
	for i := rowMin; i < rowMin+3; i++ {
		for j := columnMin; j < columnMin+3; j++ {
			val := grid[i][j].IterationValues[iteration]
			if val.Possible[valDeleted] {
				if found || *val.Val > 0 {
					return pos, false, true
				}
				found = true
				row = i
				column = j
			}
		}
	}
	if found {
		return datatypes.Position{X: row, Y: column}, true, true
	}
	return pos, false, false
}

func getBlockTopLeftFromPosition(pos datatypes.Position) (int, int) {
	row := pos.X
	column := pos.Y
	return getBlockTopLeft(row, column)
}

func getBlockTopLeft(x int, y int) (int, int) {
	return x - x%3, y - y%3
}
