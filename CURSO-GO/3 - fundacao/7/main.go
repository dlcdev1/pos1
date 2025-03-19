package main

import "fmt"

func main() {

	salarios := map[string]int{"David": 1000, "João": 2000, "Maria": 3000}
	// fmt.Println(salarios["David"])
	// delete(salarios, "David")
	// salarios["David"] = 5000

	// fmt.Println(salarios["David"])
	// sal1 := make(map[string]int)
	// sal1["David"] = 1000
	// fmt.Println(salarios["David"])
	for nome, salario := range salarios {
		fmt.Println("O salario de %s é %d\n", nome, salario)
	}

	for _, salario := range salarios {
		fmt.Println("O salario de %s é %d\n", salario)
	}

}
