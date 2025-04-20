package handlers

import (
	"fmt"
	"github.com/forester2k/metrics/internal/storage"
	"net/http"
	"strconv"
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

func ReadStoredHandler(w http.ResponseWriter, r *http.Request) {
	mType := r.PathValue("mType")
	mName := r.PathValue("mName")
	m, err := URLValidate(mType, mName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	num, err := storage.Store.Get(&m)
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
