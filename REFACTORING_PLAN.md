# REFACTORING PLAN - ghrego

> Generato: 2026-01-14 | Framework di riferimento: framegotui

---

## STATO ATTUALE

| Aspetto | Valore | Note |
|---------|--------|------|
| **Linguaggio** | Go | Completo |
| **Struttura** | Hexagonal | Ottima (ports/adapters) |
| **Test** | Presenti | Da espandere |
| **Docs** | Parziale | Solo ARCHITECTURE.md |
| **go.mod** | ERRORE | Versione Go impossibile |

---

## LACUNE CRITICHE

### 1. go.mod - Versione Go Invalida
```go
go 1.25.5  // NON ESISTE - Go attuale massimo è ~1.22
```

**Impatto**: Build failure su sistemi standard
**Fix**: Cambiare a `go 1.21` o `go 1.22`

### 2. Documentazione Incompleta
Confronto con framegotui:
- [ ] README.md aggiornato ✓
- [ ] CHANGELOG.md mancante
- [ ] docs/GETTING_STARTED.md mancante
- [ ] docs/API.md mancante
- [ ] docs/CONTRIBUTING.md mancante

### 3. Test Coverage
- `config_test.go` ✓
- `ai_test.go` ✓
- `github_test.go` ✓
- `repository_test.go` ✓
- `router_test.go` ✓
- [ ] Integration tests mancanti
- [ ] Benchmark tests mancanti

---

## PIANO DI REFACTORING

### Fase 1: Fix Critici (Priorità ALTA)
```bash
# 1. Fix go.mod version
sed -i 's/go 1.25.5/go 1.22/' go.mod

# 2. Rigenera go.sum
go mod tidy
```

### Fase 2: Documentazione (Priorità MEDIA)
1. Creare `CHANGELOG.md`
2. Creare `docs/GETTING_STARTED.md`
3. Creare `docs/API.md` con documentazione endpoints
4. Aggiungere `docs/CONTRIBUTING.md`

### Fase 3: Testing Enhancement (Priorità BASSA)
1. Aggiungere integration tests
2. Creare benchmark tests per operazioni critiche
3. Aggiungere coverage badge nel README

---

## LINTING CHECKLIST

```bash
# Eseguire questi comandi:
go vet ./...
golangci-lint run
go test -race ./...
go test -cover ./...
```

### Configurazione Suggerita (.golangci.yml)
```yaml
linters:
  enable:
    - gofmt
    - govet
    - errcheck
    - staticcheck
    - gosec
    - ineffassign
    - typecheck
```

---

## STRUTTURA TARGET (ispirata a framegotui)

```
ghrego/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── adapters/
│   ├── cache/
│   ├── config/
│   ├── core/
│   │   ├── domain/
│   │   ├── ports/
│   │   └── services/
│   └── mocks/
├── docs/
│   ├── ARCHITECTURE.md ✓
│   ├── API.md
│   ├── GETTING_STARTED.md
│   └── CONTRIBUTING.md
├── CHANGELOG.md
├── README.md ✓
├── go.mod
└── Makefile  # DA AGGIUNGERE
```

---

## PROSSIMI PASSI

1. [ ] Fix immediato go.mod (critico)
2. [ ] Aggiungere Makefile con targets standard
3. [ ] Creare CHANGELOG.md
4. [ ] Espandere documentazione
5. [ ] Aggiungere integration tests

---

## NOTE

ghrego ha una delle migliori architetture tra i submodules. L'architettura hexagonal è ben implementata con chiara separazione tra ports (interfacce) e adapters (implementazioni). Il fix principale è la versione Go nel go.mod.
