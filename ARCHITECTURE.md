# Financial Insights Service: Arquitetura e Fluxo de Dados

Este documento descreve detalhadamente a arquitetura atual do **Financial Insights Service**, com foco nas escolhas de design, padrões estruturais adotados em Go (Golang) e o fluxo do ciclo de vida da requisição desde o client HTTP até a inteligência artificial ou banco de dados.

O projeto foi desenhado utilizando princípios de **Hexagonal Architecture** (Ports and Adapters) aliados à simplicidade idiomática nativa de Go. O foco é manutenibilidade, testabilidade e delimitação estrita de fronteiras (*Boundary Segregation*).

---

## 1. Visão Geral das Camadas (Layers)

A aplicação é fragmentada usando a nomenclatura padrão da comunidade Go, focada no isolamento das regras de negócio (*Core Domain*) em repúdio ao acoplamento a frameworks (*Framework Agnostic*).

### `internal/domain` (Entidades Centrais)
A camada mais purista e interna do sistema. 
- Contém os tipos e as *structs* anêmicas (ou ricas) que regem o negócio (`Transaction`, `FinancialSummary`, `Insight`).
- **Nenhuma outra diretório ou pacote** dita regras sobre essa área. Pelo contrário, todas dependem dela. 

### `internal/transaction` & `internal/insight` (Serviços / Casos de Uso)
Módulos verticais organizados por *Feature/Bounded Context*. Dentro de cada módulo, há a subdivisão tática:
- **`service.go`**: O "Coração". Declara a interface `Repository` que necessita (*Consumer-Side Interface Pattern*). Orquestra as regras do core (ex: Validações de domínio `if t.Amount <= 0`, injeção de ID, processamento).
- **`handler.go`**: A porta de entrada HTTP primária (*Primary Adapter/Driver Port*). É inteiramente focada em fazer de-serialização do JSON (`json.NewDecoder`), delegação cega da transação hidratada ao Serviço e serializar de volta o status code.
- **`postgres_repository.go`**: O adaptador da infraestrutura primário (*Secondary Adapter/Driven Port*). Realiza conversão pragmática de objetos de Domínio em Queries SQL nativas (`jackc/pgx/v5`).

### `internal/ai` & `internal/database` (Infraestrutura Compartilhada)
Serviços perféricos (*Secondary Adapters*).
- **`openai_client.go`**: Implementa a interface exigida pelo Serviço de Insights (`AIClient`). Gerencia a inicialização da biblioteca `sashabaranov/go-openai`, a blindagem das *structs* externas de chat completion e a criação do prompt unificado. 
- **`postgres.go`**: Cria e entrega o pool flexível de conexões (`pgxpool`) isolando configurações de infra (timeouts, limites de sessão).

### `cmd/api/main.go` (Composition Root)
A zona de convergência controlada. É o único arquivo no projeto que sabe quem todos são e os conecta:
- Inicia os **Ambientes Globais (`internal/config`)** e o Pool de Conexões.
- Orquestra a injeção de dependências: Injeta o DB no Repo, o Repo e a Camada de IA no Service, e o Service no Handler.
- Liga o Mux do Servidor.
- Mantém o bloqueio assíncrono do servidor Web (`http.Server`) gerido por uma malha de **Graceful Shutdown** via Context e canais do OS (`os.Signal`).

---

## 2. Diagrama Arquitetural de Dependências

Abaixo, a representação hierárquica usando um Diagrama de Fluxo (as setas indicam injecão de conhecimento):

```mermaid
graph TD
    Client[Client (Web/Mobile)] -->|HTTP POST/GET| Handlers

    subgraph "Application Core (Immutable)"
        Handlers[HTTP Handlers] -->|Calls DTO/Domain| Services
        Services[Business Logic / Services] -->|Uses| Domain[Domain Entities]
        Services -.->|Require Interface| Ports[Ports Interfaces]
    end

    subgraph "Infrastructure Adapters"
        Ports -.->|Implemented By| Repositories[PostgreSQL Repositories]
        Ports -.->|Implemented By| AI[OpenAI Client]
        Repositories -.->|SQL| Database[(PostgreSQL)]
        AI -.->|HTTP| OpenAI[OpenAI GPT API]
    end

    Main[main.go / Composition Root] -->|Orchestrates| Handlers
    Main -->|Injects| Repositories
    Main -->|Injects| AI
    Main -->|Assembles| Services
```

> **Nota:** Repare que na injeção de dependência de Go, as **Interfaces vivem na camada de quem as consome** (nos Services), garantindo que a regra de negócio não seja acoplada a quem provê a implementação concreta (Repos ou AI).

---

## 3. Fluxo Base: Transaction Processing

1. **Recepção (`handler.go`)**: O Mux capta o `POST /transactions` e aloca para `txHandler.CreateTransaction`.
2. **De-serialização**: O Request payload é mapeado fisicamente na entity `domain.Transaction`.
3. **Validação de Domínio (`service.go`)**: A função entra no `ProcessTransaction`. Ali, *early-returns* garantem a pureza: avalia se as moedas não são negativas, se o `customer_id` é válido e formata o UUIDv4.
4. **Persistência (`postgres_repository.go`)**: Passa silenciosamente a entity preenchida ao banco. Usando `Exec` nativo, a persistência final consolida.

## 4. Fluxo Avançado: AI Financial Generation
Este fluxo incorpora operações pesadas distribuídas de forma paralela ao longo das *Boundaries*:

1. **Engatilho (`handler.go`)**: Requisição de `POST /customers/{customer_id}/generate-insight`, convertendo o URL Pattern.
2. **Heavy-Lifting Base (`postgres_repository.go` -> `GetSummaryByCustomer`)**:
   - O `insight.Service` pede ao Banco não as transações brutas, mas o **agregado**. A lógica SQL nativa é usada ativamente como cérebro de soma `COALESCE(SUM(amount)...)`. O retorno é um rápido `FinancialSummary`.
3. **Consumo LLM (`openai_client.go` -> `GenerateFinancialInsight`)**:
   - O mesmo Service dispara o `FinancialSummary` para o Adapter de Custo/AI.
   - O Adaptador da OpenAI restringe tokens limitandos MaxTokens a `150`, projeta a *Persona* financeira rígida (Prompt Engineering) e despacha a request sincrona para a OpenAI.
4. **Resumo & Registro**: 
   - A *string* é extraída da *choices list*, transformada em uma Entity `domain.Insight` validada, e guardada no banco de histórico (Postgres) pelo repositório.
5. **Encerramento HTTP**: O Handler cospe o DTO encapsulando o conselho recém gerado (HTTP 200 OK).

## Conclusão de Design
Esse modelo arquitetônico eleva a testabilidade do projeto. Evita a proliferação desenfreada de *mocks* sujos de frameworks e concentra a inteligência do desenvolvedor inteiramente na estrutura do objeto `internal/domain` e a coesão das interfaces no `service.go`.
