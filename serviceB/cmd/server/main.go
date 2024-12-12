package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/JMENDES82/Go-Expert-Open-Telemetry/serviceB/internal/handler"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	_ = godotenv.Load()

	// Inicia o tracer do OTEL
	ctx := context.Background()
	exporter, err := otlptracehttp.New(ctx)
	if err != nil {
		log.Fatal(err)
	}

	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	otel.SetTracerProvider(tp)

	router := mux.NewRouter()
	bHandler := handler.NewBHandler()
	router.HandleFunc("/weather", bHandler.GetWeather).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Service B starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
