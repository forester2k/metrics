package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGaugeMetric_SetValue(t *testing.T) {
	type fields struct {
		Name  string
		Value float64
	}
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		wantRes float64
	}{
		{
			name: "rightValueType",
			fields: fields{
				Name:  "SomeName",
				Value: float64(3.14),
			},
			args: args{
				i: float64(2.01),
			},
			wantErr: false,
			wantRes: float64(2.01),
		},
		{
			name: "wrongValueType",
			fields: fields{
				Name:  "SomeName",
				Value: float64(3.14),
			},
			args: args{
				i: float32(2.01),
			},
			wantErr: true,
			wantRes: float64(3.14),
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &GaugeMetric{
				Name:  tt.fields.Name,
				Value: tt.fields.Value,
			}
			if err := m.SetValue(tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equalf(t, tt.wantRes, m.Value, "SetValue() error. Value want - %v, got - %v", tt.wantRes, m.Value)
		})
	}
}

func TestCounterMetric_SetValue(t *testing.T) {
	type fields struct {
		Name  string
		Value int64
	}
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		wantRes int64
	}{
		{
			name: "rightValueType",
			fields: fields{
				Name:  "SomeName",
				Value: int64(3),
			},
			args: args{
				i: int64(2),
			},
			wantErr: false,
			wantRes: int64(2),
		},
		{
			name: "wrongValueType",
			fields: fields{
				Name:  "SomeName",
				Value: int64(3),
			},
			args: args{
				i: int(2),
			},
			wantErr: true,
			wantRes: int64(3),
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &CounterMetric{
				Name:  tt.fields.Name,
				Value: tt.fields.Value,
			}
			if err := m.SetValue(tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equalf(t, tt.wantRes, m.Value, "SetValue() error. Value want - %v, got - %v", tt.wantRes, m.Value)
		})
	}
}
