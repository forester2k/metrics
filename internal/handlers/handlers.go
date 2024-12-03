package handlers

import (
	"fmt"
	"net/http"
	"strings"
)

func Webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	urlSlice := strings.Split(strings.Trim(r.URL.Path, "/ "), "/")
	newMetric, err := UrlValidate(urlSlice)
	if err != nil {
		switch err.Error() {
		case "BadRequest":
			w.WriteHeader(http.StatusBadRequest)
		case "NotFound":
			w.WriteHeader(http.StatusNotFound)
		}
		return
	}

	// Отладочный блок
	w.Header().Set("Content-Type", "text/plain")
	errMsg := "Ошибок нет"
	//if err != nil {
	//	errMsg = err.Error()
	//}
	rr := fmt.Sprintf("%#v, %v", newMetric, errMsg)

	_, _ = w.Write([]byte(rr))
}
