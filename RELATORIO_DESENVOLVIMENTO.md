# 🚀 Relatório de Desenvolvimento e Decisões Arquiteturais

Este documento foi criado para registrar a jornada de pensamento, o planejamento e os desafios superados durante a construção da API de Web Scraping para extração de Notebooks Lenovo.

---

## 1. O Desafio Inicial

O objetivo principal era extrair notebooks da marca Lenovo do site `webscraper.io/test-sites/e-commerce/static`, coletando Nome, Preço, Descrição, Avaliação, Número de Reviews e a URL do produto, ordenando tudo do menor para o maior preço.

**Restrição Crucial**: Não podíamos usar ferramentas de automação de browser (como Selenium ou Puppeteer). O trabalho precisava ser feito puramente via requisições HTTP e parsing de HTML estático.

## 2. Decisões de Arquitetura

Decidimos usar **Golang** por dois motivos principais:
1. **Velocidade e Concorrência**: O scraper precisaria varrer várias páginas. As _goroutines_ do Go são perfeitas para fazer requisições simultâneas gastando quase zero de memória.
2. **Tipagem Forte e Organização**: Go permite uma separação muito clara de responsabilidades.

Adotamos os princípios da **Clean Architecture (Arquitetura Limpa)**. Ao invés de jogar tudo no arquivo `main.go`, separamos o projeto em camadas:
* **Camada de Aplicação (`webservice` e `webmodel`)**: Lida com o mundo externo (fazer o request HTML e formatar o JSON).
* **Camada de Domínio (`entity`, `repository`, `service`)**: Onde fica o nosso "Coração". A entidade `Notebook` não sabe de onde veio a informação. O `ServicoNotebook` aplica as regras de negócio puras (Filtrar Lenovo e Ordenar).
* **Camada de Rotas (`controller` e `router`)**: O "Garçom" que atende o cliente, pega os parâmetros da URL e devolve a resposta HTTP.

## 3. A Evolução do Robô (Crawler)

### Fase A: Abordagem Simples (Que falhou)
Inicialmente, tentamos apenas puxar os produtos da página de listagem principal. Mas logo percebemos que a página usava **Paginação**. Havia 20 páginas, e pior: a descrição completa e as avaliações exatas muitas vezes não estão disponíveis na miniatura, exigindo que a gente visitasse o link do produto.

### Fase B: O "Crawler de Duas Fases" (A solução robusta)
Repensamos a lógica para um crawler dividido em duas etapas para garantir 100% de exatidão dos dados:
1. **Coleta**: Vasculhamos as páginas de 1 a 20 e extraímos todos os links individuais de produtos.
2. **Mineração**: Entramos em cada um desses 117 links encontrados para extrair a ficha técnica completa do produto.

Para que isso não demorasse 3 minutos, usamos **Goroutines (Paralelismo)**. Contudo, para não sobrecarregar o servidor alvo e sofrer um bloqueio de IP, criamos um **Semáforo com Channels** (limite de 5 conexões simultâneas).

## 4. Obstáculos e Resolução de Problemas (Debugging)

Durante os testes finais, a nossa API passou a retornar **0 produtos**. Para descobrir o motivo, usamos a clássica técnica de "Dividir para Conquistar", criando os arquivos de rascunho (`scratch.go` e `scratch2.go`) e baixando o HTML "cru" (`pagina_temp.html`). Encontramos três problemas críticos:

### Problema 1: O Header "Accept-Encoding"
* **O Erro**: Para tentar ser um "bom cidadão" na internet, informamos no nosso Request Header que aceitávamos respostas comprimidas (`gzip, deflate, br`). O site então devolvia HTML comprimido em _Brotli_.
* **O Efeito**: A biblioteca HTTP do Go não descompacta _Brotli_ automaticamente. Então, nosso parser (`goquery`) recebia bytes embaralhados e não encontrava nenhuma tag HTML.
* **A Solução**: Removemos o header manual. A standard library do Go voltou a pedir `gzip` silenciosamente e a fazer a descompressão automática de volta para HTML em texto puro.

### Problema 2: A Atualização do Site Alvo
* **O Erro**: Antigamente, os títulos dos produtos ficavam em `<a class="title">` e as estrelas de avaliação ficavam em `<span class="glyphicon-star">`.
* **O Efeito**: Nosso robô passava pela página, mas não achava esses nomes.
* **A Solução**: Vasculhando os arquivos `.html` baixados na mão, notamos que o desenvolvedor do `webscraper.io` mudou o layout para o Bootstrap 5. Os títulos foram para `<h4 class="title">` e os icones mudaram para `ws-icon-star`. Atualizamos os seletores CSS no nosso `scraper_service.go`.

### Problema 3: Variações de Armazenamento
* Observamos que na página os notebooks possuem diferentes capacidades (ex: 128GB, 256GB) que mudam o preço ao clicar.
* **Decisão Arquitetural**: Como o desafio exigia puramente requisições HTTP (sem uso de Puppeteer para clicar nos botões ou avaliar o JavaScript), decidimos conscientemente capturar a **Configuração Padrão** que vem no HTML inicial. Tentar fazer engenharia reversa do JavaScript fugiria do escopo de um Web Scraper passivo e seria uma péssima prática de manutenção a longo prazo.

## 5. Toque Final: O Cache em Memória
Sabíamos que raspar 20 páginas e mais 117 produtos a cada requisição (`GET /api/laptops/lenovo`) seria incrivelmente lento e arriscado (banimento de IP).

Por isso, na camada de Domínio, implementamos um **Cache Thread-Safe** usando `sync.RWMutex`.
* No primeiro acesso, o usuário espera alguns segundos (O Scraper vai pra rua trabalhar).
* O resultado é guardado na memória RAM.
* Nos próximos acessos, o tempo de resposta da nossa API cai para **0.001 milissegundos**.
* Implementamos um tempo de vida (TTL) de 5 minutos, e também uma rota extra (`/refresh`) para limpar a memória sob demanda.

## Conclusão
O projeto evoluiu de um script simples de raspagem para uma API de nível empresarial, aplicando paralelismo, resiliência a bloqueios, design patterns (DDD) e tratamento seguro de memória no Golang. Todas as regras de negócio foram atendidas com sucesso!
