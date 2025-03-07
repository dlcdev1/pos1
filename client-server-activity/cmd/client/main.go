package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	cotacao, err := getValue(ctx)
	if err != nil {
		log.Fatal("Error get value: ", err)
	}

	if err := saveValue(cotacao); err != nil {
		log.Fatal("Error saving values: ", err)
	}

	fmt.Println("Cotação: ", cotacao)

}

func getValue(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)

	if err != nil {
		return "Error request.", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "Error get value", err
	}
	defer resp.Body.Close()

	var cotacao string

	if err := json.NewDecoder(resp.Body).Decode(&cotacao); err != nil {
		return "Error decoded.", err
	}
	return cotacao, nil
}

func saveValue(cotacao string) error {
	value := fmt.Sprintf("Cotação: %s - salvo em: %s\n", cotacao, time.Now().Format("2006-01-02 15:04:05"))
	file, err := os.OpenFile("cotacao.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(value); err != nil {
		return err
	}

	return nil
}
