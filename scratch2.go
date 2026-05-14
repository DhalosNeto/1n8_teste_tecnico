package main

import (
	"fmt"
	"log"
	
	"github.com/DhalosNeto/1n8_teste_tecnico/application/webservice"
)

// ============================================================================
// ARQUIVO: scratch2.go
// POR QUE EXISTE: Este foi o nosso segundo "Rascunho de Teste".
// Diferente do scratch.go que testava apenas se a internet e a biblioteca goquery
// estavam lendo a página 1 direito, este arquivo testava o nosso MOTOR COMPLETO
// (ScraperService) ignorando o cache, o controlador e o roteador web.
// 
// Usar scripts assim (chamados de scripts de integração ou rascunhos)
// é uma técnica excelente de "Dividir para Conquistar". Quando a API retornava 0,
// nós podíamos ter um bug em 3 lugares diferentes: no Router, no Service(Filtros)
// ou no Crawler. Rodando este script direto no Crawler, isolamos o problema
// e vimos que o próprio Crawler estava devolvendo 0 itens.
// ============================================================================

func main() {
	// 1. Instancia o nosso Motor Robô (Scraper) que acabamos de programar.
	scraper := webservice.NovoScraperService()
	
	// 2. Manda ele puxar a cordinha e executar as duas fases inteiras.
	// Ocultamente, ele vai varrer as 20 páginas e entrar nos 117 produtos.
	notebooks, err := scraper.BuscarTodosNotebooks()
	if err != nil {
		log.Fatal(err) // Se a internet cair no meio ou o site bloquear, avisa aqui.
	}
	
	// 3. Imprime linha por linha o nome de tudo que ele conseguiu extrair.
	// Isso me ajudou a confirmar se o site tinha bloqueado a gente ou se
	// os campos estavam vindo vazios.
	for _, n := range notebooks {
		fmt.Printf("Nome: %q\n", n.Nome)
	}
	
	// 4. Mostra quantos ele puxou no total (Esperávamos ~117 produtos).
	fmt.Printf("Total extraído pela máquina de raspagem: %d\n", len(notebooks))
}
