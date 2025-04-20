package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

func RequestGzipDecompressor(h http.Handler) http.Handler {
	gzipDecompressFn := func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gz.Close()
		r.Body = io.NopCloser(gz)
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(gzipDecompressFn)
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func ResponseGzipCompressor(h http.Handler) http.Handler {
	gzipCompressFn := func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gz.Close()
		w.Header().Set("Content-Encoding", "gzip")
		h.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	}
	return http.HandlerFunc(gzipCompressFn)
}
