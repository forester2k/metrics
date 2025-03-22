package main

import (
	"testing"
)

func Test_pollingInit(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "refreshRand",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pollingInit()
			first := random
			pollingInit()
			second := random
			if first == second {
				t.Errorf("before compress %v, after compress-decompress %v", first, second)
			}
		})
	}
}
