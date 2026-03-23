package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
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
	successBody := Response{
		Response: ResponseBody{
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
		method       string
		expectedCode int
		expectedBody string
	}{
		{
			method:       http.MethodGet,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
		},
		{
			method:       http.MethodPut,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
		},
		{
			method:       http.MethodDelete,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
		},
		{
			method:       http.MethodPost,
			expectedCode: http.StatusOK,
			expectedBody: string(successJson),
		},
	}

	for _, tc := range testCases {

		req := resty.New().R()
		req.Method = tc.method
		req.URL = srv.URL

		resp, err := req.Send()
		require.NoError(t, err)

		assert.Equal(t, tc.expectedCode, resp.StatusCode())

		if tc.expectedBody != "" {
			assert.JSONEq(t, tc.expectedBody, string(resp.Body()))
		}
	}

}
