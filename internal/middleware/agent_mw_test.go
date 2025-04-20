package middleware

//
//import (
//	"bytes"
//	"compress/gzip"
//	"fmt"
//	"io"
//	"log"
//	"net/http"
//	"reflect"
//	"testing"
//)
//
//func TestCompress(t *testing.T) {
//	type args struct {
//		CompressLevel int
//	}
//	tests := []struct {
//		name       string
//		args       args
//		bodyString string
//		wantErr    bool
//	}{
//		{
//			name: "rightCompressLevel",
//			args: args{
//				CompressLevel: 1,
//			},
//			bodyString: "Hello, my little calm world",
//			wantErr:    false,
//		},
//		//{
//		//	name: "wrongCompressLevel",
//		//	args: args{
//		//		CompressLevel: 10,
//		//	},
//		//	bodyString: "Hello, my little calm world",
//		//	wantErr:    true,
//		//},
//		// TODO: Add test cases.
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			b := bytes.NewBuffer([]byte(tt.bodyString))
//			req, err := http.NewRequest("POST", "localhost:8080", io.Reader(b))
//			if err != nil {
//				panic(err)
//			}
//			mw := Compress(tt.args.CompressLevel)
//
//			rt := mw(http.DefaultTransport)
//			rresp, err := rt.RoundTrip(req)
//			fmt.Println(err)
//			_ = rresp
//			//_ = rresp.Body.Close()
//			//_ = err
//			//!!! Что-то непонятное с возвращающимся err творится, она даже не печатается
//			//fmt.Println(err)
//			//if (err != nil) != tt.wantErr {
//			//	t.Errorf("Compress() error - %v, wantErr - %v", err != nil, tt.wantErr)
//			//	return
//			//}
//
//			r, err := gzip.NewReader(req.Body)
//			if err != nil {
//				log.Fatal(err)
//			}
//			var buf bytes.Buffer
//			_, err = buf.ReadFrom(r)
//			if err != nil {
//				log.Fatal(err)
//			}
//			_ = r.Close()
//			unzipped := buf.String()
//			if !reflect.DeepEqual(tt.bodyString, unzipped) {
//				t.Errorf("Compress() = %v, want %v", tt.bodyString, unzipped)
//			}
//		})
//	}
//}
