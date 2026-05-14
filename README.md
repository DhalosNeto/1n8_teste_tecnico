# Lenovo Laptop Scraper — REST API

API RESTful em Go que extrai todos os notebooks da marca **Lenovo** do site [webscraper.io](https://webscraper.io/test-sites/e-commerce/static/computers/laptops), filtrando e ordenando do mais barato ao mais caro.

---

## Como rodar

### Pré-requisitos
- [Go 1.25+](https://golang.org/dl/)

### Instalação e execução

```bash
# Clone o repositório
git clone https://github.com/DhalosNeto/1n8_teste_tecnico.git
cd 1n8_teste_tecnico

# Instale as dependências
go mod download

# Inicie o servidor (porta padrão: 3000)
go run ./cmd/main.go
```

### Variáveis de ambiente (opcionais)

| Variável    | Padrão | Descrição                                  |
|-------------|--------|--------------------------------------------|
| `PORT`      | `3000` | Porta em que a API irá escutar             |
| `CACHE_TTL` | `5m`   | Tempo de vida do cache (ex: `10m`, `1h`)   |

```bash
PORT=8080 CACHE_TTL=10m go run ./cmd/main.go
```

---

## Endpoints

### `GET /health`
Verifica se a API está no ar.

```bash
curl http://localhost:3000/health
```

```json
{ "status": "ok", "time": "2026-05-12T22:00:00Z" }
```

---

### `GET /api/laptops/lenovo`
Retorna todos os notebooks Lenovo, ordenados do mais barato ao mais caro.

**Query params opcionais:**

| Parâmetro    | Tipo    | Descrição                       |
|--------------|---------|---------------------------------|
| `min_price`  | float   | Preço mínimo em USD             |
| `max_price`  | float   | Preço máximo em USD             |
| `min_rating` | int 1–5 | Avaliação mínima em estrelas    |

```bash
# Todos os Lenovos
curl http://localhost:3000/api/laptops/lenovo

# Com filtros
curl "http://localhost:3000/api/laptops/lenovo?min_price=300&max_price=800&min_rating=3"
```

**Resposta:**
```json
{
  "total": 12,
  "source": "https://webscraper.io/test-sites/e-commerce/static/computers/laptops",
  "scraped_at": "2026-05-12T22:00:00Z",
  "cached": false,
  "data": [
    {
      "name": "Lenovo IdeaPad 100",
      "price": 356.99,
      "currency": "USD",
      "description": "15.6\", Celeron N2840 2.16GHz, 2GB, 250GB, Linux",
      "rating": 3,
      "reviews": 6,
      "url": "https://webscraper.io/test-sites/e-commerce/static/product/42"
    }
  ]
}
```

---

### `GET /api/laptops/lenovo/refresh`
Força um novo scraping, ignorando o cache. Aceita os mesmos query params.

```bash
curl http://localhost:3000/api/laptops/lenovo/refresh
```

---

## Arquitetura

O projeto segue uma **arquitetura em camadas**:

```
cmd/
└── main.go                    # Ponto de entrada

application/
├── webservice/
│   └── scraper_service.go     # HTTP client + parsing HTML com goquery
└── webmodel/
    └── laptop.go              # Structs de resposta da API (JSON tags)

domain/
├── entity/
│   └── laptop.go              # Entidade de domínio
├── repository/
│   └── laptop_repository.go   # Interface + cache em memória
└── service/
    └── laptop_service.go      # Regras de negócio: filtrar, ordenar, cachear

routes/
├── router.go                  # Registro das rotas
└── controller/
    └── laptop_controller.go   # Handlers HTTP 
```

### Fluxo de dados

```
Request HTTP
    → router.go           (roteia para o handler correto)
    → laptop_controller   (extrai e valida query params)
    → laptop_service      (verifica cache, aplica filtros e ordena)
    → laptop_repository   (lê ou invalida o cache em memória)
    → scraper_service     (scraping paralelo das 20 páginas com goquery)
    → Response JSON
```

---

## Decisões técnicas

### Por que Go?
Go oferece **concorrência nativa** com goroutines, compilação rápida, e stdlib robusta para HTTP. O scraping paralelo de 20 páginas é natural com `sync.WaitGroup` e semáforo por canal.

### Por que `net/http` + `goquery`?
- `net/http` (stdlib) oferece controle total sobre headers sem dependências extras.
- `goquery` é a biblioteca padrão para parsing HTML estático.

### Scraping paralelo com semáforo
As 20 páginas são scraped em paralelo com concorrência limitada a 5 goroutines simultâneas via canal com buffer:
```go
sem := make(chan struct{}, 5)
```
Isso evita sobrecarga no servidor.

### Cache em memória thread-safe
O cache usa `sync.RWMutex` para leituras concorrentes sem lock e escrita exclusiva, com TTL configurável. Isso evita scraping repetido a cada request.

### Filtragem por marca
A filtragem usa `strings.Contains(strings.ToLower(name), "lenovo")` para ser robusta a variações de capitalização.

---

## 📦 Dependências

| Pacote | Uso |
|--------|-----|
| [`github.com/PuerkitoBio/goquery`](https://github.com/PuerkitoBio/goquery) | Parsing HTML com seletores CSS |

Todas as outras funcionalidades (HTTP server, concorrência, cache) usam apenas a **stdlib do Go**.
