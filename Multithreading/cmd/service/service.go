package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Address struct {
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Cidade     string `json:"localidade"`
	Uf         string `json:"uf"`
}

func FindBrasilAPI(cep string, wg *sync.WaitGroup, resultChan chan<- string) {
	defer wg.Done()
	client := &http.Client{Timeout: 1 * time.Second}
	resp, err := client.Get(fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep))
	if err != nil {
		resultChan <- fmt.Sprintf("BrasilAPI: erro - %s", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		resultChan <- "BrasilAPI: erro - resposta não foi 200"
		return
	}

	var address Address
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		resultChan <- fmt.Sprintf("BrasilAPI: erro - %s", err.Error())
		return
	}
	resultChan <- fmt.Sprintf("BrasilAPI: %s, %s, %s - %s", address.Logradouro, address.Bairro, address.Cidade, address.Uf)
}

func FindViaCEP(cep string, wg *sync.WaitGroup, resultChan chan<- string) {
	defer wg.Done()
	client := &http.Client{Timeout: 1 * time.Second}
	resp, err := client.Get(fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep))
	if err != nil {
		resultChan <- fmt.Sprintf("ViaCEP: erro - %s", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		resultChan <- "ViaCEP: erro - resposta não foi 200"
		return
	}

	var address Address
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		resultChan <- fmt.Sprintf("ViaCEP: erro - %s", err.Error())
		return
	}
	resultChan <- fmt.Sprintf("ViaCEP: %s, %s, %s - %s", address.Logradouro, address.Bairro, address.Cidade, address.Uf)
}
