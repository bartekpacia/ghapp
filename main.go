package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", handleWebhook)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalln("Error listening: ", err)
	}
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
