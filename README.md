# 🕒 API Horário Mundial (Go / Golang)

API REST desenvolvida em **Go (Golang)** para consulta de data (`DD/MM/AAAA`) e hora (`HH:MM:SS`) atuais usando identificadores de fuso horário **IANA** (ex: `America/Sao_Paulo`, `Europe/London`, `Asia/Tokyo`). O endpoint de listagem mostra os fusos disponíveis; várias cidades podem compartilhar o mesmo fuso.

Projeto construído com foco em boas práticas de engenharia de software, arquitetura limpa, testes automatizados, documentação e preparação para deploy em nuvem (VM/EC2).

---tree arc

api-horario/
├── .github/
│   └── workflows/
│       └── ci.yml
├── cmd/
│   └── api/
│       ├── main.go
│       └── main_test.go
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal/
│   ├── handler/
│   │   ├── horario_handler.go
│   │   └── horario_handler_test.go
│   └── service/
│       ├── horario_service.go
│       ├── horario_service_test.go
│       ├── timezones.go
│       └── timezones.txt
├── tests/
│   └── load/
│       └── horario.js
├── .gitignore
├── go.mod
├── go.sum
├── README.md
└── sonar-project.properties


## 🛠️ Tecnologias e Padrões Utilizados

- **Linguagem:** Go 1.27+ (Biblioteca padrão `net/http`, `time`, `encoding/json`)
- **Testes:** `testing` e `net/http/httptest` com *Table-Driven Tests*
- **Teste de carga:** k6
- **Qualidade & CI/CD:** GitHub Actions + SonarCloud
- **Documentação:** Swagger (`swaggo/swag`)

---

## 🚀 Como Executar Localmente

### Pré-requisitos
- [Go 1.27 ou superior](https://go.dev) instalado.
- VS Code com a extensão oficial **Go** (`golang.go`).

### Passos

1. **Clone o repositório ou acesse a pasta do projeto:**
   ```bash
   cd api-horario
   ```

2. **Inicie o servidor localmente:**
   ```bash
   go run cmd/api/main.go
   ```
   O servidor estará disponível em `http://localhost:8080/health`
   http://localhost:8080/api/v1/horario?timezone=America/Sao_Paulo.

---

## 📡 Documentação dos Endpoints

### 1. Consultar Horário Atual
Retorna a data e a hora atual formatadas para o fuso horário IANA informado.

- **URL:** `/api/v1/horario`
- **Método:** `GET`
- **Query Parameters:**
  - `timezone` *(obrigatório)*: Identificador no padrão IANA (ex: `America/Sao_Paulo`).

#### Exemplo de Requisição:
```bash
curl -i "http://localhost:8080/api/v1/horario?timezone=America/Sao_Paulo"
```

#### Resposta de Sucesso (`200 OK`):
```json
{
  "timezone": "America/Sao_Paulo",
  "data": "12/11/2026",
  "hora": "01:05:30"
}
```

#### Resposta de Erro (`400 Bad Request`):
```json
{
  "erro": "timezone inválido: utilize o padrão IANA (ex: America/Sao_Paulo)"
}
```

---

### 2. Listar fusos horários IANA
Retorna os identificadores IANA aceitos pela API em ordem alfabética. A lista inclui os identificadores e aliases legados distribuídos com a base de fusos do Go 1.27.4; não representa todas as cidades. Mais de uma cidade pode compartilhar o mesmo fuso.

- **URL:** `/api/v1/timezones`
- **Método:** `GET`

#### Exemplo de Requisição:
```bash
curl -i "http://localhost:8080/api/v1/timezones"
```

#### Resposta de Sucesso (`200 OK`):
```json
{
  "timezones": [
    "Africa/Abidjan",
    "Africa/Accra",
    "Africa/Addis_Ababa"
  ]
}
```

Use um valor retornado no endpoint para consultar o horário, por exemplo `Africa/Tunis` ou `America/Maceio`. O arquivo `internal/service/timezones.txt` contém a lista versionada. Para atualizá-la junto com o banco de fusos do toolchain Go:

```bash
unzip -Z1 "$(go env GOROOT)/lib/time/zoneinfo.zip" | LC_ALL=C sort -u > internal/service/timezones.txt
```

---

### 3. Health Check
Verifica se a aplicação está saudável e respondendo.

- **URL:** `/health`
- **Método:** `GET`
- **Resposta (`200 OK`):** `{"status":"ok"}`

---

## 🧪 Executando os Testes Unitários

Para rodar todos os testes e visualizar o relatório de cobertura no terminal:

```bash
go test -v -cover ./...
```

Para gerar o relatório visual de cobertura em HTML no navegador:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

```bash
go test -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

A cobertura reportada pelo SonarCloud deve refletir o código Go da aplicação. Os arquivos gerados em `docs/` e o script k6 em `tests/load/` não são código de produção e ficam excluídos da análise de cobertura.

### Teste de carga com k6

Os testes unitários validam os casos funcionais. Para medir a API por HTTP sob carga local, instale o [k6](https://grafana.com/docs/k6/latest/set-up/install-k6/) e inicie a aplicação:

```bash
go run ./cmd/api
```

Em outro terminal, execute o perfil padrão de 600 requisições por minuto durante 1 minuto:

```bash
k6 run tests/load/horario.js
```

Taxa, duração e URL podem ser ajustadas por variáveis de ambiente:

```bash
TARGET_RPM=1200 DURATION=2m BASE_URL=http://localhost:8080 k6 run tests/load/horario.js
```

O teste verifica status e formato da resposta, além de aplicar limites de falhas, latência p95 e iterações descartadas. Execute-o contra a própria máquina ou um ambiente autorizado, não como parte dos testes unitários rápidos.

---

## 📖 Documentação Swagger

Com a API em execução, acesse a interface interativa em:
`http://localhost:8080/swagger/index.html`

Para regenerar os arquivos da especificação após alterar as anotações:

1. Instale a CLI do Swag:
   ```bash
   go install github.com/swaggo/swag/cmd/swag@v1.16.4
   ```
   Se o comando `swag` não for encontrado, adicione `$(go env GOPATH)/bin` ao `PATH`.
2. Gere os arquivos de documentação na pasta `docs/`:
   ```bash
   swag init -d ./cmd/api,./internal/handler,./internal/service -g main.go -o docs
   ```

Os arquivos gerados devem ser mantidos junto ao código para que a interface Swagger funcione durante a execução.

---

## ☁️ Como Compilar para Deploy em VM Linux (EC2 / Compute Engine)

O Go compila toda a aplicação em um **único binário executável**, dispensando a instalação de runtime na máquina virtual.

1. **Gere o binário para Linux (64-bit):**
   ```bash
   GOOS=linux GOARCH=amd64 go build -o api-horario cmd/api/main.go
   ```

2. **Envie o arquivo para sua VM via SCP:**
   ```bash
   scp -i sua-chave.pem api-horario ubuntu@IP_DA_SUA_VM:/home/ubuntu/
   ```

3. **Execute na VM:**
   ```bash
   chmod +x api-horario
   PORT=8080 ./api-horario
   ```