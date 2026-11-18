package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zzRaphazz/api-horario/internal/service"
)

func TestGetHorarioHandler(t *testing.T) {
	testes := []struct {
		nome           string
		url            string
		statusEsperado int
		timezone       string
		erroEsperado   string
	}{
		{
			nome:           "timezone valido",
			url:            "/api/v1/horario?timezone=Europe/London",
			statusEsperado: http.StatusOK,
			timezone:       "Europe/London",
		},
		{
			nome:           "timezone com espacos nas extremidades",
			url:            "/api/v1/horario?timezone=%20America%2FSao_Paulo%20",
			statusEsperado: http.StatusOK,
			timezone:       "America/Sao_Paulo",
		},
		{
			nome:           "timezone ausente",
			url:            "/api/v1/horario",
			statusEsperado: http.StatusBadRequest,
			erroEsperado:   "o parâmetro timezone é obrigatório",
		},
		{
			nome:           "timezone invalido",
			url:            "/api/v1/horario?timezone=Planeta%2FMarte",
			statusEsperado: http.StatusBadRequest,
			erroEsperado:   "timezone inválido: utilize o padrão IANA (ex: America/Sao_Paulo)",
		},
		{
			nome:           "nome de localidade que nao e identificador IANA",
			url:            "/api/v1/horario?timezone=America%2FSergipe",
			statusEsperado: http.StatusBadRequest,
			erroEsperado:   "timezone inválido: utilize o padrão IANA (ex: America/Sao_Paulo)",
		},
	}

	for _, teste := range testes {
		t.Run(teste.nome, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, teste.url, nil)
			rec := httptest.NewRecorder()

			GetHorarioHandler(rec, req)

			if rec.Code != teste.statusEsperado {
				t.Fatalf("status esperado %d, recebeu %d", teste.statusEsperado, rec.Code)
			}
			if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
				t.Errorf("Content-Type esperado application/json, recebeu %q", contentType)
			}

			if teste.statusEsperado == http.StatusOK {
				var resp service.HorarioResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
					t.Fatalf("erro ao decodificar JSON: %v", err)
				}
				if resp.Timezone != teste.timezone || resp.Data == "" || resp.Hora == "" {
					t.Errorf("resposta incompleta ou inesperada: %+v", resp)
				}
				return
			}

			var resp ErroResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
				t.Fatalf("erro ao decodificar JSON de erro: %v", err)
			}
			if resp.Erro != teste.erroEsperado {
				t.Errorf("erro esperado %q, recebeu %q", teste.erroEsperado, resp.Erro)
			}
		})
	}
}

func TestHealthCheckHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	HealthCheckHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status esperado %d, recebeu %d", http.StatusOK, rec.Code)
	}

	if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type esperado application/json, recebeu %q", contentType)
	}

	var resposta map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &resposta); err != nil {
		t.Fatalf("erro ao decodificar JSON: %v", err)
	}
	if resposta["status"] != "ok" {
		t.Errorf("status esperado %q, recebeu %q", "ok", resposta["status"])
	}
}

func TestGetTimezonesHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/timezones", nil)
	rec := httptest.NewRecorder()

	GetTimezonesHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status esperado %d, recebeu %d", http.StatusOK, rec.Code)
	}
	if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type esperado application/json, recebeu %q", contentType)
	}

	var resposta service.TimezonesResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resposta); err != nil {
		t.Fatalf("erro ao decodificar JSON: %v", err)
	}
	if len(resposta.Timezones) < 500 {
		t.Errorf("esperava lista ampla de fusos IANA, recebeu %d", len(resposta.Timezones))
	}
}
