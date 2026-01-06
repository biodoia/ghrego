package main

import (
	"context"
	"os"

	"github.com/biodoia/ghrego/internal/adapters/ai"
	"github.com/biodoia/ghrego/internal/adapters/github"
	"github.com/biodoia/ghrego/internal/adapters/handler/http"
	"github.com/biodoia/ghrego/internal/adapters/storage/postgres"
	"github.com/biodoia/ghrego/internal/config"
	"github.com/biodoia/ghrego/internal/core/ports"
	"github.com/biodoia/ghrego/internal/core/services"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Setup Logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Load Config
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatal().Err(err).Msg("Invalid configuration")
	}

	// Connect to DB
	db, err := postgres.NewDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	// Initialize Repositories
	userRepo := postgres.NewUserRepository(db)
	repoStore := postgres.NewRepositoryStore(db)
	analysisRepo := postgres.NewAnalysisRepository(db)
	featureRepo := postgres.NewFeatureRepository(db)
	techRepo := postgres.NewTechnologyRepository(db)
	suggestionRepo := postgres.NewSuggestionRepository(db)

	// Initialize Adapters
	ghClient := github.NewClient(cfg.APIKey) // Use APIKey as GitHub Token for now
	
	// Setup Gemini Client
	geminiClient, err := ai.NewGeminiClient(context.Background(), os.Getenv("GEMINI_API_KEY"))
	if err != nil {
		log.Warn().Err(err).Msg("Failed to initialize Gemini Client (AI features disabled)")
	} else {
		defer geminiClient.Close()
	}

	// Initialize Services
	ghService := services.NewGitHubService(ghClient, repoStore, userRepo)
	
	var aiService ports.AIAnalysisService
	if geminiClient != nil {
		aiService = services.NewAIAnalysisService(geminiClient, repoStore, analysisRepo, featureRepo, techRepo, suggestionRepo)
	} else {
		log.Warn().Msg("AI Service not initialized - Using NoOp or failing calls")
		// Ideally pass a NoOp implementation here to avoid nil pointer in Handler
	}

	// Initialize HTTP Server
	server := http.NewServer(cfg, ghService, aiService, repoStore, userRepo, suggestionRepo)
	
	if err := server.Run(); err != nil {
		log.Fatal().Err(err).Msg("Server failed")
	}
}