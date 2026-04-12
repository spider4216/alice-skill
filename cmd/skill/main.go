package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spider4216/alice-skill/internal/logger"
	"github.com/spider4216/alice-skill/internal/models"
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

	logger.Log.Debug("decoding request")

	req := models.Request{}
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&req); err != nil {
		logger.Log.Debug("cannot decod json request body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if req.Request.Type != models.TypeSimpleUtterance {
		logger.Log.Debug("unsupported request type", zap.String("type", req.Request.Type))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	resp := models.Response{
		Response: models.ResponsePayload{
			Text: "Извините, я пока ничего не умею",
		},
		Version: ApiVer,
	}

	w.Header().Set("content-type", "application/json")

	enc := json.NewEncoder(w)

	if err := enc.Encode(resp); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		return
	}

	logger.Log.Debug("sending HTTP 200 response")
}
