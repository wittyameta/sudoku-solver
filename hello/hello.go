package main

import (
	"fmt"
)

func main() {
	var yo,zo string
	fmt.Scanf("%s%s",&yo,&zo)
	fmt.Printf("Hello, world.\n")
	fmt.Println(yo,zo)
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
