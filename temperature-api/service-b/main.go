package main

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	otlptrace "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"log"
	"net/http"
	"time"
)

const viaCEPURL = "https://viacep.com.br/ws/"
const weatherAPIURL = "http://api.weatherapi.com/v1/current.json?key=d12a877c0d564ef784b14428250606&q="

func main() {
	tp, err := initTracer()
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatalf("Error shutting down tracer: %v", err)
		}
	}()

	r := gin.New()
	r.Use(otelgin.Middleware("service-b"))

	r.POST("/temperature", handleTemperature)

	if err := r.Run(":3001"); err != nil {
		log.Fatalf("Unable to start server: %v", err)
	}
}

func initTracer() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()

	exporter, err := otlptrace.New(ctx,
		otlptrace.WithInsecure(),
		otlptrace.WithEndpoint("otel-collector:4317"),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter, sdktrace.WithBatchTimeout(500*time.Millisecond)),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("service-b"),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

func handleTemperature(c *gin.Context) {
	_, span := otel.Tracer("service-b").Start(c.Request.Context(), "HandleTemperature")
	defer span.End()

	type CepRequest struct {
		Cep string `json:"cep"`
	}
	var req CepRequest
	if err := c.ShouldBindJSON(&req); err != nil || len(req.Cep) != 8 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "invalid zipcode"})
		span.SetStatus(codes.Error, "Invalid zipcode")
		return
	}

	locationResp, err := http.Get(viaCEPURL + req.Cep + "/json/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		span.SetStatus(codes.Error, "Failed to fetch location")
		return
	}
	defer locationResp.Body.Close()

	if locationResp.StatusCode != http.StatusOK {
		c.JSON(http.StatusNotFound, gin.H{"message": "can not find zipcode"})
		span.SetStatus(codes.Error, "Zipcode not found")
		return
	}

	var locationData map[string]interface{}
	json.NewDecoder(locationResp.Body).Decode(&locationData)
	city := locationData["localidade"].(string)

	tempResp, err := http.Get(weatherAPIURL + city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		span.SetStatus(codes.Error, "Failed to fetch temperature")
		return
	}
	defer tempResp.Body.Close()

	var weatherData struct {
		Current struct {
			TempC float64 `json:"temp_c"`
		} `json:"current"`
	}
	json.NewDecoder(tempResp.Body).Decode(&weatherData)

	tempC := weatherData.Current.TempC
	tempF := (tempC * 1.8) + 32
	tempK := tempC + 273.15

	c.JSON(http.StatusOK, gin.H{
		"city":   city,
		"temp_C": tempC,
		"temp_F": tempF,
		"temp_K": tempK,
	})
}
