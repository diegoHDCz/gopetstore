package main

import (
	"encoding/json"
	"net/http"
)

type Message struct {
	Greeting string `json:"greeting"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", sayHello)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	server.ListenAndServe()

}

func sayHello(w http.ResponseWriter, r *http.Request) {
	message := Message{
		Greeting: "say hello!",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}
