package repository

import (
	"sync"
	"time"

	"github.com/DhalosNeto/1n8_teste_tecnico/domain/entity"
)

// RepositorioNotebook define o contrato.
type RepositorioNotebook interface {
	ObterTodos() ([]entity.Notebook, bool)
	Salvar(notebooks []entity.Notebook)
	Invalidar()
}

type cacheMemoria struct {
	mu        sync.RWMutex
	dados     []entity.Notebook
	buscadoEm time.Time
	ttl       time.Duration
}

// NovoCacheMemoria é o construtor do cache.
func NovoCacheMemoria(ttl time.Duration) RepositorioNotebook {
	return &cacheMemoria{ttl: ttl}
}

// ObterTodos devolve a lista de notebooks.
func (c *cacheMemoria) ObterTodos() ([]entity.Notebook, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.dados == nil || time.Since(c.buscadoEm) > c.ttl {
		return nil, false
	}
	return c.dados, true
}

// Salvar guarda a lista nova de notebooks.
func (c *cacheMemoria) Salvar(notebooks []entity.Notebook) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.dados = notebooks
	c.buscadoEm = time.Now()
}

// Invalidar apaga os dados do cache
func (c *cacheMemoria) Invalidar() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.dados = nil
	c.buscadoEm = time.Time{}
}
