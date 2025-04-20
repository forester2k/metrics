package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/forester2k/metrics/internal/logger"
	"github.com/forester2k/metrics/internal/service"
	"github.com/forester2k/metrics/internal/storage"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
	"time"
)

func getStatusError(err error) int {
	errorMsg := err.Error()
	switch errorMsg {
	case "BadRequest":
		return http.StatusBadRequest
	case "NotFound":
		return http.StatusNotFound
	default:
		return http.StatusNotImplemented
	}
}

func Webhook(w http.ResponseWriter, r *http.Request) {
	mType := r.PathValue("mType")
	mName := r.PathValue("mName")
	mValue := r.PathValue("mValue")
	m, err := URLValidate(mType, mName)
	if err != nil {
		w.WriteHeader(getStatusError(err))
		return
	}
	m, err = ValueValidate(mValue, m)
	if err != nil {
		w.WriteHeader(getStatusError(err))
		return
	}
	if m == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//err = newMetric.Send(*storage.Store)
	err = storage.Store.Save(&m)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	w.Header().Set("Content-Type", "text/plain")
	rr := ""
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(rr))
}

func WriteJSONMetricHandler(w http.ResponseWriter, r *http.Request) {
	// десериализуем запрос в структуру модели
	logger.Log.Debug("decoding request")
	var req service.Metrics
	//if r.Body == io.ReadCloser(") {
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m, err := JSONTypeAndValueValidate(&req)
	if err != nil {
		w.WriteHeader(getStatusError(err))
		return
	}
	err = storage.Store.Save(&m)
	if err != nil {
		w.WriteHeader(http.StatusNotFound) // Не правильнее ли InternalServerError здесь?
	}
	freshM, err := storage.Store.GetMetric(&m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp := service.ConvToMetrics(freshM)
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		return
	}
	w.WriteHeader(http.StatusOK)
	// _, _ = w.Write([]byte(""))
}

func ReadJSONMetricHandler(w http.ResponseWriter, r *http.Request) {
	// десериализуем запрос в структуру модели
	logger.Log.Debug("decoding request")
	var req service.Metrics
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError) // Не правильнее ли BadRequest здесь?
		return
	}
	m, err := JSONTypeValidate(&req)
	if err != nil {
		w.WriteHeader(getStatusError(err))
		return
	}
	// Возможно надо убрать ! ! !
	if m == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	freshM, err := storage.Store.GetMetric(&m)
	if err != nil {
		w.WriteHeader(http.StatusNotFound) // Не правильнее ли InternalServerError здесь?
		return
	}
	resp := service.ConvToMetrics(freshM)
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		logger.Log.Debug("error encoding response", zap.Error(err))
		return
	}
	w.WriteHeader(http.StatusOK)
	// _, _ = w.Write([]byte(""))
}

func ReadStoredHandler(w http.ResponseWriter, r *http.Request) {
	mType := r.PathValue("mType")
	mName := r.PathValue("mName")
	m, err := URLValidate(mType, mName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	num, err := storage.Store.GetValue(&m)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/plain")

	var res string
	switch n := num.(type) {
	case float64:
		res = strconv.FormatFloat(n, 'g', -1, 64)
	case int64:
		res = strconv.FormatInt(n, 10)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(res))
}

func ListStoredHandler(w http.ResponseWriter, r *http.Request) {
	l := storage.Store.MakeList()
	w.Header().Set("Content-Type", "text/html")
	body := `<!doctype html>
		<html>
    	<head>
    	<title></title>
    	</head>
    	<body>
		`
	end := `     
		</body>
		</html>`
	for k, v := range l.Gauges {
		body += fmt.Sprintf("<br> %s   %f\n", k, v)
	}
	for k, v := range l.Counters {
		body += fmt.Sprintf("<br> %s   %d\n", k, v)
	}
	body += end
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(body))
}

// ToDo: наполнить функцию ! сделано!
func CheckDBConnection(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	status := http.StatusOK
	if err := storage.DB.PingContext(ctx); err != nil {
		status = http.StatusInternalServerError
	}
	w.WriteHeader(status)
	// _, _ = w.Write([]byte(""))
}
