package main

import "fmt"

type Endereco struct {
	Logradouro string
	Numero     int
	Cidade     string
	Estado     string
}

type Pessoa interface {
	Desativar()
}

type Client struct {
	Nome    string
	Idade   int
	Ativo   bool
	Adddres Endereco
}

func (cliente Client) Desativar() {
	cliente.Ativo = false
	fmt.Printf("O Client %s foi desativado", cliente.Nome)

}

type Empresa struct {
	Nome string
}

func (e Empresa) Desativar() {

}

func Destivacao(pessoa Pessoa) {
	pessoa.Desativar()
}

func main() {

	david := Client{
		Nome:  "David",
		Idade: 42,
		Ativo: true,
	}
	minhaEmpresa := Empresa{}
	fmt.Println(david)
	Destivacao(minhaEmpresa)

	// fmt.Printf("Nome %s, Idade: %d, Ativo %t", david.Nome, david.Idade, david.Ativo)
	// david.Adddres.Cidade = "Minas Gerias"
	// david.Adddres.Logradouro = "Walter Diniz"
}
