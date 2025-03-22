package handlers

//func TestWebhook(t *testing.T) {
//	storage.Store = &storage.MemStorage{}
//	storage.Store.Init()
//
//	type want struct {
//		contentType string
//		statusCode  int
//	}
//	tests := []struct {
//		name    string
//		method  string
//		request string
//		mType   string
//		mName   string
//		mValue  string
//		want    want
//	}{
//		{
//			name:    "Good gauge metriccccc",
//			method:  http.MethodPost,
//			request: "/update/gauge/Alloc/1.1",
//			mType:   "gauge",
//			mName:   "Alloc",
//			mValue:  "1.1",
//			want:    want{contentType: "text/plain", statusCode: http.StatusOK},
//		},
//		{
//			name:    "Good counter metric",
//			method:  http.MethodPost,
//			request: "/update/counter/PollCount/3",
//			mType:   "counter",
//			mName:   "PollCount",
//			mValue:  "3",
//			want:    want{contentType: "text/plain", statusCode: http.StatusOK},
//		},
//		{
//			name:    "Bad gauge metric value",
//			method:  http.MethodPost,
//			request: "/update/gauge/metricName/1a",
//			mType:   "gauge",
//			mName:   "metricName",
//			mValue:  "1a",
//			want:    want{contentType: "", statusCode: http.StatusBadRequest},
//		},
//		{
//			name:    "Bad counter metric value",
//			method:  http.MethodPost,
//			request: "/update/counter/metricName/3.1",
//			mType:   "counter",
//			mName:   "metricName",
//			mValue:  "3.1",
//			want:    want{contentType: "", statusCode: http.StatusBadRequest},
//		},
//		{
//			name:    "Bad metric type",
//			method:  http.MethodPost,
//			request: "/update/somethingWrong/metricName/1.1",
//			mType:   "somethingWrong",
//			mName:   "metricName",
//			mValue:  "1.1",
//			want:    want{contentType: "", statusCode: http.StatusBadRequest},
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			request := httptest.NewRequest(tt.method, tt.request, nil)
//			response := httptest.NewRecorder()
//			request.SetPathValue("mType", tt.mType)
//			request.SetPathValue("mName", tt.mName)
//			request.SetPathValue("mValue", tt.mValue)
//			Webhook(response, request)
//			result := response.Result()
//			assert.Equal(t, tt.want.statusCode, result.StatusCode)
//			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
//			result.Body.Close()
//		})
//	}
//}

//func TestWriteJSONMetricHandler(t *testing.T) {
//	handler := http.HandlerFunc(WriteJSONMetricHandler)
//	srv := httptest.NewServer(handler)
//	defer srv.Close()
//
//	testCases := []struct {
//		name         string // добавляем название тестов
//		method       string
//		body         string // добавляем тело запроса в табличные тесты
//		expectedCode int
//		expectedBody string
//	}{
//		{
//			name:         "request_without_body",
//			method:       http.MethodPost,
//			expectedCode: http.StatusBadRequest,
//			expectedBody: "",
//		},
//		{
//			name:         "request_unsupported_type",
//			method:       http.MethodPost,
//			body:         `{"id": "testSetGet103","type": "Counter","delta": 333}`,
//			expectedCode: http.StatusBadRequest,
//			expectedBody: "",
//		},
//		{
//			name:         "request_counter_without_delta",
//			method:       http.MethodPost,
//			body:         `{"id": "testSetGet103","type": "counter"}`,
//			expectedCode: http.StatusBadRequest,
//			expectedBody: "",
//		},
//		{
//			name:         "request_counter_with_wrong_delta",
//			method:       http.MethodPost,
//			body:         `{"id": "testSetGet103","type": "counter","delta": 333.33}`,
//			expectedCode: http.StatusBadRequest,
//			expectedBody: "",
//		},
//		{
//			name:         "request_gauge_with_wrong_value",
//			method:       http.MethodPost,
//			body:         `{"id": "testSetGet103","type": "gauge","value": "333.33"}`,
//			expectedCode: http.StatusBadRequest,
//			expectedBody: "",
//		},
//		{
//			name:         "request_#1_gauge_with_OK_value",
//			method:       http.MethodPost,
//			body:         `{"id": "testSetGet103","type": "gauge","value": 333.33}`,
//			expectedCode: http.StatusOK,
//			expectedBody: `{"id": "testSetGet103","type": "gauge","value": 333.33}`,
//		},
//		{
//			name:         "request_#2_gauge_with_OK_value",
//			method:       http.MethodPost,
//			body:         `{"id": "testSetGet103","type": "gauge","value": 333.33}`,
//			expectedCode: http.StatusOK,
//			expectedBody: `{"id": "testSetGet103","type": "gauge","value": 333.33}`,
//		},
//		{
//			name:         "request_#1_counter_with_OK_delta",
//			method:       http.MethodPost,
//			body:         `{"id": "testSetGet104","type": "counter","delta": 333}`,
//			expectedCode: http.StatusOK,
//			expectedBody: `{"id": "testSetGet104","type": "counter","delta": 333}`,
//		},
//		{
//			name:         "request_#2_counter_with_OK_delta",
//			method:       http.MethodPost,
//			body:         `{"id": "testSetGet104","type": "counter","delta": 333}`,
//			expectedCode: http.StatusOK,
//			expectedBody: `{"id": "testSetGet104","type": "counter","delta": 666}`,
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			req := resty.New().R()
//			req.Method = tc.method
//			req.URL = srv.URL
//			if len(tc.body) > 0 {
//				req.SetHeader("Content-Type", "application/json")
//				req.SetBody(tc.body)
//			}
//			resp, err := req.Send()
//			assert.NoError(t, err, "error making HTTP request")
//			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
//			if tc.expectedBody != "" {
//				assert.JSONEq(t, tc.expectedBody, string(resp.Body()), "Response body didn't match expected")
//			} else {
//				assert.Equal(t, int64(0), resp.Size(), "Response should not have a body")
//			}
//		})
//	}
//}
