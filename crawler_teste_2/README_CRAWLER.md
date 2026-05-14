# Crawler — Análise e Correções

## Problemas encontrados

### 1. Começando na página 2 e usando htm em vez de html — bug
**Código original**
```go
page := "page-2.htm"
```
**Impacto em produção:** acaba pulando a pagina 1, ou seja, perdendo informações da pagina e está usando htm em vez de html na variavel page isso pode confundir e até mesmo dar erro na execução.

**Correção:** `page = "page-1.html"`

---

### 2. Chamando parseBooks duas vezes
**Código original**
```go
next := getNextPage(doc)
if next == "" {
    parseBooks(doc) 
}
page = next
```
**Impacto em produção:** faz a leitura da ultima pagina duas vezes

**Correção:** removi o bloco `if next == ""`

---

### 3. `Accept-Encoding: zstd`
**Código original**
```go
req.Header.Set("Accept-Encoding", "zstd")
```
**Impacto em produção:** Go já gerencia gzip automaticamente então não teria necessidade de adicionar esse header. Sobre o impacto na produção poderia dar algum conflito retornando algum dado corrompido ou não lido.

**Correção:** cabeçalho removido.

---

### 4. Nenhum tratamento de erros HTTP
**Código original**
```go
resp, err := http.DefaultClient.Do(req)
if err != nil {
	return nil, err
}
```
**Impacto em produção:** faltou verificar o status http de resposta então em caso de erro 429 (rate limit) 404 ou 500 o processo não é parado aceitando a pagina de erro e continuando o processo, perdendo informações de erro.

**Correção:** verifica `resp.StatusCode` e retorna erro para qualquer status não 200.

---

### 5. Sem timeout no HTTP client
**Código original**
```go
resp, err := http.DefaultClient.Do(req)
```
**Impacto em produção:** `http.DefaultClient` não tem timeout. Em produção o
processo nunca termina e fica consumindo recursos.

**Correção:**
```go
var httpClient = &http.Client{Timeout: 30 * time.Second}
```

### 6. Navegação para a próxima página 
**Código original**
```go
if n.FirstChild != nil && n.FirstChild.NextSibling != nil {
    for _, a := range n.FirstChild.NextSibling.Attr { ... }
}
```
**Impacto em produção:** qualquer espaço extra entre o `<li>`
e o `<a>` pode quebrar a navegação, não mudando a página em que está. 

**Correção:** `getNextPage` itera pelos filhos do `<li class="next">` buscando
especificamente um nó `<a>`, independente de posição.

---

### 7. Avaliação não era extraída corretamente
**Código original**
```go
parts := strings.Split(a.Val, " ")
```
**Impacto em produção:** O split não é muito seguro, pois se o html tivesse um espaço extra entre as classes, o split iria retornar mais de 2 elementos, fazendo com que o programa pegasse o valor errado.

**Correção:** O código foi alterado para usar o `strings.Fields` que retorna os elementos separados por espaço.


---

### 8. `strings.Contains` no `hasClass`
**Código original**
```go
strings.Contains(a.Val, class)
```
**Impacto em produção:** `"price_color_old"` seria detectado como
`"price_color"`. Dados extraídos de elementos errados podem corromper o resultado.

**Correção:** Usando strings.Fields() para verificar a classe. Evitamos falsos positivos, pois ele separa as classes por espaço e verifica se a classe é igual.

---
## Bibliotecas utilizadas

| Biblioteca | Motivo |
|---|---|
| `time` | Novo import usei para adicionar o timeout do request |


---

## Como executar

```bash
# Instalar dependência
go mod init crawler
go get golang.org/x/net/html

# Executar
go run crawler.go

```
