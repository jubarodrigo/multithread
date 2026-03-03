package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"multithread/internal/domain"
	"multithread/internal/service"
)

const timeout = 1 * time.Second

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "uso: cepfinder <CEP>")
		fmt.Fprintln(os.Stderr, "exemplo: cepfinder 01310100")
		os.Exit(1)
	}

	cep := strings.TrimSpace(os.Args[1])
	if cep == "" {
		fmt.Fprintln(os.Stderr, "erro: CEP não pode ser vazio")
		os.Exit(1)
	}

	cep = strings.ReplaceAll(cep, "-", "")
	if len(cep) != 8 {
		fmt.Fprintln(os.Stderr, "erro: CEP deve ter 8 dígitos")
		os.Exit(1)
	}

	svc := service.NewCEPService(timeout)
	ctx := context.Background()

	addr, err := svc.FetchByCEP(ctx, cep)
	if err != nil {
		if errors.Is(err, service.ErrTimeout) {
			fmt.Fprintln(os.Stderr, "erro: timeout - nenhuma API respondeu em 1 segundo")
		} else {
			fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		}
		os.Exit(1)
	}

	printAddress(addr)
}

func printAddress(a *domain.Address) {
	fmt.Println("--- Endereço ---")
	fmt.Printf("CEP:         %s\n", a.CEP)
	fmt.Printf("Logradouro:  %s\n", a.Street)
	if a.Complement != "" {
		fmt.Printf("Complemento: %s\n", a.Complement)
	}
	fmt.Printf("Bairro:      %s\n", a.Neighborhood)
	fmt.Printf("Localidade:  %s\n", a.City)
	fmt.Printf("UF:          %s\n", a.State)
	fmt.Println("---------------")
	fmt.Printf("Fonte:       %s\n", a.Source)
}
