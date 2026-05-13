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

// MockRepo simula nosso cache em memória.
// Ele sempre retorna que "tem dados no cache" (true).
type MockRepo struct {
	DadosFalsos []entity.Notebook
}

func (m *MockRepo) ObterTodos() ([]entity.Notebook, bool) {
	return m.DadosFalsos, true // True = finge que o cache está válido!
}
func (m *MockRepo) Salvar(n []entity.Notebook) {}
func (m *MockRepo) Invalidar()                 {}

// ----------------------------------------------------------------------------
// TESTES
// ----------------------------------------------------------------------------

func TestObterLenovo_SucessoEOrdenacao(t *testing.T) {
	// 1. Preparação (Arrange)
	// Vamos criar alguns dados misturados (marcas diferentes e preços fora de ordem)
	dadosSimulados := []entity.Notebook{
		{Nome: "Lenovo Caro", Preco: 2000.0, Avaliacao: 5},
		{Nome: "Dell XPS", Preco: 1500.0, Avaliacao: 4}, // Dell não deve aparecer!
		{Nome: "Lenovo Barato", Preco: 500.0, Avaliacao: 3},
	}

	repoFalso := &MockRepo{DadosFalsos: dadosSimulados}
	// Passamos nil no scraper, pois o MockRepo vai retornar true e o scraper não será chamado!
	servico := service.NovoServicoNotebook(repoFalso, nil)
	controlador := controller.NovoControladorNotebook(servico)

	// Criamos uma requisição e um "Gravador de Resposta" (Recorder) nativos do Go
	req := httptest.NewRequest(http.MethodGet, "/api/laptops/lenovo", nil)
	w := httptest.NewRecorder()

	// 2. Execução (Act)
	controlador.ObterLenovo(w, req)

	// 3. Verificação (Assert)
	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Esperava status 200 OK, recebeu %d", res.StatusCode)
	}

	var respostaJSON webmodel.RespostaAPINotebooks
	if err := json.NewDecoder(res.Body).Decode(&respostaJSON); err != nil {
		t.Fatalf("Erro ao ler JSON: %v", err)
	}

	// O Dell deve ser filtrado. Restam apenas 2 Lenovos.
	if respostaJSON.Total != 2 {
		t.Errorf("Esperava 2 produtos, recebeu %d", respostaJSON.Total)
	}

	// Eles devem vir ordenados do menor para o maior preço.
	// O Lenovo Barato ($500) deve ser o índice 0.
	if respostaJSON.Dados[0].Nome != "Lenovo Barato" {
		t.Errorf("Esperava 'Lenovo Barato' na primeira posição, recebeu '%s'", respostaJSON.Dados[0].Nome)
	}
}
