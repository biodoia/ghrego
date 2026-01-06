# Ghrego - GitHub Repo Analyzer (Go Backend)

[![Go Report Card](https://goreportcard.com/badge/github.com/biodoia/ghrego)](https://goreportcard.com/report/github.com/biodoia/ghrego)
[![Coverage Status](https://img.shields.io/badge/coverage-16.5%25-red)](https://github.com/biodoia/ghrego)

Refactoring completo in **Go** del backend di GitHub Repo Analyzer.
Progettato per performance, scalabilità e manutenibilità utilizzando la **Clean Architecture**.

## 🚀 Funzionalità

*   **Sincronizzazione GitHub**: Recupero rapido di repository e metadati.
*   **Analisi AI**: Integrazione con **Google Gemini 1.5** per analisi architetturale e suggerimenti di codice.
*   **API REST**: Interfaccia HTTP moderna e veloce.
*   **Persistenza**: Utilizzo efficiente di PostgreSQL tramite driver nativo `pgx`.

## 🛠 Prerequisiti

*   **Go** 1.22+
*   **PostgreSQL** 15+
*   **API Keys**:
    *   GitHub Personal Access Token
    *   Google Gemini API Key

## ⚙️ Configurazione

Crea un file `.env` nella root del progetto o esporta le variabili d'ambiente:

```bash
PORT=8080
LOG_LEVEL=info
DATABASE_URL="postgres://user:password@localhost:5432/ghrego?sslmode=disable"
API_KEY="tuo-github-token"
GEMINI_API_KEY="tua-gemini-key"
SKIP_BACKEND_CHECK=true
```

> **Nota**: `API_KEY` è usata temporaneamente sia per l'auth interna che come token GitHub di default se non presente nell'utente.

## 🏃‍♂️ Avvio Rapido

1.  **Installa dipendenze**:
    ```bash
    go mod tidy
    ```

2.  **Esegui i Test**:
    ```bash
    go test ./...
    ```

3.  **Compila ed Esegui**:
    ```bash
    go build -o server_bin ./cmd/server
    ./server_bin
    ```

## 🏗 Architettura

Vedi [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) per i dettagli completi su Clean Architecture, Layer e decisioni progettuali.

## 🧪 Testing

Il progetto utilizza:
*   `testify`: Per asserzioni e mocking suite.
*   `pgxmock`: Per simulare PostgreSQL nei test unitari.

```bash
# Esegui test con coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📝 Roadmap

- [x] Core Domain & Models
- [x] PostgreSQL Adapters (`pgx`)
- [x] GitHub Adapter
- [x] AI Adapter (Gemini)
- [x] REST API (`go-chi`)
- [ ] Integrazione completa Frontend React
- [ ] WebSocket per progressi real-time
- [ ] Brain Service (Vector DB logic)
