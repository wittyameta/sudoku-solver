package datatypes

import (
	"fmt"
	"sync"
)

//Position yo
type Position struct {
	// X is yo
	X int
	// Y is yo
	Y int
}

type Value struct {
	Val      *int
	Possible map[int]bool
}

func InitValue() *Value {
	possible := make(map[int]bool)
	for i := 1; i < 10; i++ {
		possible[i] = true
	}
	val := 0
	v := Value{&val, possible}
	return &v
}

func SetValue(val int) *Value {
	possible := make(map[int]bool)
	possible[val] = true
	v := Value{&val, possible}
	return &v
}

func CopyValue(value Value) *Value {
	possible := make(map[int]bool)
	for key := range value.Possible {
		possible[key] = true
	}
	val := *value.Val
	v := Value{&val, possible}
	return &v
}

type Cell struct {
	Val             *int
	IterationValues map[int]Value
	Mutex           sync.Mutex
}

func NewCell(x int, y int) *Cell {
	value := *InitValue()
	iterationValues := make(map[int]Value)
	iterationValues[0] = value
	val := 0
	cell := Cell{&val, iterationValues, sync.Mutex{}}
	return &cell
}

type Grid [9][9]Cell

func InitGrid() *Grid {
	grid := Grid{}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			grid[i][j] = *NewCell(i, j)
		}
	}
	return &grid
}

func (grid *Grid) Print() {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			fmt.Printf("%d ", *grid[i][j].Val)
		}
		fmt.Println()
	}
	fmt.Println()
}
