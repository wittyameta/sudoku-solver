package datatypes

type Position struct {
	x int
	y int
}

type Value struct {
	val      int
	possible [9]bool
}

func NewValue() *Value {
	possible := [9]bool{true, true, true, true, true, true, true, true, true}
	v := Value{0, possible}
	return &v
}

type Cell struct {
	pos             Position
	iterationValues map[int]Value
}

func NewCell(x int, y int) *Cell {
	pos := Position{x, y}
	value := *NewValue()
	iterationValues := make(map[int]Value)
	iterationValues[0] = value
	cell := Cell{pos, iterationValues}
	return &cell
}

type Backtrack struct {
	currentPos         Position
	currentVal         Value
	backtrackPositions map[Position]bool
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
