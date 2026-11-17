package service

import (
	"testing"
	"time"
)

func TestObterHorarioAtual(t *testing.T) {
	testes := []struct {
		nome       string
		timezone   string
		esperaErro bool
	}{
		{
			nome:       "Sucesso com America/Sao_Paulo",
			timezone:   "America/Sao_Paulo",
			esperaErro: false,
		},
		{
			nome:       "Sucesso com Asia/Tokyo",
			timezone:   "Asia/Tokyo",
			esperaErro: false,
		},
		{
			nome:       "Sucesso com UTC",
			timezone:   "UTC",
			esperaErro: false,
		},
		{
			nome:       "Erro quando timezone vazio",
			timezone:   "",
			esperaErro: true,
		},
		{
			nome:       "Erro quando timezone apenas espacos",
			timezone:   "   ",
			esperaErro: true,
		},
		{
			nome:       "Erro quando timezone inexistente",
			timezone:   "Planeta/Marte",
			esperaErro: true,
		},
	}

	for _, tc := range testes {
		t.Run(tc.nome, func(t *testing.T) {
			resp, err := ObterHorarioAtual(tc.timezone)

			if tc.esperaErro {
				if err == nil {
					t.Errorf("esperava erro para timezone '%s', mas não ocorreu", tc.timezone)
				}
				return
			}

			if err != nil {
				t.Fatalf("não esperava erro, mas recebeu: %v", err)
			}

			if resp.Timezone != tc.timezone {
				t.Errorf("timezone esperado %s, recebeu %s", tc.timezone, resp.Timezone)
			}

			// Valida o formato DD/MM/AAAA
			if _, err := time.Parse("02/01/2006", resp.Data); err != nil {
				t.Errorf("formato de data inválido (%s): %v", resp.Data, err)
			}

			// Valida o formato HH:MM:SS
			if _, err := time.Parse("15:04:05", resp.Hora); err != nil {
				t.Errorf("formato de hora inválido (%s): %v", resp.Hora, err)
			}
		})
	}
}
