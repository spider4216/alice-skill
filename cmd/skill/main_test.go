package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/spider4216/alice-skill/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type successBody struct {
	Version  string `json:"version"`
	Response response
}

type response struct {
	Text string `json:"text"`
}

func TestWebhook(t *testing.T) {
	successBody := models.Response{
		Response: models.ResponsePayload{
			Text: "Извините, я пока ничего не умею",
		},
		Version: ApiVer,
	}

	successJson, err := json.Marshal(successBody)

	assert.NoError(t, err)

	handler := http.HandlerFunc(webhook)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	testCases := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "method_get",
			method:       http.MethodGet,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
		},
		{
			name:         "methid_put",
			method:       http.MethodPut,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
		},
		{
			name:         "method_delete",
			method:       http.MethodDelete,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
		},
		{
			name:         "method_post_without_body",
			method:       http.MethodPost,
			expectedCode: http.StatusInternalServerError,
			expectedBody: "",
		},
		{
			name:         "method_post_unsupported_type",
			method:       http.MethodPost,
			body:         `{"request": {"type": "idunno", "command": "do something"}, "version": "1.0"}`,
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: "",
		},
		{
			name:         "method_post_success",
			method:       http.MethodPost,
			body:         `{"request": {"type": "SimpleUtterance", "command": "sudo do something"}, "version": "1.0"}`,
			expectedCode: http.StatusOK,
			expectedBody: string(successJson),
		},
	}

	for _, tc := range testCases {

		req := resty.New().R()
		req.Method = tc.method
		req.URL = srv.URL

		if len(tc.body) > 0 {
			req.SetHeader("Content-Type", "application/json")
			req.SetBody(tc.body)
		}

		resp, err := req.Send()
		require.NoError(t, err)

		assert.Equal(t, tc.expectedCode, resp.StatusCode())

		if tc.expectedBody != "" {
			assert.JSONEq(t, tc.expectedBody, string(resp.Body()))
		}
	}

}
