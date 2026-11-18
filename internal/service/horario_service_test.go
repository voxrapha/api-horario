package service

import (
	"slices"
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
		{
			nome:       "Erro quando timezone local nao pertence a lista IANA",
			timezone:   "Local",
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

func TestListarTimezones(t *testing.T) {
	resposta := ListarTimezones()
	timezones := resposta.Timezones

	if len(timezones) < 500 {
		t.Fatalf("esperava pelo menos 500 identificadores IANA, recebeu %d", len(timezones))
	}
	if !slices.IsSorted(timezones) {
		t.Fatal("identificadores de timezone devem estar ordenados")
	}

	encontrados := make(map[string]bool, len(timezones))
	for _, timezone := range timezones {
		if encontrados[timezone] {
			t.Errorf("identificador duplicado: %s", timezone)
		}
		encontrados[timezone] = true

		if resposta, err := ObterHorarioAtual(timezone); err != nil {
			t.Errorf("identificador listado não pode ser usado na consulta (%s): %v", timezone, err)
		} else if resposta.Timezone != timezone {
			t.Errorf("timezone esperado %s, recebeu %s", timezone, resposta.Timezone)
		}
	}

	for _, timezone := range []string{
		"Africa/Tunis",
		"America/Maceio",
		"Europe/Madrid",
		"UTC",
		"US/Pacific",
	} {
		if !encontrados[timezone] {
			t.Errorf("identificador esperado não encontrado: %s", timezone)
		}
	}
	if encontrados["America/Sergipe"] {
		t.Error("lista oficial não deve incluir aliases próprios da API")
	}
}

func TestListarTimezonesReturnsIndependentSlice(t *testing.T) {
	primeiraResposta := ListarTimezones()
	primeiroTimezone := primeiraResposta.Timezones[0]
	primeiraResposta.Timezones[0] = "alterado"

	segundaResposta := ListarTimezones()
	if segundaResposta.Timezones[0] != primeiroTimezone {
		t.Fatal("modificar a resposta não deve alterar a lista compartilhada")
	}
}
