package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
)

type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
}

func GetCityFromCEP(ctx context.Context, cep string) (string, error) {
	tracer := otel.Tracer("serviceB")
	ctx, span := tracer.Start(ctx, "GetCityFromCEP")
	defer span.End()

	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 400 || resp.StatusCode == 404 {
		return "", errors.New("cep not found")
	}

	var r ViaCEPResponse
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return "", err
	}

	if r.Localidade == "" {
		return "", errors.New("cep not found")
	}

	return r.Localidade, nil
}
