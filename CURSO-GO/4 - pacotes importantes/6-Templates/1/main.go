package main

import (
	"os"
	"text/template"
)

type Curso struct {
	Nome         string
	CargaHoraria int
}

func main() {
	curso := Curso{"Go", 40}
	tmp := template.New("CursoTemaplate")
	tmp, err := tmp.Parse("Cruso: {{.Nome}} - Carga Horária: {{.CargaHoraria}}")
	if err != nil {
		panic(err)
	}
	err = tmp.Execute(os.Stdout, curso)
	if err != nil {
		panic(err)
	}

}
