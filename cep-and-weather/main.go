package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	viaCepURL     = "https://viacep.com.br/ws/%s/json/"
	weatherAPIURL = "http://api.weatherapi.com/v1/current.json?key=%s&q=%s"
	weatherAPIKey = "d12a877c0d564ef784b14428250606"
)

type WeatherResponse struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
}

type Location struct {
	Name string `json:"name"`
}

type Current struct {
	TempC float64 `json:"temp_c"`
}

type TemperatureResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/weather/{cep}", WeatherHandler).Methods("GET")
	http.Handle("/", r)
	fmt.Println("Servidor rodando na porta 8080...")
	http.ListenAndServe(":8080", nil)
}

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cep := vars["cep"]

	if len(cep) != 8 {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	location, err := getLocationByCEP(cep)
	if err != nil || location.Name == "" {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	weather, err := getWeather(location.Name)
	if err != nil {
		http.Error(w, "failed to get weather", http.StatusInternalServerError)
		return
	}

	response := TemperatureResponse{
		TempC: weather.Current.TempC,
		TempF: CToF(weather.Current.TempC),
		TempK: CToK(weather.Current.TempC),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func getLocationByCEP(cep string) (Location, error) {
	var location Location
	resp, err := http.Get(fmt.Sprintf(viaCepURL, cep))
	if err != nil {
		return location, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return location, fmt.Errorf("cep not found")
	}

	var result interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if loc, ok := result.(map[string]interface{}); ok {
		if loc["localidade"] != nil {
			location.Name = loc["localidade"].(string)
		}
	}
	return location, nil
}

func getWeather(city string) (WeatherResponse, error) {
	var weather WeatherResponse
	resp, err := http.Get(fmt.Sprintf(weatherAPIURL, "d12a877c0d564ef784b14428250606", city))
	if err != nil {
		return weather, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return weather, fmt.Errorf("weather not found")
	}

	json.NewDecoder(resp.Body).Decode(&weather)
	return weather, nil
}

func CToF(c float64) float64 {
	return c*1.8 + 32
}

func CToK(c float64) float64 {
	return c + 273.15
}
