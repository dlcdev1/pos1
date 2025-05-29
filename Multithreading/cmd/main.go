package main

import (
	"fmt"
	"github.com/dlcdev1/pos1/multithreadings/cmd/service"
	"sync"
)

func main() {
	var cep string
	fmt.Print("Digite o CEP: ")
	fmt.Scanln(&cep)

	var wg sync.WaitGroup
	resultChan := make(chan string)

	wg.Add(2)
	go service.FindBrasilAPI(cep, &wg, resultChan)
	go service.FindViaCEP(cep, &wg, resultChan)

	// Criar um goroutine para encerrar o resultado após o WaitGroup
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		fmt.Println(result)
	}
	main()
}
