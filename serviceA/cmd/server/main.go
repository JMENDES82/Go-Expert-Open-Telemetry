package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/JMENDES82/Go-Expert-Open-Telemetry/serviceA/internal/handler"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	_ = godotenv.Load()

	// OTEL tracer
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
	aHandler := handler.NewAHandler()
	router.HandleFunc("/input", aHandler.PostCEP).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Service A starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
