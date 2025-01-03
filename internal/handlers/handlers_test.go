package handlers

import (
	"github.com/forester2k/metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebhook(t *testing.T) {
	storage.Store = &storage.MemStorage{}
	storage.Store.Init()

	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name    string
		method  string
		request string
		mType   string
		mName   string
		mValue  string
		want    want
	}{
		{
			name:    "Good gauge metriccccc",
			method:  http.MethodPost,
			request: "/update/gauge/Alloc/1.1",
			mType:   "gauge",
			mName:   "Alloc",
			mValue:  "1.1",
			want:    want{contentType: "text/plain", statusCode: http.StatusOK},
		},
		{
			name:    "Good counter metric",
			method:  http.MethodPost,
			request: "/update/counter/PollCount/3",
			mType:   "counter",
			mName:   "PollCount",
			mValue:  "3",
			want:    want{contentType: "text/plain", statusCode: http.StatusOK},
		},
		{
			name:    "Bad gauge metric value",
			method:  http.MethodPost,
			request: "/update/gauge/metricName/1a",
			mType:   "gauge",
			mName:   "metricName",
			mValue:  "1a",
			want:    want{contentType: "", statusCode: http.StatusBadRequest},
		},
		{
			name:    "Bad counter metric value",
			method:  http.MethodPost,
			request: "/update/counter/metricName/3.1",
			mType:   "counter",
			mName:   "metricName",
			mValue:  "3.1",
			want:    want{contentType: "", statusCode: http.StatusBadRequest},
		},
		{
			name:    "Bad metric type",
			method:  http.MethodPost,
			request: "/update/somethingWrong/metricName/1.1",
			mType:   "somethingWrong",
			mName:   "metricName",
			mValue:  "1.1",
			want:    want{contentType: "", statusCode: http.StatusBadRequest},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, nil)
			response := httptest.NewRecorder()
			request.SetPathValue("mType", tt.mType)
			request.SetPathValue("mName", tt.mName)
			request.SetPathValue("mValue", tt.mValue)
			Webhook(response, request)
			result := response.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			result.Body.Close()
		})
	}
}
