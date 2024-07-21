package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
)

func main() {
	log.Println("server is starting")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", index)
	mux.HandleFunc("POST /webhook", handleWebhook)
	mux.HandleFunc("POST /check-runs", handleCheckRuns)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalln("Error listening: ", err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	slog.Info("request received", slog.String("path", r.URL.Path))
	fmt.Fprintln(w, "hello world")
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	l := slog.With("path", r.URL.Path)

	// Print request information, such as headers and body
	l.Info("new request")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		l.Error("error reading body", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error reading body: %v", err)))
		return
	}

	var payload map[string]interface{}
	err = json.Unmarshal([]byte(body), &payload)
	if err != nil {
		l.Error("error reading body", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error reading body: %v", err)
		return
	}

	// Check type of webhook event
	eventType := r.Header.Get("X-GitHub-Event")
	switch eventType {
	case "installation":
		// https://docs.github.com/en/webhooks/webhook-events-and-payloads?actionType=created#installation
		if payload["action"] == "created" {
			repository := payload["repository"].(map[string]interface{})
			fullName := repository["full_name"].(string)
			l.Info("installation created", slog.Any("repo", fullName))
		}
	case "check_suite":
		// https://docs.github.com/en/webhooks/webhook-events-and-payloads?actionType=requested#check_suite
		if payload["action"] == "requested" {
			repository := payload["repository"].(map[string]interface{})
			fullName := repository["full_name"].(string)

			checkSuite := payload["check_suite"].(map[string]interface{})
			headSHA := checkSuite["head_sha"].(string)
			l.Info("check suite requested", slog.Any("repo", fullName), slog.String("commit", headSHA))
		}
	default:
		l.Error("unknown event type", slog.String("type", eventType))
	}
}

func handleCheckRuns(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("Request Headers: ", r.Header)
	// reqBody, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	fmt.Println("Error reading body: ", err)
	// }
	//
	// params := r.URL.Query()
	// owner := params.Get("owner")
	// repo := params.Get("repo")
	// name := params.Get("name")
	//
	// body := []byte(`{
	// 	"name": "linter",
	// 	"head_sha": "1234567890abcdef",
	// }`)
	//
	// http.NewRequest("POST", "https://api.github.com/repos/:owner/:repo/check-runs", body)
}
