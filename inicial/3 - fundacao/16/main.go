package main

import "fmt"

func main() {
	var minhaVar1 interface{} = "David Leandro"
	println(minhaVar1.(string))
	res, ok := minhaVar1.(int)

	fmt.Printf("O valor de res é %v e o resultado de ok é %v", res, ok)

	res2 := minhaVar1.(int)
	fmt.Printf("O Valor de res2 é %v", res2)
}