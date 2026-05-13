package webmodel

import "time"

// RespostaNotebook para ter o JSON de um único notebook.
type RespostaNotebook struct {
	Nome       string  `json:"name"`
	Preco      float64 `json:"price"`
	Moeda      string  `json:"currency"`
	Descricao  string  `json:"description"`
	Avaliacao  int     `json:"rating"`
	Avaliacoes int     `json:"reviews"`
	URL        string  `json:"url"`
}

// RespostaAPINotebooks resposta completa da API.
type RespostaAPINotebooks struct {
	Total     int                `json:"total"`
	Fonte     string             `json:"source"`
	RaspadoEm time.Time          `json:"scraped_at"`
	EmCache   bool               `json:"cached"`
	Dados     []RespostaNotebook `json:"data"`
}

// RespostaErro resposta de erro (status 400 ou 500).
type RespostaErro struct {
	Erro string `json:"error"`
}
