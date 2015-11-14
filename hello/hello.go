package main

import (
	"fmt"

	"github.com/wittyameta/sudoku-solver/stringutil"
)

func main() {
	fmt.Printf("Hello, world.\n")
	fmt.Println(stringutil.Reverse("Hello, world"))
	var x,y string
	fmt.Scanf("%s%s\n",&x,&y)
	fmt.Println("first",x,y)
	fmt.Printf("y is %s",y)
	for a :=0;a<9;a++ {
		for b:=0;b<9;b++ {
			fmt.Scanf("%s%s", &x)
			fmt.Printf("%s ", x)
		}
		fmt.Println("done",a)
	}
}
