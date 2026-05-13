package service

import (
	"sort"
	"strings"
	"time"

	"github.com/DhalosNeto/1n8_teste_tecnico/application/webservice"
	"github.com/DhalosNeto/1n8_teste_tecnico/domain/entity"
	"github.com/DhalosNeto/1n8_teste_tecnico/domain/repository"
)

const urlFonte = "https://webscraper.io/test-sites/e-commerce/static/computers/laptops"

// OpcoesFiltro define os filtros para consulta.
type OpcoesFiltro struct {
	PrecoMinimo     *float64
	PrecoMaximo     *float64
	AvaliacaoMinima *int
}

// ResultadoNotebook resultado retornado.
type ResultadoNotebook struct {
	Notebooks []entity.Notebook
	RaspadoEm time.Time
	EmCache   bool
	Fonte     string
}

// ServicoNotebook contém as regras de negócio.
type ServicoNotebook struct {
	repo     repository.RepositorioNotebook
	raspador *webservice.ScraperService
}

// NovoServicoNotebook cria uma nova instância de ServicoNotebook.
func NovoServicoNotebook(repo repository.RepositorioNotebook, raspador *webservice.ScraperService) *ServicoNotebook {
	return &ServicoNotebook{repo: repo, raspador: raspador}
}

func (s *ServicoNotebook) ObterLenovo(opcoes OpcoesFiltro) (*ResultadoNotebook, error) {
	notebooks, emCache, raspadoEm, err := s.carregarNotebooks()
	if err != nil {
		return nil, err
	}

	filtrados := filtrarLenovo(notebooks, opcoes)

	ordenarPorPreco(filtrados)

	return &ResultadoNotebook{
		Notebooks: filtrados,
		RaspadoEm: raspadoEm,
		EmCache:   emCache,
		Fonte:     urlFonte,
	}, nil
}

func (s *ServicoNotebook) Atualizar(opcoes OpcoesFiltro) (*ResultadoNotebook, error) {
	s.repo.Invalidar()

	notebooks, _, raspadoEm, err := s.carregarNotebooks()
	if err != nil {
		return nil, err
	}

	filtrados := filtrarLenovo(notebooks, opcoes)
	ordenarPorPreco(filtrados)

	return &ResultadoNotebook{
		Notebooks: filtrados,
		RaspadoEm: raspadoEm,
		EmCache:   false,
		Fonte:     urlFonte,
	}, nil
}

// carregarNotebooks decide se busca do cache ou do raspador.
func (s *ServicoNotebook) carregarNotebooks() ([]entity.Notebook, bool, time.Time, error) {
	if dados, ok := s.repo.ObterTodos(); ok {
		return dados, true, time.Now(), nil
	}

	notebooks, err := s.raspador.BuscarTodosNotebooks()
	if err != nil {
		return nil, false, time.Time{}, err
	}

	s.repo.Salvar(notebooks)

	return notebooks, false, time.Now(), nil
}

// filtrarLenovo aplica os filtros sobre os notebooks Lenovo.
func filtrarLenovo(notebooks []entity.Notebook, opcoes OpcoesFiltro) []entity.Notebook {
	resultado := make([]entity.Notebook, 0)

	for _, n := range notebooks {
		if !strings.Contains(strings.ToLower(n.Nome), "lenovo") {
			continue
		}
		// Verificações opcionais: Ignora notebooks abaixo do preço mínimo ou acima do preço máximo.
		if opcoes.PrecoMinimo != nil && n.Preco < *opcoes.PrecoMinimo {
			continue
		}
		if opcoes.PrecoMaximo != nil && n.Preco > *opcoes.PrecoMaximo {
			continue
		}
		// Ignora notebooks abaixo da avaliação mínima.
		if opcoes.AvaliacaoMinima != nil && n.Avaliacao < *opcoes.AvaliacaoMinima {
			continue
		}
		resultado = append(resultado, n)
	}

	return resultado
}

// ordenarPorPreco modifica a lista original para ficar na ordem exigida.
func ordenarPorPreco(notebooks []entity.Notebook) {
	sort.Slice(notebooks, func(i, j int) bool {
		return notebooks[i].Preco < notebooks[j].Preco
	})
}
