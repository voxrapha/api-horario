package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	_ "time/tzdata"

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
	mux := http.NewServeMux()

	// Rotas usando o novo padrão do Go 1.22+ (Método + Caminho)
	mux.HandleFunc("GET /health", handler.HealthCheckHandler)
	mux.HandleFunc("GET /api/v1/horario", handler.GetHorarioHandler)
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	porta := os.Getenv("PORT")
	if porta == "" {
		porta = "8080"
	}

	endereco := fmt.Sprintf(":%s", porta)
	log.Printf("Servidor iniciado na porta http://localhost:%s...", porta)

	if err := http.ListenAndServe(endereco, mux); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
