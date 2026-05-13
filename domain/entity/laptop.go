package entity

// Notebook representa um produto notebook extraído do site alvo.
type Notebook struct {
	Nome       string
	Preco      float64
	Moeda      string
	Descricao  string
	Avaliacao  int
	Avaliacoes int
	URL        string
}
