package main

import (
	"fmt"
	"net/http"
)

var number uint64 = 0

func main() {

	//m := sync.Mutex{}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//m.Lock()
		number++
		//m.Unlock()
		w.Write([]byte(fmt.Sprintf("Você é o visitante numero %d", number)))
	})

	http.ListenAndServe(":3000", nil)
}

// go run -race main.go
