package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spider4216/alice-skill/internal/logger"
	"go.uber.org/zap"
)

const (
	ApiVer = "1.0"
)

func main() {
	parseFlags()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	if err := logger.Initialize(flagLogLevel); err != nil {
		return err
	}

	logger.Log.Info("Running server", zap.String("address", flagRunAddr))

	fmt.Println("Run on", flagRunAddr)
	return http.ListenAndServe(flagRunAddr, http.HandlerFunc(logger.RequestLogger(webhook)))
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

	logger.Log.Debug("sending HTTP 200 response")
}

type Response struct {
	Response ResponseBody `json:"response"`
	Version  string       `json:"version"`
}

type ResponseBody struct {
	Text string `json:"text"`
}
