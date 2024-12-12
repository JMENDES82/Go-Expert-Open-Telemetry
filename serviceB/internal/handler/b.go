package handler

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/JMENDES82/Go-Expert-Open-Telemetry/serviceB/internal/model"
	"github.com/JMENDES82/Go-Expert-Open-Telemetry/serviceB/internal/service"
	"github.com/JMENDES82/Go-Expert-Open-Telemetry/serviceB/internal/util"

	"go.opentelemetry.io/otel"
)

type BHandler struct{}

func NewBHandler() *BHandler {
	return &BHandler{}
}

func (h *BHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("serviceB")
	ctx, span := tracer.Start(ctx, "GetWeatherHandler")
	defer span.End()

	cep := r.URL.Query().Get("cep")
	if cep == "" || !isValidCEP(cep) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("invalid zipcode"))
		return
	}

	city, err := service.GetCityFromCEP(ctx, cep)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("can not find zipcode"))
		return
	}

	tempC, err := service.GetCurrentTemperature(ctx, city)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("can not find zipcode"))
		return
	}

	resp := model.WeatherResponse{
		City:  city,
		TempC: tempC,
		TempF: util.CelsiusToFahrenheit(tempC),
		TempK: util.CelsiusToKelvin(tempC),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func isValidCEP(cep string) bool {
	match, _ := regexp.MatchString(`^\d{8}$`, cep)
	return match
}
