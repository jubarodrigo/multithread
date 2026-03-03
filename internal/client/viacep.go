package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"multithread/internal/domain"
)

const viaCEPBaseURL = "http://viacep.com.br/ws"

// viaCEPResponse representa o JSON retornado pela ViaCEP.
type viaCEPResponse struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
}

// ViaCEP client para consulta de CEP na ViaCEP.
type ViaCEP struct {
	client  *http.Client
	baseURL string
}

// NewViaCEP cria um novo client da ViaCEP.
func NewViaCEP(timeout time.Duration) *ViaCEP {
	return &ViaCEP{
		client:  &http.Client{Timeout: timeout},
		baseURL: viaCEPBaseURL,
	}
}

// FetchByCEP consulta o CEP e retorna o endereço no formato do domínio.
func (c *ViaCEP) FetchByCEP(cep string) (*domain.Address, error) {
	url := c.baseURL + "/" + cep + "/json/"
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

	var body viaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("decodificar resposta: %w", err)
	}

	return &domain.Address{
		CEP:         body.CEP,
		Street:      body.Logradouro,
		Complement:  body.Complemento,
		Neighborhood: body.Bairro,
		City:        body.Localidade,
		State:       body.UF,
		Source:      "ViaCEP",
	}, nil
}
