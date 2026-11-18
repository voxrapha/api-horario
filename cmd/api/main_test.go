package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewServerUsesDefaultPort(t *testing.T) {
	server := newServer("")

	if server.Addr != ":8080" {
		t.Fatalf("endereco esperado :8080, recebeu %q", server.Addr)
	}
	if server.Handler == nil {
		t.Fatal("handler do servidor não deve ser nil")
	}
}

func TestNewServerUsesConfiguredPort(t *testing.T) {
	server := newServer("9090")

	if server.Addr != ":9090" {
		t.Fatalf("endereco esperado :9090, recebeu %q", server.Addr)
	}
}

func TestServerRoutes(t *testing.T) {
	server := newServer("8080")

	testes := []struct {
		nome           string
		path           string
		statusEsperado int
		contentType    string
	}{
		{
			nome:           "health",
			path:           "/health",
			statusEsperado: http.StatusOK,
			contentType:    "application/json",
		},
		{
			nome:           "consulta de horario",
			path:           "/api/v1/horario?timezone=America%2FSao_Paulo",
			statusEsperado: http.StatusOK,
			contentType:    "application/json",
		},
		{
			nome:           "lista de fusos IANA",
			path:           "/api/v1/timezones",
			statusEsperado: http.StatusOK,
			contentType:    "application/json",
		},
		{
			nome:           "interface swagger",
			path:           "/swagger/index.html",
			statusEsperado: http.StatusOK,
			contentType:    "text/html; charset=utf-8",
		},
		{
			nome:           "rota inexistente",
			path:           "/rota-inexistente",
			statusEsperado: http.StatusNotFound,
		},
	}

	for _, teste := range testes {
		t.Run(teste.nome, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, teste.path, nil)
			response := httptest.NewRecorder()

			server.Handler.ServeHTTP(response, request)

			if response.Code != teste.statusEsperado {
				t.Fatalf("status esperado %d, recebeu %d", teste.statusEsperado, response.Code)
			}
			if teste.contentType != "" && response.Header().Get("Content-Type") != teste.contentType {
				t.Errorf(
					"Content-Type esperado %q, recebeu %q",
					teste.contentType,
					response.Header().Get("Content-Type"),
				)
			}
		})
	}
}
