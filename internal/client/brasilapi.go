package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"multithread/internal/domain"
)

const brasilAPIBaseURL = "https://brasilapi.com.br/api/cep/v1"

// brasilAPIResponse representa o JSON retornado pela BrasilAPI.
type brasilAPIResponse struct {
	CEP         string `json:"cep"`
	State       string `json:"state"`
	City        string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street      string `json:"street"`
	Service     string `json:"service"`
}

// BrasilAPI client para consulta de CEP na BrasilAPI.
type BrasilAPI struct {
	client  *http.Client
	baseURL string
}

// NewBrasilAPI cria um novo client da BrasilAPI.
func NewBrasilAPI(timeout time.Duration) *BrasilAPI {
	return &BrasilAPI{
		client:  &http.Client{Timeout: timeout},
		baseURL: brasilAPIBaseURL,
	}
}

// FetchByCEP consulta o CEP e retorna o endereço no formato do domínio.
func (c *BrasilAPI) FetchByCEP(cep string) (*domain.Address, error) {
	url := c.baseURL + "/" + cep
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("criar requisição: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	var body brasilAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("decodificar resposta: %w", err)
	}

	return &domain.Address{
		CEP:         body.CEP,
		Street:      body.Street,
		Neighborhood: body.Neighborhood,
		City:        body.City,
		State:       body.State,
		Source:      "BrasilAPI",
	}, nil
}
