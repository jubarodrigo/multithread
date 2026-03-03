# Multithread – Busca de CEP em paralelo

Sistema em Go que consulta endereço por CEP em **duas APIs ao mesmo tempo** e retorna o resultado da que responder **mais rápido**, com timeout de 1 segundo.

## APIs utilizadas

- **BrasilAPI:** `https://brasilapi.com.br/api/cep/v1/{cep}`
- **ViaCEP:** `http://viacep.com.br/ws/{cep}/json/`

## Requisitos atendidos

- **Requisições simultâneas:** as duas APIs são chamadas em paralelo (goroutines).
- **Race (corrida):** apenas a primeira resposta é aceita; a outra é descartada.
- **Saída no terminal:** exibe os dados do endereço e qual API respondeu (BrasilAPI ou ViaCEP).
- **Timeout de 1 segundo:** se nenhuma API responder a tempo, é exibido erro de timeout.

## Estrutura do projeto (padrão Go)

```
multithread/
├── cmd/
│   └── cepfinder/       # Binário principal
│       └── main.go
├── internal/
│   ├── client/          # Clientes HTTP das APIs
│   │   ├── brasilapi.go
│   │   └── viacep.go
│   ├── domain/          # Modelo de endereço
│   │   └── address.go
│   └── service/         # Regra de negócio (race + timeout)
│       └── cep.go
├── go.mod
└── README.md
```

## Como rodar

### Pré-requisitos

- Go 1.21 ou superior

### Executar com `go run`

```bash
go run ./cmd/cepfinder <CEP>
```

Exemplo:

```bash
go run ./cmd/cepfinder 01310100
```

Ou com CEP formatado:

```bash
go run ./cmd/cepfinder 01310-100
```

### Compilar e executar o binário

```bash
go build -o cepfinder ./cmd/cepfinder
./cepfinder 01310100
```

## Exemplo de saída (sucesso)

```
--- Endereço ---
CEP:         01310100
Logradouro:  Avenida Paulista
Bairro:      Bela Vista
Localidade:  São Paulo
UF:          SP
---------------
Fonte:       BrasilAPI
```

(ou `Fonte: ViaCEP`, dependendo de qual API responder primeiro.)

## Exemplo de saída (timeout)

Se nenhuma API responder em 1 segundo:

```
erro: timeout - nenhuma API respondeu em 1 segundo
```

## Como testar

1. **CEP válido:** use um CEP real (ex.: `01310100`) e verifique que o endereço e a fonte (BrasilAPI ou ViaCEP) são exibidos.
2. **Timeout:** para simular timeout, desligue a rede ou use um CEP que demore (em ambiente controlado); o programa deve sair com a mensagem de timeout após 1 segundo.
3. **CEP inválido:** `go run ./cmd/cepfinder 00000000` pode retornar erro das APIs; o programa repassa o erro ao usuário.

## Tecnologias e conceitos

- **Linguagem:** Go (Golang)
- **Conceitos:** Goroutines, Channels, `select`, pacote `net/http`, timeout com `time.After`.
