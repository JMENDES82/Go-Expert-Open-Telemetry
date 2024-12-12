package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"regexp"

	"go.opentelemetry.io/otel"
)

type AHandler struct {
	BServiceURL string
}

func NewAHandler() *AHandler {
	// Supondo que o Service B esteja acessível via ENV ou URL fixa
	// Ajuste conforme necessário.
	return &AHandler{
		BServiceURL: "http://service_b:8080/weather",
	}
}

type CEPInput struct {
	CEP string `json:"cep"`
}

func (h *AHandler) PostCEP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("serviceA")
	ctx, span := tracer.Start(ctx, "PostCEPHandler")
	defer span.End()

	var input CEPInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("invalid zipcode"))
		return
	}

	if !isValidCEP(input.CEP) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("invalid zipcode"))
		return
	}

	// Chamar serviço B
	result, err := h.callServiceB(ctx, input.CEP)
	if err != nil {
		// Se erro, apenas repassamos o status e a mensagem
		if errResp, ok := err.(HttpError); ok {
			w.WriteHeader(errResp.Code)
			w.Write([]byte(errResp.Msg))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal error"))
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func (h *AHandler) callServiceB(ctx context.Context, cep string) ([]byte, error) {
	tracer := otel.Tracer("serviceA")
	ctx, span := tracer.Start(ctx, "callServiceB")
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, "GET", h.BServiceURL+"?cep="+cep, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, HttpError{Code: resp.StatusCode, Msg: string(bodyBytes)}
	}

	return bodyBytes, nil
}

func isValidCEP(cep string) bool {
	match, _ := regexp.MatchString(`^\d{8}$`, cep)
	return match
}

type HttpError struct {
	Code int
	Msg  string
}

func (e HttpError) Error() string {
	return e.Msg
}
