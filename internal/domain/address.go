package domain

// Address representa o endereço unificado retornado por qualquer API de CEP.
type Address struct {
	CEP         string `json:"cep"`
	Street      string `json:"logradouro"`
	Complement  string `json:"complemento,omitempty"`
	Neighborhood string `json:"bairro"`
	City        string `json:"localidade"`
	State       string `json:"uf"`
	Source      string `json:"source"` // "BrasilAPI" ou "ViaCEP"
}
