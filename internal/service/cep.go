package service

import (
	"context"
	"errors"
	"time"

	"multithread/internal/client"
	"multithread/internal/domain"
)

const defaultTimeout = 1 * time.Second

// ErrTimeout é retornado quando nenhuma API responde dentro do tempo limite.
var ErrTimeout = errors.New("timeout: nenhuma API respondeu em 1 segundo")

// CEPService orquestra a busca de CEP em múltiplas APIs com corrida (race).
type CEPService struct {
	brasilAPI *client.BrasilAPI
	viaCEP    *client.ViaCEP
	timeout   time.Duration
}

// NewCEPService cria um novo serviço de consulta de CEP.
func NewCEPService(timeout time.Duration) *CEPService {
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	return &CEPService{
		brasilAPI: client.NewBrasilAPI(timeout),
		viaCEP:    client.NewViaCEP(timeout),
		timeout:   timeout,
	}
}

// FetchByCEP consulta o CEP nas duas APIs em paralelo e retorna a resposta
// da que responder primeiro. Descarta a resposta da mais lenta.
// Retorna ErrTimeout se nenhuma responder dentro do tempo limite.
func (s *CEPService) FetchByCEP(ctx context.Context, cep string) (*domain.Address, error) {
	// Canal com buffer 2 para que ambas as goroutines possam enviar sem bloquear
	// (aceitamos apenas a primeira resposta; a segunda pode ser descartada no buffer).
	ch := make(chan *domain.Address, 2)

	fetch := func(do func() (*domain.Address, error)) {
		addr, err := do()
		if err != nil {
			return
		}
		select {
		case ch <- addr:
		default:
			// Canal já preenchido ou fechado; descarta (outra API já ganhou).
		}
	}

	go fetch(func() (*domain.Address, error) { return s.brasilAPI.FetchByCEP(cep) })
	go fetch(func() (*domain.Address, error) { return s.viaCEP.FetchByCEP(cep) })

	select {
	case addr := <-ch:
		return addr, nil
	case <-time.After(s.timeout):
		return nil, ErrTimeout
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
