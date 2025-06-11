package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCToF(t *testing.T) {
	type args struct {
		c float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "0°C para 32°F",
			args: args{c: 0},
			want: 32,
		},
		{
			name: "100°C para 212°F",
			args: args{c: 100},
			want: 212,
		},
		{
			name: "-40°C para -40°F",
			args: args{c: -40},
			want: -40,
		},
		{
			name: "37°C para 98.6°F",
			args: args{c: 37},
			want: 98.60000000000001,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, CToF(tt.args.c), "CToF(%v)", tt.args.c)
		})
	}
}

func TestCToK(t *testing.T) {
	type args struct {
		c float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "0°C para 273.15K",
			args: args{c: 0},
			want: 273.15,
		},
		{
			name: "100°C para 373.15K",
			args: args{c: 100},
			want: 373.15,
		},
		{
			name: "-273.15°C para 0K",
			args: args{c: -273.15},
			want: 0,
		},
		{
			name: "25°C para 298.15K",
			args: args{c: 25},
			want: 298.15,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, CToK(tt.args.c), "CToK(%v)", tt.args.c)
		})
	}
}

func TestWeatherHandler(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "CEP inválido (menos de 8 dígitos)",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/weather/1234567", nil),
			},
		},
		{
			name: "CEP válido, mas não encontrado",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/weather/00000000", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WeatherHandler(tt.args.w, tt.args.r)
		})
	}
}

func Test_getLocationByCEP(t *testing.T) {
	type args struct {
		cep string
	}
	tests := []struct {
		name    string
		args    args
		want    Location
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "CEP válido retorna localidade",
			args:    args{cep: "30140071"}, // CEP de Belo Horizonte
			want:    Location{Name: "Belo Horizonte"},
			wantErr: assert.NoError,
		},
		{
			name:    "CEP inválido retorna erro",
			args:    args{cep: "-00000001"},
			want:    Location{},
			wantErr: assert.Error,
		},
		{
			name:    "CEP com menos de 8 dígitos retorna erro",
			args:    args{cep: "123"},
			want:    Location{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getLocationByCEP(tt.args.cep)
			if !tt.wantErr(t, err, fmt.Sprintf("getLocationByCEP(%v)", tt.args.cep)) {
				return
			}
			assert.Equalf(t, tt.want, got, "getLocationByCEP(%v)", tt.args.cep)
		})
	}
}

func Test_getWeather(t *testing.T) {
	type args struct {
		city string
	}
	tests := []struct {
		name    string
		args    args
		want    WeatherResponse
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "Cidade válida retorna dados do tempo",
			args:    args{city: "Belo Horizonte"},
			want:    WeatherResponse{Location: Location{Name: "Belo Horizonte"}, Current: Current{TempC: 20.2}},
			wantErr: assert.NoError,
		},
		{
			name:    "Cidade inválida retorna erro",
			args:    args{city: "CidadeInexistente"},
			want:    WeatherResponse{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getWeather(tt.args.city)
			if !tt.wantErr(t, err, fmt.Sprintf("getWeather(%v)", tt.args.city)) {
				return
			}
			assert.Equalf(t, tt.want, got, "getWeather(%v)", tt.args.city)
		})
	}
}
