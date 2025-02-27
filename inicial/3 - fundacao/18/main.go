package main

import (
	"fmt"
	"pos/matematica"

	"github.com/google/uuid"
)

func main() {
	s := matematica.Soma(10, 20)
	fmt.Println("Resultado: ", s)
	carro := matematica.Carro{Marca: "Ford"}
	fmt.Println(carro.Andar())

	fmt.Println(carro)
	fmt.Println(matematica.A)
	fmt.Println(uuid.New())
}
