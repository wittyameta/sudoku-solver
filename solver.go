package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/wittyameta/sudoku-solver/datatypes"
)

// TODO
// for all loops check if current cell needs to be excluded
// function return instead of value
// when num is set, check if after elimination in nearby block, the num is specific to a row/col
const max int = 9
const easy, medium, hard string = "easy", "medium", "hard"

var numSolutions int

func main() {
	numSolutions = 0
	grid := *datatypes.InitGrid()
	count := 0
	inputValues := make(map[int]bool)
	for i := 0; i < max; i++ {
		count += readRow(&grid, i, inputValues)
	}
	if count < 17 || len(inputValues) < 8 {
		handleError("Too few input values given. At least 17 values, and 8 distinct values must be given.", nil)
	}
	positions := solve(&grid, count)
	solveByGuessing(&grid, positions, 0)
	fmt.Println("Total solutions:", numSolutions)
	difficultyLevel := easy
	if len(positions) > 0 {
		if len(positions) < max {
			difficultyLevel = medium
		} else {
			difficultyLevel = hard
		}
	}
	fmt.Println("Difficulty level:", difficultyLevel)
}

func readRow(grid *datatypes.Grid, rownum int, inputValues map[int]bool) (count int) {
	var row [max]string
	format := ""
	for i := 0; i < max; i++ {
		format += "%s"
	}
	format += "\n"
	n, err := fmt.Scanf(format, &row[0], &row[1], &row[2], &row[3], &row[4], &row[5], &row[6], &row[7], &row[8])
	if n != max || err != nil {
		handleError("", err)
	}
	for i, elem := range row {
		val := grid[rownum][i].IterationValues
		input := verifyElement(elem)
		if input > 0 {
			if val[0].Possible[input] {
				val[0] = *datatypes.SetValue(input)
				*grid[rownum][i].Val = input
				count++
				inputValues[input] = true
			} else {
				handleError("no solution possible", nil)
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
		handleError("", err)
	}
	if n < 1 || n > max {
		handleError("number should be from 1 to "+strconv.Itoa(max)+".", nil)
	}
	return n
}

func solve(grid *datatypes.Grid, count int) map[datatypes.Position]bool {
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
				go initialElimination(grid, i, j, val, &wg)
			}
		}
	}
	potentialCountDiff := count - verificationCount
	if potentialCountDiff > 0 {
		wg.Add(potentialCountDiff)
	}
	wg.Wait()
	return initPositions(grid)
}

func initialElimination(grid *datatypes.Grid, row int, column int, val int, wg *sync.WaitGroup) {
	defer wg.Done()
	if eliminateUsingGivenValues(grid, 0, row, column, val) {
		handleError("No solution", nil)
	}
}

func eliminateUsingGivenValues(grid *datatypes.Grid, iteration int, row int, column int, val int) bool {
	if eliminatePossibilities(grid, iteration, row, column, val) {
		return true
	}
	for i := 1; i <= max; i++ {
		if i != val {
			if checkIfUniqueAndEliminate(grid, iteration, row, column, i) {
				return true
			}
		}
	}
	return false
}

// Called when a number is set in a cell.
// Eliminates the possibilities from the grid given a number at a row,col
// TODO eliminate only for row/col/block required
func eliminatePossibilities(grid *datatypes.Grid, iteration int, row int, column int, val int) (backtrack bool) {
	for i := 0; i < max; i++ {
		if i == row {
			continue
		}
		backtrack = eliminatePossibilitiesForPosition(grid, iteration, i, column, val)
		if backtrack {
			return
		}
	}
	for j := 0; j < max; j++ {
		if j == column {
			continue
		}
		backtrack = eliminatePossibilitiesForPosition(grid, iteration, row, j, val)
		if backtrack {
			return
		}
	}
	rowMin, columnMin := getBlockTopLeft(row, column)
	for i := rowMin; i < rowMin+3; i++ {
		for j := columnMin; j < columnMin+3; j++ {
			if i == row && j == column {
				continue
			}
			backtrack = eliminatePossibilitiesForPosition(grid, iteration, i, j, val)
			if backtrack {
				return
			}
		}
	}
	return
}

// Called when a number is set in a cell.
// Eliminates the possibilities from the grid given a number at a i,j
func eliminatePossibilitiesForPosition(grid *datatypes.Grid, iteration int, i int, j int, val int) bool {
	cell := &grid[i][j]
	cell.Mutex.Lock()
	setValue, updated, backtrack := updateCell(cell, iteration, val)
	cell.Mutex.Unlock()
	if backtrack {
		return true
	}
	if setValue > 0 {
		if eliminatePossibilities(grid, iteration, i, j, setValue) {
			return true
		}
	}
	if updated {
		if checkIfUniqueAndEliminate(grid, iteration, i, j, val) {
			return true
		}
	}
	return false
}

// return valuesEliminated,isValueSet
func setValueForCell(cell *datatypes.Cell, iteration int, setValue int) (eliminatedValues []int, isValueSet bool) {
	existingValue := cell.IterationValues[iteration]
	if (existingValue.Possible[setValue] && *cell.Val == 0) || *cell.Val == setValue {
		*existingValue.Val = setValue
		*cell.Val = setValue
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
	if *cell.Val == valToDelete {
		return 0, false, true
	}
	updated := false
	setValue := 0
	if *cell.Val == 0 && existingValue.Possible[valToDelete] {
		updated = true
		delete(existingValue.Possible, valToDelete)
		if len(existingValue.Possible) == 1 {
			for key := range existingValue.Possible {
				setValue = key
				*existingValue.Val = key
				*cell.Val = key
			}
		}
	}
	return setValue, updated, false
}

func checkIfUniqueAndEliminate(grid *datatypes.Grid, iteration int, i int, j int, val int) bool {
	uniquePositions, conflict := checkIfUnique(grid, iteration, val, datatypes.Position{X: i, Y: j})
	if conflict {
		return true
	}
	for _, pos := range uniquePositions {
		setCell := &grid[pos.X][pos.Y]
		setCell.Mutex.Lock()
		eliminatedValues, isValueSet := setValueForCell(setCell, iteration, val)
		setCell.Mutex.Unlock()
		if isValueSet {
			for _, eliminatedVal := range eliminatedValues {
				if checkIfUniqueAndEliminate(grid, iteration, pos.X, pos.Y, eliminatedVal) {
					return true
				}
			}
			if eliminatePossibilities(grid, iteration, pos.X, pos.Y, val) {
				return true
			}
		} else {
			return true
		}
	}
	return false
}

// TODO add support for row/col/block
// Called when a cell is updated
// If the value deleted now exists once in the row/block/col, then the cell is returned
// returns uniquePositions array, backtrack
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
		cell := grid[row][j]
		if *cell.Val == valDeleted {
			return pos, false, true
		}
		if val.Possible[valDeleted] {
			if found {
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
		cell := grid[i][column]
		if *cell.Val == valDeleted {
			return pos, false, true
		}
		if val.Possible[valDeleted] {
			if found {
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
			cell := grid[i][j]
			if *cell.Val == valDeleted {
				return pos, false, true
			}
			if val.Possible[valDeleted] {
				if found {
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

func remainingPositions(grid *datatypes.Grid, positions map[datatypes.Position]bool) map[datatypes.Position]bool {
	emptyPositions := make(map[datatypes.Position]bool)
	for pos := range positions {
		if *grid[pos.X][pos.Y].Val == 0 {
			emptyPositions[pos] = true
		}
	}
	return emptyPositions
}

func initPositions(grid *datatypes.Grid) map[datatypes.Position]bool {
	positions := make(map[datatypes.Position]bool)
	index := 0
	for i := 0; i < max; i++ {
		for j := 0; j < max; j++ {
			positions[datatypes.Position{X: i, Y: j}] = true
			index++
		}
	}
	return remainingPositions(grid, positions)
}

func solveByGuessing(grid *datatypes.Grid, positions map[datatypes.Position]bool, iteration int) {
	// grid copy
	if len(positions) == 0 {
		numSolutions++
		grid.Print()
		return
	}
	pos := copyValuesForNextIteration(grid, positions, iteration)
	existingValue := grid[pos.X][pos.Y].IterationValues[iteration]
	for val := range existingValue.Possible {
		nextValue := grid[pos.X][pos.Y].IterationValues[iteration+1]
		*grid[pos.X][pos.Y].Val = val
		*nextValue.Val = val
		for key := range nextValue.Possible {
			if key != val {
				delete(nextValue.Possible, key)
			}
		}
		if eliminateUsingGivenValues(grid, iteration+1, pos.X, pos.Y, val) {
		} else {
			updatedPositions := remainingPositions(grid, positions)
			solveByGuessing(grid, updatedPositions, iteration+1)
		}
		copyValuesForNextIteration(grid, positions, iteration)
	}
	return
}

func copyValuesForNextIteration(grid *datatypes.Grid, positions map[datatypes.Position]bool, iteration int) (minPos datatypes.Position) {
	minPossibilities := max + 1
	for pos := range positions {
		cell := grid[pos.X][pos.Y]
		cell.IterationValues[iteration+1] = *datatypes.CopyValue(cell.IterationValues[iteration])
		*cell.Val = 0
		countPossibilities := len(cell.IterationValues[iteration].Possible)
		if countPossibilities < minPossibilities {
			minPossibilities = countPossibilities
			minPos = pos
		}
	}
	return
}

func handleError(msg string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	} else {
		fmt.Fprintf(os.Stderr, "error: %v\n", msg)
	}
	os.Exit(1)
}
