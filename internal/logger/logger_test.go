package logger

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInitialize(t *testing.T) {
	tests := []struct {
		name    string
		level   string
		wantErr bool
	}{{name: "RightlevelString",
		level:   "Info",
		wantErr: false,
	},
		{name: "WronglevelString",
			level:   "WrongLevel",
			wantErr: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Initialize(tt.level); (err != nil) != tt.wantErr {
				t.Errorf("Initialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRequestResponseLogger(t *testing.T) {
	type args struct {
		method string
		uri    string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "firstTest",
			args: args{method: "GET",
				uri: ""},
			want: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {

		//logger := zaptest.NewLogger(t)

		t.Run(tt.name, func(t *testing.T) {

			r, _ := http.NewRequest(tt.args.method, tt.args.uri, nil)
			w := httptest.NewRecorder()
			//homeHandle.ServeHTTP(w, r)

			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				//fmt.Println("i'am in handl")
			})
			h.ServeHTTP(w, r)

			if got := RequestResponseLogger(h); !true {
				fmt.Printf("%#v", got)
				t.Errorf("RequestResponseLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loggingResponseWriter_Write(t *testing.T) {
	type fields struct {
		ResponseWriter http.ResponseWriter
		responseData   *responseData
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &loggingResponseWriter{
				ResponseWriter: tt.fields.ResponseWriter,
				responseData:   tt.fields.responseData,
			}
			got, err := r.Write(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Write() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loggingResponseWriter_WriteHeader(t *testing.T) {
	type fields struct {
		ResponseWriter http.ResponseWriter
		responseData   *responseData
	}
	type args struct {
		statusCode int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &loggingResponseWriter{
				ResponseWriter: tt.fields.ResponseWriter,
				responseData:   tt.fields.responseData,
			}
			r.WriteHeader(tt.args.statusCode)
		})
	}
}
