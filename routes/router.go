package routes

import (
	"net/http"

	"github.com/DhalosNeto/1n8_teste_tecnico/routes/controller"
)

// Registrar monta as Rotas.
func Registrar(roteador *http.ServeMux, controlador *controller.ControladorNotebook) {
	// Health Check
	roteador.HandleFunc("GET /health", controlador.Saude)
	// Atualização manual (Refresh).
	roteador.HandleFunc("GET /api/laptops/lenovo/refresh", controlador.Atualizar)
	// Obter Notebooks Lenovo.
	roteador.HandleFunc("GET /api/laptops/lenovo", controlador.ObterLenovo)
}
