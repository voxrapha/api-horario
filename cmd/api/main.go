package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/zzRaphazz/api-horario/docs"
	"github.com/zzRaphazz/api-horario/internal/handler"
)

// @title           API Horário Mundial
// @version         1.0
// @description     API REST para consulta de data e hora atual por timezone IANA.
// @host            localhost:8080
// @BasePath        /
func main() {
	server := newServer(os.Getenv("PORT"))
	log.Printf("Servidor iniciado na porta http://localhost%s...", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}

func newServer(port string) *http.Server {
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.HealthCheckHandler)
	mux.HandleFunc("GET /api/v1/horario", handler.GetHorarioHandler)
	mux.HandleFunc("GET /api/v1/timezones", handler.GetTimezonesHandler)
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	return &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}
}
