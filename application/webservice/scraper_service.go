package webservice

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/DhalosNeto/1n8_teste_tecnico/domain/entity"
	"github.com/PuerkitoBio/goquery"
)

const (
	urlBase        = "https://webscraper.io/test-sites/e-commerce/static/computers/laptops"
	urlBaseProduto = "https://webscraper.io"
	totalPaginas   = 20
	concorrencia   = 5
)

// Cliente http
var clienteHTTP = &http.Client{
	Timeout: 30 * time.Second,
}

// Struct para os dados
type ScraperService struct{}

// Construtor
func NovoScraperService() *ScraperService {
	return &ScraperService{}
}

func (s *ScraperService) BuscarTodosNotebooks() ([]entity.Notebook, error) {
	// Coletar links de listagem
	links, err := s.coletarLinksDeListagem()
	if err != nil {
		return nil, fmt.Errorf("fase 1 (coleta de links): %w", err)
	}
	// Raspar as páginas
	notebooks, err := s.scraperPaginas(links)
	if err != nil {
		return nil, fmt.Errorf("fase 2 (scraper de produtos): %w", err)
	}

	return notebooks, nil
}

func (s *ScraperService) coletarLinksDeListagem() ([]string, error) {
	semaforo := make(chan struct{}, concorrencia)

	var mu sync.Mutex
	var wg sync.WaitGroup
	var erroColeta error

	// Sobrescreve o link encontrado em vez de duplicar.
	linksSet := make(map[string]struct{})

	for pagina := 1; pagina <= totalPaginas; pagina++ {
		wg.Add(1) //Incrementa o contador de tarefas pendentes.
		semaforo <- struct{}{}

		go func(p int) {
			defer wg.Done()
			defer func() { <-semaforo }()

			url := fmt.Sprintf("%s?page=%d", urlBase, p)
			links, erro := s.extrairLinksDeListagem(url)
			if erro != nil {
				mu.Lock()
				erroColeta = fmt.Errorf("página %d: %w", p, erro)
				mu.Unlock()
				return
			}

			mu.Lock()
			for _, l := range links {
				linksSet[l] = struct{}{}
			}
			mu.Unlock()
		}(pagina)
	}

	wg.Wait()

	if erroColeta != nil {
		return nil, erroColeta
	}

	// Converte o map para um array
	links := make([]string, 0, len(linksSet))
	for link := range linksSet {
		links = append(links, link)
	}
	return links, nil
}

func (s *ScraperService) extrairLinksDeListagem(url string) ([]string, error) {
	doc, err := s.buscarHTML(url)
	if err != nil {
		return nil, err
	}

	var links []string

	doc.Find(".thumbnail a.title").Each(func(_ int, sel *goquery.Selection) {
		// Obtém o atributo href do link.
		if href, existe := sel.Attr("href"); existe {
			href = strings.TrimSpace(href)

			// Se a URL for relativa, adiciona a URL base.
			if !strings.HasPrefix(href, "http") {
				href = urlBaseProduto + href
			}
			links = append(links, href)
		}
	})

	return links, nil
}

func (s *ScraperService) scraperPaginas(urls []string) ([]entity.Notebook, error) {
	// Varrer a lista de URLs que voltaram.
	semaforo := make(chan struct{}, concorrencia)
	var mu sync.Mutex
	var wg sync.WaitGroup
	var erroScraper error

	// Recebe os dados de cada página.
	notebooks := make([]entity.Notebook, 0, len(urls))

	for _, url := range urls {
		wg.Add(1)
		semaforo <- struct{}{}

		go func(u string) {
			defer wg.Done()
			defer func() { <-semaforo }()

			notebook, erro := s.scraperProduto(u)
			if erro != nil {
				mu.Lock()
				erroScraper = fmt.Errorf("%s: %w", u, erro)
				mu.Unlock()
				return
			}

			mu.Lock()
			notebooks = append(notebooks, notebook)
			mu.Unlock()
		}(url)
	}

	wg.Wait()

	if erroScraper != nil {
		return nil, erroScraper
	}

	return notebooks, nil
}

// scraperProduto entra no link do notebook e pega as informações nas tags HTML.
func (s *ScraperService) scraperProduto(url string) (entity.Notebook, error) {
	doc, err := s.buscarHTML(url)
	if err != nil {
		return entity.Notebook{}, err
	}

	thumbnail := doc.Find(".thumbnail")

	nome := strings.TrimSpace(thumbnail.Find("a.title").Text())
	if nome == "" {
		nome = strings.TrimSpace(thumbnail.Find("h4.title").Text())
	}

	// Pega os textos
	textoPreco := strings.TrimSpace(thumbnail.Find("h4.price").Text())
	descricao := strings.TrimSpace(thumbnail.Find("p.description").Text())
	textoAvaliacoes := strings.TrimSpace(thumbnail.Find("p.review-count").Text())

	preco := converterPreco(textoPreco)

	avaliacao := thumbnail.Find("span.ws-icon-star").Length()
	avaliacoes := converterAvaliacoes(textoAvaliacoes)

	// Retorna a entidade limpa.
	return entity.Notebook{
		Nome:       nome,
		Preco:      preco,
		Moeda:      "USD",
		Descricao:  descricao,
		Avaliacao:  avaliacao,
		Avaliacoes: avaliacoes,
		URL:        url,
	}, nil
}

func (s *ScraperService) buscarHTML(url string) (*goquery.Document, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")

	req.Header.Set("Connection", "keep-alive")

	resp, err := clienteHTTP.Do(req) // Aperta o "Enter" pra acessar o site
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // Fecha a conexão

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status inesperado: %d", resp.StatusCode)
	}
	// Transforma o HTML em "Árvore" que podemos pesquisar.
	return goquery.NewDocumentFromReader(resp.Body)
}

func converterPreco(s string) float64 {
	s = strings.ReplaceAll(s, "$", "")
	s = strings.TrimSpace(s)
	valor, _ := strconv.ParseFloat(s, 64)
	return valor
}

func converterAvaliacoes(s string) int {
	partes := strings.Fields(s)
	if len(partes) == 0 {
		return 0
	}
	n, _ := strconv.Atoi(partes[0]) // Converte String para Int
	return n
}
