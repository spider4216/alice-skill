package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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
		r := httptest.NewRequest(tc.method, "/", nil)
		w := httptest.NewRecorder()

		webhook(w, r)

		assert.Equal(t, tc.expectedCode, w.Code)

		if tc.expectedBody != "" {
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		}
	}

}
