package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebhook(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name    string
		method  string
		request string
		want    want
	}{
		{
			name:    "Good gauge metric",
			method:  http.MethodPost,
			request: "/update/gauge/metricName/1.1",
			want:    want{contentType: "text/plain", statusCode: http.StatusOK},
		},
		{
			name:    "Good counter metric",
			method:  http.MethodPost,
			request: "/update/counter/metricName/3",
			want:    want{contentType: "text/plain", statusCode: http.StatusOK},
		},
		{
			name:    "Bad gauge metric",
			method:  http.MethodPost,
			request: "/update/gauge/metricName/1a",
			want:    want{contentType: "", statusCode: http.StatusBadRequest},
		},
		{
			name:    "Bad counter metric",
			method:  http.MethodPost,
			request: "/update/counter/metricName/3.1",
			want:    want{contentType: "", statusCode: http.StatusBadRequest},
		},
		{
			name:    "Bad metric type",
			method:  http.MethodPost,
			request: "/update/somethingWrong/metricName/1.1",
			want:    want{contentType: "", statusCode: http.StatusBadRequest},
		},
		{
			name:    "Not update in path",
			method:  http.MethodPost,
			request: "/updater/gauge/metricName/1.1",
			want:    want{contentType: "", statusCode: http.StatusBadRequest},
		},
		{
			name:    "Short path",
			method:  http.MethodPost,
			request: "/update/gauge/",
			want:    want{contentType: "", statusCode: http.StatusNotFound},
		},
		{
			name:    "Wrong http-method (Get)",
			method:  http.MethodGet,
			request: "/update/gauge/metricName/1.1",
			want:    want{contentType: "", statusCode: http.StatusMethodNotAllowed},
		},
		{
			name:    "Wrong http-method (Delete)",
			method:  http.MethodDelete,
			request: "/update/gauge/metricName/1.1",
			want:    want{contentType: "", statusCode: http.StatusMethodNotAllowed},
		},
		{
			name:    "Wrong http-method (Patch)",
			method:  http.MethodPatch,
			request: "/update/gauge/metricName/1.1",
			want:    want{contentType: "", statusCode: http.StatusMethodNotAllowed},
		},
		{
			name:    "Wrong http-method (Put)",
			method:  http.MethodPut,
			request: "/update/gauge/metricName/1.1",
			want:    want{contentType: "", statusCode: http.StatusMethodNotAllowed},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, nil)
			response := httptest.NewRecorder()
			Webhook(response, request)
			result := response.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			result.Body.Close()
		})
	}
}
