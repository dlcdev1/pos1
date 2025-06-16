package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	ohttp "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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

const serviceBURL = "http://service-b:3001/temperature"

type CepRequest struct {
	Cep string `json:"cep"`
}

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
	r.Use(otelgin.Middleware("service-a"))
	r.POST("/cep", handleCep)

	if err := r.Run(":3000"); err != nil {
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
			semconv.ServiceNameKey.String("service-a"),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

func handleCep(c *gin.Context) {
	ctx, span := otel.Tracer("service-a").Start(c.Request.Context(), "HandleCep")
	defer span.End()

	var req CepRequest
	if err := c.ShouldBindJSON(&req); err != nil || len(req.Cep) != 8 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "invalid zipcode"})
		span.SetStatus(codes.Error, "Invalid zipcode")
		return
	}

	body, _ := json.Marshal(req)

	client := http.Client{Transport: ohttp.NewTransport(http.DefaultTransport)}

	httpreq, err := http.NewRequestWithContext(ctx, "POST", serviceBURL, bytes.NewBuffer(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		span.SetStatus(codes.Error, "Failed to prepare request")
		return
	}
	httpreq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(httpreq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		span.SetStatus(codes.Error, "HTTP Post error")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse gin.H
		json.NewDecoder(resp.Body).Decode(&errorResponse)
		c.JSON(resp.StatusCode, errorResponse)
		span.SetStatus(codes.Error, "Service B returned error")
		return
	}

	var locationData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&locationData)
	c.JSON(http.StatusOK, locationData)
}
