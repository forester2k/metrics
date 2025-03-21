package compression

import (
	"bytes"
	"compress/gzip"
	"log"
	"reflect"
	"testing"
)

func TestCompressGzip(t *testing.T) {
	type args struct {
		data          []byte
		compressLevel int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "compres-decompres",
			args: args{data: []byte{72, 101, 108, 108, 11},
				compressLevel: 1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zipped, err := CompressGzip(tt.args.data, tt.args.compressLevel)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompressGzip() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			b := bytes.NewReader(zipped)
			r, err := gzip.NewReader(b)
			if err != nil {
				log.Fatal(err)
			}
			var res bytes.Buffer
			_, err = res.ReadFrom(r)
			if err != nil {
				log.Fatal(err)
			}
			_ = r.Close()
			unzipped := res.Bytes()
			if !reflect.DeepEqual(tt.args.data, unzipped) {
				t.Errorf("before compress %v, after compress-decompress %v", zipped, unzipped)
			}
		})
	}
}
