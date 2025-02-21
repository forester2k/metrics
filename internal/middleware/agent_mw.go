package middleware

import (
	"bytes"
	"fmt"
	"github.com/forester2k/metrics/internal/compression"
	"io"
	"net/http"
	"net/http/httputil"
	"time"
)

type internalRT func(*http.Request) (*http.Response, error)

func (rt internalRT) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt(req)
}

type Middleware func(http.RoundTripper) http.RoundTripper

func Сonveyer(rt http.RoundTripper, middlewares ...Middleware) http.RoundTripper {
	if rt == nil {
		rt = http.DefaultTransport
	}
	for _, m := range middlewares {
		rt = m(rt)
	}
	return rt
}

// CustomTimer takes a writer and will output a request duration.
func CustomTimer(w io.Writer) Middleware {
	return func(rt http.RoundTripper) http.RoundTripper {
		return internalRT(func(req *http.Request) (*http.Response, error) {
			startTime := time.Now()
			defer func() {
				o, err := httputil.DumpRequest(req, true)
				if err != nil {
					panic(err)
				}
				fmt.Println(string(o))
				fmt.Fprintf(w, ">>> request duration: %s",
					time.Since(startTime))
				fmt.Println(" <<<<<<")
			}()
			return rt.RoundTrip(req)
		})
	}
}

// DumpResponse uses dumps the response body to console.
func DumpResponse(includeBody bool) Middleware {
	return func(rt http.RoundTripper) http.RoundTripper {
		return internalRT(func(req *http.Request) (resp *http.Response, err error) {
			h := req.Header
			for k, v := range h {
				fmt.Println(">>>", k, v)
			}

			defer func() {

				if err == nil {
					h := resp.Header
					for k, v := range h {
						fmt.Println("<<<<", k, v)
					}
					o, err := httputil.DumpResponse(resp, includeBody)
					if err != nil {
						panic(err)
					}
					fmt.Println(string(o))
				}
			}()
			return rt.RoundTrip(req)
		})
	}
}

func Compress(gzipCompressLevel int) Middleware {
	return func(rt http.RoundTripper) http.RoundTripper {
		return internalRT(func(req *http.Request) (resp *http.Response, err error) {
			req.Header.Set("Content-Encoding", "gzip")
			reqBody, err := io.ReadAll(req.Body)
			req.Body.Close()
			gz, err := compression.CompressGzip(reqBody, gzipCompressLevel)
			if err != nil {
				//ToDo: сдклать логгирование
				fmt.Printf("!!! Ошибка при сжатии запроса %w \n", err.Error())
				return rt.RoundTrip(req)
			}
			req.Body = io.NopCloser(bytes.NewBuffer(gz))
			req.ContentLength = int64(len(gz))
			req.Body.Close()
			return rt.RoundTrip(req)
		})
	}
}
