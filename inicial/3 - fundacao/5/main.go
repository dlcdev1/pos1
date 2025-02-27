package main

import (
	"fmt"
)

const a = "Hello, World!"

type ID int

var (
	b bool    = true
	c int     = 10
	d string  = "David"
	e float64 = 1.2
	f ID      = 1
)

func main() {
	var myArray [3]int
	myArray[0] = 10
	myArray[1] = 20
	myArray[2] = 30
	fmt.Println(myArray[len(myArray)-1])

	for i, v := range myArray {
		fmt.Printf("O valor de %d é %d\n", i, v)
	}
}
