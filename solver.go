package main

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/wittyameta/sudoku-solver/datatypes"
)

// TODO
// min 17 digits and 8 diff digits
// for all loops check if current cell needs to be excluded
// backtrack
// iter count should be max instaed of 0
// function return instead of value
// when num is set, check if after elimination in nearby block, the num is specific to a row/col
const max int = 9

var numSolutions int

func main() {
	numSolutions = 0
	grid := *datatypes.InitGrid()
	count := 0
	for i := 0; i < max; i++ {
		count += readRow(&grid, i)
	}
	solve(&grid, count)
	grid.Print()
	grid.PrintAll(0)
	stack := initStack(&grid)
	fmt.Println("starting stack")
	backtrack(&grid, stack[0].BacktrackPositions, 0)
	fmt.Println(numSolutions)
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
		val := grid[rownum][i].IterationValues
		input := verifyElement(elem)
		if input > 0 {
			if val[0].Possible[input] {
				val[0] = *datatypes.SetValue(input)
				*grid[rownum][i].Val = input
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
				go initialElimination(grid, i, j, val, &wg)
			}
		}
	}
	potentialCountDiff := count - verificationCount
	if potentialCountDiff > 0 {
		wg.Add(potentialCountDiff)
	}
	wg.Wait()
}

func initialElimination(grid *datatypes.Grid, row int, column int, val int, wg *sync.WaitGroup) {
	defer wg.Done()
	b := *datatypes.NewBacktrack(datatypes.Position{X: row, Y: column}, *datatypes.SetValue(val))
	// TODO backtrack necessary?
	if eliminateUsingGivenValues(grid, 0, row, column, val, &b) {
		panic("No solution")
	}
}

func eliminateUsingGivenValues(grid *datatypes.Grid, iteration int, row int, column int, val int, b *datatypes.Backtrack) bool {
	if eliminatePossibilities(grid, iteration, row, column, val, b) {
		fmt.Println("true from A")
		return true
	}
	for i := 1; i <= max; i++ {
		if i != val {
			if checkIfUniqueAndEliminate(grid, iteration, row, column, i, b) {
				fmt.Println("true from B")
				return true
			}
		}
	}
	return false
}

// Called when a number is set in a cell.
// Eliminates the possibilities from the grid given a number at a row,col
// TODO eliminate only for row/col/block required
func eliminatePossibilities(grid *datatypes.Grid, iteration int, row int, column int, val int, b *datatypes.Backtrack) (backtrack bool) {
	for i := 0; i < max; i++ {
		if i == row {
			continue
		}
		backtrack = eliminatePossibilitiesForPosition(grid, iteration, i, column, val, b)
		if backtrack {
			return
		}
	}
	for j := 0; j < max; j++ {
		if j == column {
			continue
		}
		backtrack = eliminatePossibilitiesForPosition(grid, iteration, row, j, val, b)
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
			backtrack = eliminatePossibilitiesForPosition(grid, iteration, i, j, val, b)
			if backtrack {
				return
			}
		}
	}
	return
}

// Called when a number is set in a cell.
// Eliminates the possibilities from the grid given a number at a i,j
func eliminatePossibilitiesForPosition(grid *datatypes.Grid, iteration int, i int, j int, val int, b *datatypes.Backtrack) bool {
	cell := &grid[i][j]
	cell.Mutex.Lock()
	setValue, updated, backtrack := updateCell(cell, iteration, val)
	cell.Mutex.Unlock()
	if backtrack {
		return true
	}
	if setValue > 0 {
		if eliminatePossibilities(grid, iteration, i, j, setValue, b) {
			return true
		}
	}
	if updated {
		b.BacktrackPositions[datatypes.Position{X: i, Y: j}] = true
		if checkIfUniqueAndEliminate(grid, iteration, i, j, val, b) {
			return true
		}
	}
	return false
}

// return valuesEliminated,isValueSet
func setValueForCell(cell *datatypes.Cell, iteration int, setValue int, b *datatypes.Backtrack) (eliminatedValues []int, isValueSet bool) {
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
		if len(eliminatedValues) > 0 {
			b.BacktrackPositions[cell.Pos] = true
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

func checkIfUniqueAndEliminate(grid *datatypes.Grid, iteration int, i int, j int, val int, b *datatypes.Backtrack) bool {
	uniquePositions, conflict := checkIfUnique(grid, iteration, val, datatypes.Position{X: i, Y: j})
	if conflict {
		fmt.Println("from C")
		return true
	}
	for _, pos := range uniquePositions {
		setCell := &grid[pos.X][pos.Y]
		setCell.Mutex.Lock()
		eliminatedValues, isValueSet := setValueForCell(setCell, iteration, val, b)
		setCell.Mutex.Unlock()
		if isValueSet {
			b.BacktrackPositions[pos] = true
			for _, eliminatedVal := range eliminatedValues {
				if checkIfUniqueAndEliminate(grid, iteration, pos.X, pos.Y, eliminatedVal, b) {
					return true
				}
			}
			if eliminatePossibilities(grid, iteration, pos.X, pos.Y, val, b) {
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
		fmt.Println("from D")
		return uniquePositions, true
	}
	uniquePos, foundUnique, atLeastOnce = checkIfUniqueInColumn(grid, iteration, valDeleted, pos)
	if foundUnique {
		uniquePositions = append(uniquePositions, uniquePos)
	} else if !atLeastOnce {
		fmt.Println("from E")
		return uniquePositions, true
	}
	uniquePos, foundUnique, atLeastOnce = checkIfUniqueInBlock(grid, iteration, valDeleted, pos)
	if foundUnique {
		uniquePositions = append(uniquePositions, uniquePos)
	} else if !atLeastOnce {
		fmt.Println("from F")
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
	fmt.Println("from G-start")
	fmt.Println(iteration, valDeleted, pos)
	fmt.Println("from G-end")
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

func initStack(grid *datatypes.Grid) (stack datatypes.Stack) {
	positions := make(map[datatypes.Position]bool)
	index := 0
	for i := 0; i < max; i++ {
		for j := 0; j < max; j++ {
			positions[grid[i][j].Pos] = true
			index++
		}
	}
	backtrack := *datatypes.NewBacktrack(datatypes.Position{X: max, Y: max}, *datatypes.InitValue())
	backtrack.BacktrackPositions = remainingPositions(grid, positions)
	stack.Push(backtrack)
	return
}

func backtrack(grid *datatypes.Grid, positions map[datatypes.Position]bool, iteration int) {
	// grid copy
	copyValuesForNextIteration(grid, positions, iteration)
	for pos := range positions {
		fmt.Println("pos", pos)
		existingValue := grid[pos.X][pos.Y].IterationValues[iteration]
		fmt.Println(existingValue)
		fmt.Println(grid[pos.X][pos.Y].IterationValues[iteration])
		for val := range existingValue.Possible {
			nextValue := grid[pos.X][pos.Y].IterationValues[iteration+1]
			fmt.Println("val", val)
			*grid[pos.X][pos.Y].Val = val        // TODO 0
			setValue := *datatypes.SetValue(val) // TODO 1
			*nextValue.Val = val                 // TODO 2
			for key := range nextValue.Possible {
				if key != val {
					delete(nextValue.Possible, key) // TODO 3
				}
			}
			b := *datatypes.NewBacktrack(pos, setValue) // TODO 4
			if eliminateUsingGivenValues(grid, iteration+1, pos.X, pos.Y, val, &b) {
				grid.Print()
				fmt.Println("backtrack after", pos.X, pos.Y, val)
				// TODO revert todo4,1
				copyValuesForNextIteration(grid, positions, iteration)
				*grid[pos.X][pos.Y].Val = 0
			} else {
				updatedPositions := remainingPositions(grid, positions)
				if len(updatedPositions) == 0 {
					if numSolutions > 0 {
						panic("multiple solutions")
					}
					numSolutions++
					fmt.Println("success at ", iteration)
					grid.Print()
				} else {
					fmt.Println("go next level. remaining cells:", len(updatedPositions))
					grid.Print()
					backtrack(grid, updatedPositions, iteration+1)
				}
				// next val in stack
			}
		}
		return // TODO remove
	}
}

func copyValuesForNextIteration(grid *datatypes.Grid, positions map[datatypes.Position]bool, iteration int) {
	for pos := range positions {
		cell := grid[pos.X][pos.Y]
		cell.IterationValues[iteration+1] = *datatypes.CopyValue(cell.IterationValues[iteration])
		*cell.Val = 0
	}
}

func backtrack1() {
	// increment iter count
	// TODO sorted with length of possible
	// select a pos
	// TODO sort possibilities
	// select a value
	// set the value to pos
	// call eliminateUsingGivenValues
	//
}
