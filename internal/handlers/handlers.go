package handlers

import (
	"fmt"
	"github.com/forester2k/metrics/internal/service"
	"github.com/forester2k/metrics/internal/storage"
	"net/http"
	"strings"
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
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	urlSlice := strings.Split(strings.Trim(r.URL.Path, "/ "), "/")
	newMetric, err := URLValidate(urlSlice)
	if err != nil {
		w.WriteHeader(getStatusError(err))
		return
	}
	newMetric, err = ValueValidate(urlSlice, newMetric)

	if err != nil {
		w.WriteHeader(getStatusError(err))
		return
	}
	if newMetric == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = newMetric.Save(storage.Store)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "text/plain")

	rr := ""
	_, _ = w.Write([]byte(rr))
}

func ReadStoredHandler(w http.ResponseWriter, r *http.Request) {

	urlSlice := strings.Split(strings.Trim(r.URL.Path, "/ "), "/")
	m, err := URLValidate(urlSlice)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m, err = m.Get(storage.Store)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	var res string
	switch m := m.(type) {
	case service.GaugeMetric:
		res = fmt.Sprint(m.Value)
	case service.CounterMetric:
		res = fmt.Sprint(m.Value)
	}
	_, _ = w.Write([]byte(res))
}

func ListStoredHandler(w http.ResponseWriter, r *http.Request) {
	l := &service.MetricList{Gauges: map[string]float64{}, Counters: map[string]int64{}}
	l.TakeList(storage.Store)

	w.Header().Set("Content-Type", "text/html")

	body := `<html>
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
	_, _ = w.Write([]byte(body))
}
