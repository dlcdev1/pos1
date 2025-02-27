package main

func main() {

	//Memoria -> endereço -> valor
	a := 10
	var pointeiro *int = &a
	*pointeiro = 20
	b := &a
	*b = 30
	println(*b)
}
