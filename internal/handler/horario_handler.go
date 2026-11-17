package handler

import (
	"encoding/json"
	"net/http"

	"github.com/zzRaphazz/api-horario/internal/service"
)

// ErroResponse padroniza o JSON de erro da API
type ErroResponse struct {
	Erro string `json:"erro"`
}

// GetHorarioHandler godoc
// @Summary      Consulta data e hora atual
// @Description  Retorna a data (DD/MM/AAAA) e hora (HH:MM:SS) de um timezone no padrão IANA.
// @Tags         horario
// @Produce      json
// @Param        timezone  query     string  true  "Timezone IANA (ex: America/Sao_Paulo)"
// @Success      200       {object}  service.HorarioResponse
// @Failure      400       {object}  ErroResponse
// @Router       /api/v1/horario [get]
// GetHorarioHandler processa requisições GET /api/v1/horario
func GetHorarioHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	timezone := r.URL.Query().Get("timezone")

	resultado, err := service.ObterHorarioAtual(timezone)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErroResponse{Erro: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resultado)
}

// HealthCheckHandler godoc
// @Summary      Verifica a saúde da API
// @Description  Informa se a API está disponível.
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
// HealthCheckHandler verifica se a API está online (útil para monitoramento na VM)
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
