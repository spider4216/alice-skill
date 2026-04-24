package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/spider4216/alice-skill/internal/logger"
	"github.com/spider4216/alice-skill/internal/models"
	"github.com/spider4216/alice-skill/internal/store"
	"go.uber.org/zap"
)

type app struct {
	store store.MessageStore
}

func newApp(s store.MessageStore) *app {
	return &app{
		store: s,
	}
}

func (a *app) webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method %s not allowed", r.Method)
		return
	}

	ctx := r.Context()

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

	messages, err := a.store.ListMessages(ctx, req.Session.User.UserID)

	if err != nil {
		logger.Log.Error("cannot load messages for user", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	text := "Для вас нет новых сообщений"

	if len(messages) > 0 {
		text = fmt.Sprintf("Для вас %d новых сообщений", len(messages))
	}

	if req.Session.New {
		tz, err := time.LoadLocation(req.Timezone)

		if err != nil {
			logger.Log.Error("cannot parse timezone")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		now := time.Now().In(tz)
		hour, min, _ := now.Clock()

		text = fmt.Sprintf("Точное время %d часов, %d минут. %s", hour, min, text)
	}

	resp := models.Response{
		Response: models.ResponsePayload{
			Text: text,
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
