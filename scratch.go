package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// ============================================================================
// ARQUIVO: scratch.go
// POR QUE EXISTE: Este arquivo serviu como um "Rascunho" ou "Ambiente de Testes" (Sandbox).
// Durante o desenvolvimento, nossa API estava retornando 0 produtos. 
// Para descobrir o porquê, eu criei este script minúsculo e isolado que testa 
// apenas UMA COISA: "Consigo ler a página 1 e encontrar os links?".
// 
// Como aqui ele usou `http.Get` simples (sem os headers do nosso scraper real),
// ele conseguiu achar os links. Isso foi o que me deu a pista de que o erro
// não estava no seletor CSS do link, mas sim nos Headers da nossa requisição principal!
// (Onde eu descobri que o Accept-Encoding estava quebrando a descompressão do HTML).
// ============================================================================

func main() {
	// 1. Faz uma requisição GET super simples e limpa pra página 1.
	resp, err := http.Get("https://webscraper.io/test-sites/e-commerce/static/computers/laptops?page=1")
	if err != nil {
		log.Fatal(err) // Se a internet cair, mata o script e mostra o erro
	}
	defer resp.Body.Close() // Fecha a porta de rede quando o script acabar
	
	// 2. Transforma o HTML bruto em uma "Árvore" navegável usando o pacote goquery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// 3. Procura no HTML qualquer tag <a class="title"> que esteja dentro de algo <... class="thumbnail">
	// O método Each faz um loop por todos que ele encontrar.
	doc.Find(".thumbnail a.title").Each(func(i int, s *goquery.Selection) {
		
		// 4. Extrai o link de fato (o atributo href="...") do botão clicado
		href, _ := s.Attr("href")
		
		// 5. Imprime na tela do terminal para provar que a extração funciona.
		fmt.Println("Link encontrado na página 1:", href)
	})
	
	fmt.Println("Script de teste de seletores finalizado!")
}
