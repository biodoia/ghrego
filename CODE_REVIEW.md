# CODE REVIEW - ghrego

> Generato: 2026-01-14 | Reviewer: Claude Code

---

## EXECUTIVE SUMMARY

| Metrica | Risultato |
|---------|-----------|
| **Build Status** | ✅ PASSA |
| **go vet** | ✅ PASSA |
| **Architettura** | ✅ Hexagonal (Ports/Adapters) |
| **Test Coverage** | ⚠️ Presente ma incompleto |
| **Sicurezza** | ⚠️ Mock auth, da implementare |
| **Severity** | 🟡 MEDIO |

ghrego è il **submodule più maturo** dell'ecosistema autoschei. Compila, ha architettura pulita, e segue best practices Go.

---

## RISULTATI LINTING

### go vet
```
✅ NESSUN ERRORE
```

### go build
```
✅ BUILD SUCCESSFUL
```

---

## ANALISI DETTAGLIATA

### 1. ARCHITETTURA ✅

**Struttura Hexagonal ben implementata:**

```
internal/
├── adapters/           # Implementations
│   ├── handler/http/   # HTTP handlers (primary adapter)
│   ├── storage/postgres/# Database (secondary adapter)
│   ├── ai/             # AI client (secondary adapter)
│   └── github/         # GitHub client (secondary adapter)
├── core/
│   ├── domain/         # Business entities
│   ├── ports/          # Interfaces (contracts)
│   └── services/       # Business logic
├── cache/              # Redis cache
└── config/             # Configuration
```

**Punti di forza:**
- Chiara separazione tra core business logic e adapters
- Interfaces definite in `ports/` → facile testing e swap implementazioni
- Domain models isolati dalla persistenza

### 2. PROBLEMI TROVATI

#### 2.1 Context Key Anti-pattern 🟡
**File:** `internal/adapters/handler/http/router.go:116`
```go
ctx := context.WithValue(r.Context(), "user_id", userID)
```

**Problema:** Usare stringa come context key può causare collisioni.

**Best Practice:**
```go
type contextKey string
const userIDKey contextKey = "user_id"
ctx := context.WithValue(r.Context(), userIDKey, userID)
```

#### 2.2 Mock Authentication 🟠
**File:** `internal/adapters/handler/http/router.go:107-119`
```go
func (s *Server) authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        userID := 1 // Default for MVP/Dev  ← HARDCODED!
        // ... Auth logic ...
```

**Problema:** Authentication hardcoded - qualsiasi richiesta è autenticata come user 1.

**Rischio:** In produzione, bypass totale dell'autenticazione.

#### 2.3 Error Handling Parziale 🟡
**File:** `internal/core/services/github.go:62-66`
```go
_, err := s.repoStore.Upsert(ctx, repo)
if err != nil {
    log.Error().Err(err).Str("repo", repo.Name).Msg("Failed to upsert repository")
    // Continue with others  ← ERRORE SILENZIOSO
}
```

**Osservazione:** Gli errori vengono loggati ma non propagati. OK per resilienza, ma il chiamante non sa quanti repo hanno fallito.

**Suggerimento:** Restituire conteggio errori o lista di fallimenti.

#### 2.4 TODO Comments in Production Code 🟢
**File:** `internal/core/services/github.go:28-40`
```go
// TO FIX: Pass token or username logic.
```

**Osservazione:** Ci sono diversi commenti TODO che indicano lavoro incompleto. Non è un bug, ma indica features incomplete.

#### 2.5 go.mod Version 🔴
**File:** `go.mod:3`
```go
go 1.25.5  // VERSIONE IMPOSSIBILE
```

**Problema:** Go 1.25.5 non esiste. La versione stabile più recente è ~1.22.

**Fix:** Cambiare a `go 1.21` o `go 1.22`

---

## TEST ANALYSIS

### File di Test Presenti
- `internal/config/config_test.go` ✅
- `internal/core/services/ai_test.go` ✅
- `internal/core/services/github_test.go` ✅
- `internal/adapters/storage/postgres/repository_test.go` ✅
- `internal/adapters/handler/http/router_test.go` ✅

### Mancanze
- [ ] Integration tests end-to-end
- [ ] Benchmark tests per operazioni critiche
- [ ] Test coverage report

---

## DEPENDENCY ANALYSIS

### Dipendenze Principali
| Package | Versione | Note |
|---------|----------|------|
| go-chi/chi | v5.2.3 | ✅ Recente |
| jackc/pgx | v5.8.0 | ✅ Recente |
| rs/zerolog | v1.34.0 | ✅ Recente |
| google/generative-ai-go | v0.20.1 | ✅ Recente |
| redis/go-redis | v9.17.2 | ✅ Recente |

**Nota:** Dipendenze sono aggiornate. Nessuna vulnerabilità nota segnalata da GitHub.

---

## SECURITY CHECKLIST

| Check | Status | Note |
|-------|--------|------|
| SQL Injection | ✅ | Usa pgx con parametrizzazione |
| Auth Implementation | 🔴 | Mock - IMPLEMENTARE |
| CORS Configuration | ✅ | Configurato |
| Rate Limiting | ⚠️ | Non implementato |
| Input Validation | ⚠️ | Parziale |
| Secrets Management | ✅ | Usa env vars |

---

## RACCOMANDAZIONI

### Priorità ALTA
1. **Implementare autenticazione reale** (JWT o OAuth2)
2. **Fixare go.mod version** (`go 1.25.5` → `go 1.22`)

### Priorità MEDIA
3. Usare typed context keys invece di stringhe
4. Aggiungere rate limiting
5. Migliorare error handling con conteggio fallimenti

### Priorità BASSA
6. Risolvere TODO comments
7. Aggiungere integration tests
8. Aggiungere benchmarks

---

## METRICHE CODICE

| Metrica | Valore |
|---------|--------|
| File Go | ~25 |
| Lines of Code | ~2000 |
| Test Files | 5 |
| Packages | 10 |
| External Dependencies | 12 |

---

## CONCLUSIONE

ghrego è un progetto **ben strutturato** con architettura hexagonal pulita. I problemi principali sono:
1. Auth mock (da implementare per produzione)
2. go.mod con versione impossibile

Una volta fixati questi, il progetto è production-ready per un MVP.

**Rating: 7/10** ⭐⭐⭐⭐⭐⭐⭐☆☆☆
