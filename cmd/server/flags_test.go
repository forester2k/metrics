package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// Валидирует входящий путь [path] и проверяет существует ли файл по заданному пути.
// Если файл существует, то отдает строку с путем и именем файла.
// Если файл не существует или произошла ошибка во время выполнения функции, то
// отдает пустую строку
func checkFile(path string) (string, error) {
	path = filepath.Clean(path)
	var err error
	path, err = filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("checkFile: Ошибка при получении абсолютного пути %s: %w", path, err)
	}
	dir := filepath.Dir(path)
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return "", fmt.Errorf("checkFile: Ошибка при создании директории %s: %w", dir, err)
	}
	if _, err = os.Stat(path); err != nil {
		//файл не найден, но и ошибок не произошло
		return "", nil
	}
	return path, nil
}

func deleteFile(path string) error {
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("deleteFile: не удается удалить файл %s, %w", path, err)
	}
	return nil
}

func Test_createPathAndFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "rightPath",
			args: args{
				path: "./temp9352.tmp",
			},
			want:    "./temp9352.tmp",
			wantErr: false,
		},
		//{
		//	name: "wrongPath",
		//	args: args{
		//		path: string([]byte{0}),
		//	},
		//	want:    string([]byte{0}),
		//	wantErr: true,
		//},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Проверим, не осталось ли файла от прежних тестов, удалим его если о есть
			path, err := checkFile(tt.args.path)
			if err != nil {
				fmt.Printf("Во время проверки наличия файла %v произошла ошибка %v", path, err)
				panic(err)
			}
			if path != "" {
				if err := deleteFile(path); err != nil {
					fmt.Printf("не удалось удалить файл %v, ошибка %v", path, err)
					panic(err)
				}
			}
			gotPath, err := createPathAndFile(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("createPathAndFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			absWant, err := filepath.Abs(filepath.Clean(tt.want))
			if err != nil {
				fmt.Printf("Test_createPathAndFile: Ошибка при получении абсолютного пути %s: %v\n", tt.want, err)
				panic(err)
			}
			if gotPath != absWant {
				t.Errorf("createPathAndFile() got = %v, want %v", gotPath, absWant)
			}
			err = deleteFile(gotPath)
			if err != nil {
				fmt.Printf("cant delete file %v afer testing\n", gotPath)
			}
		})
	}
}

//func Test_parseFlags(t *testing.T) {
//	tests := []struct {
//		name    string
//		wantErr bool
//	}{
//		{
//			name:    "standardExecute",
//			wantErr: false,
//		},
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if err := parseFlags(); (err != nil) != tt.wantErr {
//				t.Errorf("parseFlags() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
