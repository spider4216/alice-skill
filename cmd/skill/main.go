package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	ApiVer = "1.0"
)

func main() {
	run()
}

func run() {
	http.ListenAndServe(":8080", http.HandlerFunc(webhook))
}

func webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method %s not allowed", r.Method)
		return
	}

	w.Header().Set("content-type", "application/json")

	resp := Response{
		Response: ResponseBody{
			Text: "Извините, я пока ничего не умею",
		},
		Version: ApiVer,
	}

	json, err := json.Marshal(resp)

	if err != nil {
		fmt.Fprint(w, "Something went wrong while preparing response")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, string(json))
}

type Response struct {
	Response ResponseBody `json:"response"`
	Version  string       `json:"version"`
}

type ResponseBody struct {
	Text string `json:"text"`
}
