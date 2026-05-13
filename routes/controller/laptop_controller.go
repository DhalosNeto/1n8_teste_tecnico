package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/DhalosNeto/1n8_teste_tecnico/application/webmodel"
	"github.com/DhalosNeto/1n8_teste_tecnico/domain/entity"
	"github.com/DhalosNeto/1n8_teste_tecnico/domain/service"
)

// ControladorNotebook gerencia os endpoints relacionados a notebooks.
type ControladorNotebook struct {
	servico *service.ServicoNotebook
}

// NovoControladorNotebook cria uma nova instância de ControladorNotebook.
func NovoControladorNotebook(servico *service.ServicoNotebook) *ControladorNotebook {
	return &ControladorNotebook{servico: servico}
}

// ObterLenovo retorna notebooks Lenovo.
func (c *ControladorNotebook) ObterLenovo(w http.ResponseWriter, r *http.Request) {
	opcoes, err := extrairOpcoesFiltro(r)
	if err != nil {
		escreverErro(w, http.StatusBadRequest, err.Error())
		return
	}

	resultado, err := c.servico.ObterLenovo(opcoes)
	if err != nil {
		escreverErro(w, http.StatusInternalServerError, "falha no scraping: "+err.Error())
		return
	}

	escreverJSON(w, http.StatusOK, converterParaRespostaAPI(resultado))
}

// Atualizar apaga o cache pra realizar outro scraper.
func (c *ControladorNotebook) Atualizar(w http.ResponseWriter, r *http.Request) {
	opcoes, err := extrairOpcoesFiltro(r)
	if err != nil {
		escreverErro(w, http.StatusBadRequest, err.Error())
		return
	}

	resultado, err := c.servico.Atualizar(opcoes)
	if err != nil {
		escreverErro(w, http.StatusInternalServerError, "falha no scraping: "+err.Error())
		return
	}

	escreverJSON(w, http.StatusOK, converterParaRespostaAPI(resultado))
}

// Saude retorna o status atual da aplicação.
func (c *ControladorNotebook) Saude(w http.ResponseWriter, r *http.Request) {
	escreverJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

func extrairOpcoesFiltro(r *http.Request) (service.OpcoesFiltro, error) {
	q := r.URL.Query()
	opcoes := service.OpcoesFiltro{}

	// valida o preço mínimo
	if v := q.Get("min_price"); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return opcoes, &erroParametro{"min_price", "deve ser um número válido"}
		}
		opcoes.PrecoMinimo = &f
	}

	// valida o preço máximo
	if v := q.Get("max_price"); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return opcoes, &erroParametro{"max_price", "deve ser um número válido"}
		}
		opcoes.PrecoMaximo = &f
	}

	// valida a avaliação mínima
	if v := q.Get("min_rating"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 || n > 5 {
			return opcoes, &erroParametro{"min_rating", "deve ser um inteiro entre 1 e 5"}
		}
		opcoes.AvaliacaoMinima = &n
	}

	return opcoes, nil
}

func converterParaRespostaAPI(resultado *service.ResultadoNotebook) webmodel.RespostaAPINotebooks {
	dados := make([]webmodel.RespostaNotebook, len(resultado.Notebooks))
	for i, n := range resultado.Notebooks {
		dados[i] = converterParaWebModel(n)
	}
	return webmodel.RespostaAPINotebooks{
		Total:     len(dados),
		Fonte:     resultado.Fonte,
		RaspadoEm: resultado.RaspadoEm.UTC(),
		EmCache:   resultado.EmCache,
		Dados:     dados,
	}
}

func converterParaWebModel(n entity.Notebook) webmodel.RespostaNotebook {
	return webmodel.RespostaNotebook{
		Nome:       n.Nome,
		Preco:      n.Preco,
		Moeda:      n.Moeda,
		Descricao:  n.Descricao,
		Avaliacao:  n.Avaliacao,
		Avaliacoes: n.Avaliacoes,
		URL:        n.URL,
	}
}

func escreverJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func escreverErro(w http.ResponseWriter, status int, mensagem string) {
	escreverJSON(w, status, webmodel.RespostaErro{Erro: mensagem})
}

type erroParametro struct {
	parametro string
	mensagem  string
}

func (e *erroParametro) Error() string {
	return "parâmetro inválido '" + e.parametro + "': " + e.mensagem
}
