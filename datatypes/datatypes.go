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

type Cell struct {
	Pos             Position
	IterationValues map[int]Value
	Mutex           sync.Mutex
}

func NewCell(x int, y int) *Cell {
	pos := Position{x, y}
	value := *InitValue()
	iterationValues := make(map[int]Value)
	iterationValues[0] = value
	cell := Cell{pos, iterationValues, sync.Mutex{}}
	return &cell
}

func NewCellWithValue(x int, y int, val Value) *Cell {
	pos := Position{x, y}
	iterationValues := make(map[int]Value)
	iterationValues[0] = val
	cell := Cell{pos, iterationValues, sync.Mutex{}}
	return &cell
}

type Backtrack struct {
	CurrentPos         Position
	CurrentVal         Value
	BacktrackPositions map[Position]bool
}

func NewBacktrack(x int, y int, val Value) *Backtrack {
	currentPos := Position{x, y}
	backTrackPositions := make(map[Position]bool)
	backtrack := Backtrack{currentPos, val, backTrackPositions}
	return &backtrack
}

type Stack []Backtrack

func (s *Stack) Push(v Backtrack) {
	*s = append(*s, v)
}

func (s *Stack) Pop() Backtrack {
	ret := (*s)[len(*s)-1]
	*s = (*s)[0 : len(*s)-1]
	return ret
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

func (grid *Grid) Print(iter int) {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			fmt.Printf("%d ", *grid[i][j].IterationValues[iter].Val)
		}
		fmt.Println()
	}
	fmt.Println()
}

func (grid *Grid) PrintAll(iter int) {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			possibilities := grid[i][j].IterationValues[iter].Possible
			for k := range possibilities {
				fmt.Printf("%d,", k)
			}
			fmt.Printf(" ")
		}
		fmt.Println()
	}
	fmt.Println()
}
