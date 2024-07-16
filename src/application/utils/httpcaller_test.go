package utils

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpCaller(t *testing.T) {
	t.Run("HttpCaller happy path", func(t *testing.T) {
		tests := []struct {
			name         string
			method       string
			path         string
			headers      map[string]string
			requestBody  []byte
			responseBody []byte
			status       int
		}{
			{
				name:   "Normal case",
				method: "GET",
				path:   "/",
				headers: map[string]string{
					"Content-Type": "application/json",
				},
				status:       http.StatusOK,
				requestBody:  []byte(`{"test":"body"}`),
				responseBody: []byte(`{"test":"response"}`),
			},
		}
		for _, tt := range tests {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				w.Write(tt.responseBody)
			}))
			defer server.Close()

			h := &HttpCaller{
				Method:  tt.method,
				Url:     server.URL + tt.path,
				Headers: tt.headers,
				Body:    tt.requestBody,
			}

			resp, err := h.Call()
			if err != nil {
				t.Errorf("Expected no error, got: %s", err)
			}

			respBody, err := io.ReadAll(resp.Body)
			assert.Nil(t, err)
			assert.Equal(t, resp.StatusCode, tt.status)
			assert.Equal(t, respBody, tt.responseBody)
		}
	})

	t.Run("HttpCaller error paths", func(t *testing.T) {
		tests := []struct {
			name   string
			path   string
			method string
			url    string
		}{
			{
				name:   "Incorrect Method",
				method: "*?",
			},
			{
				name:   "Incorrect url",
				method: "GET",
			},
		}
		for _, tt := range tests {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"status":"error"}`))
			}))
			defer server.Close()

			h := &HttpCaller{
				Method: tt.method,
			}

			resp, err := h.Call()

			assert.Error(t, err)
			assert.Nil(t, resp)
		}
	})
}
