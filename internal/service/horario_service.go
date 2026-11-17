package service

import (
	"errors"
	"strings"
	"time"
)

// HorarioResponse representa o payload de retorno da API
type HorarioResponse struct {
	Timezone string `json:"timezone"`
	Data     string `json:"data"`
	Hora     string `json:"hora"`
}

// ObterHorarioAtual busca a data e hora atual para um timezone IANA válido
func ObterHorarioAtual(timezone string) (HorarioResponse, error) {
	tz := strings.TrimSpace(timezone)
	if tz == "" {
		return HorarioResponse{}, errors.New("o parâmetro timezone é obrigatório")
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return HorarioResponse{}, errors.New("timezone inválido: utilize o padrão IANA (ex: America/Sao_Paulo)")
	}

	agora := time.Now().In(loc)

	return HorarioResponse{
		Timezone: tz,
		Data:     agora.Format("02/01/2006"), // Padrão DD/MM/AAAA
		Hora:     agora.Format("15:04:05"),   // Padrão HH:MM:SS
	}, nil
}
