package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	log.Println("server is starting")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", index)
	mux.HandleFunc("POST /webhook", handleWebhook)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalln("Error listening: ", err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	// Print request information, such as headers and body
	fmt.Println("Request Headers: ", r.Header)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading body: ", err)
		return
	} else {
		fmt.Println("Request Body: ", string(body))
	}
}
