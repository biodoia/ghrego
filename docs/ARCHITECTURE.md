# Architettura di Ghrego (Go Backend)

Questo documento descrive le scelte architetturali, la struttura del codice e i pattern di design adottati per il refactoring del backend di GitHub Repo Analyzer in Go.

## 🏛 Clean Architecture (Hexagonal)

Il progetto segue rigorosamente i principi della **Clean Architecture** (o Architettura Esagonale). L'obiettivo è separare la logica di business dai dettagli implementativi (Database, API esterne, Framework HTTP).

### I Layer

Il codice è organizzato in cerchi concentrici, dove le dipendenze puntano solo verso l'interno.

#### 1. Core (Il Centro)
Situato in `internal/core`. Non ha **nessuna dipendenza** esterna (niente SQL, niente HTTP, niente librerie di terze parti complesse).
*   **`domain/`**: Contiene le `struct` pure (Entity) che rappresentano i dati (es. `User`, `Repository`).
*   **`ports/`**: Definisce le **Interfacce** (Contratti) che il mondo esterno deve soddisfare.
    *   *Primary Ports (Input)*: Servizi usati dagli handler (es. `GitHubService`).
    *   *Secondary Ports (Output)*: Interfacce per database o API esterne (es. `UserRepository`, `GitHubClient`).

#### 2. Application Logic
Situato in `internal/core/services`.
*   Contiene l'implementazione concreta della logica di business.
*   Implementa le interfacce dei servizi definite nei *Ports*.
*   Orchestra i dati: chiama i Repository, elabora i dati, invoca client esterni.
*   **Esempio**: `SyncUserRepositories` scarica i repo da GitHub (tramite adapter) e li salva su DB (tramite adapter), senza sapere *come* questi funzionino.

#### 3. Adapters (L'Esterno)
Situato in `internal/adapters`. Qui risiedono le implementazioni concrete che "sporcano" le mani con tecnologie specifiche.
*   **`handler/http`**: Layer di presentazione. Usa `go-chi` per gestire routing REST e JSON marshalling.
*   **`storage/postgres`**: Layer di persistenza. Implementa i Repository usando `pgx` e SQL puro.
*   **`github/`**: Client API verso GitHub.
*   **`ai/`**: Client verso Google Gemini.

#### 4. Configuration & Wiring
Situato in `cmd/server` e `internal/config`.
*   Il `main.go` è l'unico punto dove tutti i pezzi vengono assemblati (Dependency Injection).

## 🛠 Decisioni Tecniche

### Database: `pgx` vs GORM
Abbiamo scelto **`pgx/v5`** invece di un ORM come GORM.
*   **Performance**: `pgx` è significativamente più veloce.
*   **Controllo**: SQL esplicito evita query "magiche" N+1 e permette ottimizzazioni fini (es. `COPY FROM` per bulk insert).
*   **Astrazione**: Usiamo interfacce per il pool di connessioni, facilitando il mocking (`pgxmock`).

### Routing: `go-chi`
*   Leggero, idiomatico (compatibile con `net/http` standard).
*   Middleware chain robusta.

### Testing
*   **Unit Tests**: Servizi e Configurazione testati isolatamente.
*   **Mocks**: Generati manualmente (o con `mockery`) in `internal/mocks`.
*   **Integration Tests**: Repository testati contro un mock del database (`pgxmock`) per verificare la correttezza dell'SQL generato.

## 🔄 Flusso dei Dati

Esempio: **Richiesta Analisi Repository**

1.  **HTTP Request**: `POST /api/analysis/start` arriva a `handler/http`.
2.  **Handler**: Valida il JSON e chiama `aiService.AnalyzeRepository`.
3.  **Service**:
    *   Chiama `repoStore.GetByID` (Porta Secondaria) -> `postgres` esegue SELECT.
    *   Chiama `aiClient.AnalyzeRepository` (Porta Secondaria) -> `ai/gemini` chiama Google API.
    *   Riceve il risultato, lo mappa nel Dominio.
    *   Chiama `analysisRepo.Create` per salvare.
4.  **Response**: L'Handler restituisce JSON 200 OK.

## 📂 Struttura Cartelle

```
ghrego/
├── cmd/server/         # Entry point (main.go)
├── internal/
│   ├── core/           # Logica pura
│   │   ├── domain/     # Structs (User, Analysis...)
│   │   ├── ports/      # Interfacce (Service, Repository)
│   │   └── services/   # Implementazione Business Logic
│   ├── adapters/       # Tecnologie concrete
│   │   ├── ai/         # Gemini Client
│   │   ├── github/     # GitHub Client
│   │   ├── handler/    # HTTP Router & Controllers
│   │   └── storage/    # PostgreSQL Implementation
│   ├── config/         # Gestione Env Vars
│   └── mocks/          # Mock objects per testing
└── docs/               # Documentazione
```
