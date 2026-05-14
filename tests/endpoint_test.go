package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DhalosNeto/1n8_teste_tecnico/application/webmodel"
	"github.com/DhalosNeto/1n8_teste_tecnico/domain/entity"
	"github.com/DhalosNeto/1n8_teste_tecnico/domain/service"
	"github.com/DhalosNeto/1n8_teste_tecnico/routes/controller"
)

type MockRepo struct {
	DadosFalsos []entity.Notebook
}

func (m *MockRepo) ObterTodos() ([]entity.Notebook, bool) {
	return m.DadosFalsos, true
}
func (m *MockRepo) Salvar(n []entity.Notebook) {}
func (m *MockRepo) Invalidar()                 {}

func TestObterLenovo_SucessoEOrdenacao(t *testing.T) {
	dadosSimulados := []entity.Notebook{
		{Nome: "Lenovo Caro", Preco: 2000.0, Avaliacao: 5},
		{Nome: "Dell XPS", Preco: 1500.0, Avaliacao: 4}, // Dell não deve aparecer!
		{Nome: "Lenovo Barato", Preco: 500.0, Avaliacao: 3},
	}

	repoFalso := &MockRepo{DadosFalsos: dadosSimulados}

	servico := service.NovoServicoNotebook(repoFalso, nil)
	controlador := controller.NovoControladorNotebook(servico)

	req := httptest.NewRequest(http.MethodGet, "/api/laptops/lenovo", nil)
	w := httptest.NewRecorder()

	controlador.ObterLenovo(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Esperava status 200 OK, recebeu %d", res.StatusCode)
	}

	var respostaJSON webmodel.RespostaAPINotebooks
	if err := json.NewDecoder(res.Body).Decode(&respostaJSON); err != nil {
		t.Fatalf("Erro ao ler JSON: %v", err)
	}

	if respostaJSON.Total != 2 {
		t.Errorf("Esperava 2 produtos, recebeu %d", respostaJSON.Total)
	}

	if respostaJSON.Dados[0].Nome != "Lenovo Barato" {
		t.Errorf("Esperava 'Lenovo Barato' na primeira posição, recebeu '%s'", respostaJSON.Dados[0].Nome)
	}
}
