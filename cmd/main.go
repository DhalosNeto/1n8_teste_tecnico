package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/DhalosNeto/1n8_teste_tecnico/application/webservice"
	"github.com/DhalosNeto/1n8_teste_tecnico/domain/repository"
	"github.com/DhalosNeto/1n8_teste_tecnico/domain/service"
	"github.com/DhalosNeto/1n8_teste_tecnico/routes"
	"github.com/DhalosNeto/1n8_teste_tecnico/routes/controller"
)

func main() {
	porta := obterEnv("PORT", "3000")
	ttlCache := converterDuracao(obterEnv("CACHE_TTL", "5m"))

	cache := repository.NovoCacheMemoria(ttlCache)
	scraper := webservice.NovoScraperService()

	servicoNotebook := service.NovoServicoNotebook(cache, scraper)

	controladorNotebook := controller.NovoControladorNotebook(servicoNotebook)

	roteador := http.NewServeMux()
	routes.Registrar(roteador, controladorNotebook)

	endereco := fmt.Sprintf(":%s", porta)
	// Logs de inicialização
	log.Printf("Servidor rodando em http://localhost%s", endereco)
	log.Printf("TTL do cache: %s", ttlCache)
	log.Printf("Endpoints disponíveis:")
	log.Printf("GET http://localhost%s/health", endereco)
	log.Printf("GET http://localhost%s/api/laptops/lenovo", endereco)
	log.Printf("GET http://localhost%s/api/laptops/lenovo?min_price=300&max_price=800&min_rating=3", endereco)
	log.Printf("GET http://localhost%s/api/laptops/lenovo/refresh", endereco)

	servidor := &http.Server{
		Addr:         endereco,
		Handler:      roteador,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := servidor.ListenAndServe(); err != nil {
		log.Fatalf("Erro no servidor: %v", err)
	}
}

func obterEnv(chave, padrao string) string {
	if v := os.Getenv(chave); v != "" {
		return v
	}
	return padrao
}

func converterDuracao(s string) time.Duration {
	duracao, err := time.ParseDuration(s)
	if err != nil {
		log.Printf("CACHE_TTL inválido %q, usando padrão de 5m", s)
		return 5 * time.Minute
	}
	return duracao
}
