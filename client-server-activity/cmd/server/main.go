package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {
	http.HandleFunc("/cotacao", valueHandler)
	http.ListenAndServe(":8080", nil)
}

func valueHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	values, err := getValue(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctxDB, cancelDB := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancelDB()

	if err := saveValue(ctxDB, values); err != nil {
		log.Println("Error saving values: ", err)
	}

	log.Println("Start request.")
	defer log.Println("Finished request.")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(values)

}

func getValue(ctx context.Context) (string, error) {
	request, err := http.NewRequestWithContext(
		ctx, "GET", "https://economia.awesomeapi.com.br/USD-BRL", nil)
	if err != nil {
		return "Error connected server", err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "Error client", err
	}
	defer response.Body.Close()

	var cotacao []Cotacao

	if err := json.NewDecoder(response.Body).Decode(&cotacao); err != nil {
		return "Error decoded", err
	}

	return cotacao[0].Bid, nil
}

func saveValue(ctx context.Context, cotacao string) error {
	db, err := sql.Open("sqlite3", "values.db")

	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS cotacoes (id INTEGER PRIMARY KEY, bid TEXT, timestamp DATETIME DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, "INSERT INTO cotacoes (bid) VALUES (?)", cotacao)
	log.Println("Value saved successfully!")
	return err
}
