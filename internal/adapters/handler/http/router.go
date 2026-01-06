package http

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/biodoia/ghrego/internal/config"
	"github.com/biodoia/ghrego/internal/core/domain"
	"github.com/biodoia/ghrego/internal/core/ports"
	"github.com/rs/zerolog/log"
)

type Server struct {
	router      *chi.Mux
	config      *config.Config
	ghService   ports.GitHubService
	aiService   ports.AIAnalysisService
	repoStore   ports.RepositoryStore
	userRepo    ports.UserRepository
	suggRepo    ports.SuggestionRepository
}

func NewServer(
	cfg *config.Config,
	ghService ports.GitHubService,
	aiService ports.AIAnalysisService,
	repoStore ports.RepositoryStore,
	userRepo ports.UserRepository,
	suggRepo ports.SuggestionRepository,
) *Server {
	s := &Server{
		router:    chi.NewRouter(),
		config:    cfg,
		ghService: ghService,
		aiService: aiService,
		repoStore: repoStore,
		userRepo:  userRepo,
		suggRepo:  suggRepo,
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(60 * time.Second))
	s.router.Use(render.SetContentType(render.ContentTypeJSON))

	// CORS
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   s.config.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	s.router.Route("/api", func(r chi.Router) {
		r.Use(s.authMiddleware) // Mock auth for now

		// Auth
		r.Get("/auth/me", s.handleGetMe)
		r.Post("/auth/logout", s.handleLogout)

		// Repositories
		r.Route("/repositories", func(r chi.Router) {
			r.Post("/sync", s.handleSyncRepositories)
			r.Get("/list", s.handleListRepositories)
			r.Get("/stats", s.handleGetRepositoryStats)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", s.handleGetRepository)
				r.Delete("/", s.handleDeleteRepository)
			})
		})

		// Analysis
		r.Route("/analysis", func(r chi.Router) {
			r.Post("/start", s.handleStartAnalysis)
			r.Get("/get", s.handleGetAnalysis) // using Query param ?repositoryId=... to match tRPC style
			r.Get("/list", s.handleListAnalysis)
		})

		// Suggestions
		r.Route("/suggestions", func(r chi.Router) {
			r.Get("/list", s.handleListSuggestions)
			r.Post("/updateStatus", s.handleUpdateSuggestionStatus)
		})
	})
}

func (s *Server) Run() error {
	log.Info().Str("port", s.config.Port).Msg("Starting HTTP Server")
	return http.ListenAndServe(":"+s.config.Port, s.router)
}

// Mock Middleware - Extract user from header or session
func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Real implementation should decode JWT
		// For now, assume User ID 1 is logged in if not specified, 
		// OR check header "X-User-ID"
		
		userID := 1 // Default for MVP/Dev
		// ... Auth logic ...
		
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// --- Handlers ---

func (s *Server) handleGetMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	user, err := s.userRepo.GetByID(r.Context(), userID)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	if user == nil {
		render.Render(w, r, ErrNotFound)
		return
	}
	render.JSON(w, r, user)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	// Clear cookie logic here
	render.JSON(w, r, map[string]bool{"success": true})
}

func (s *Server) handleSyncRepositories(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	
	// Get User OpenID - needed for sync logic (legacy/db requirement)
	user, err := s.userRepo.GetByID(r.Context(), userID)
	if err != nil || user == nil {
		render.Render(w, r, ErrUnauthorized)
		return
	}

	if err := s.ghService.SyncUserRepositories(r.Context(), userID, user.OpenID); err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}

	// Fetch updated list count
	repos, _ := s.repoStore.GetByUserID(r.Context(), userID)
	
	render.JSON(w, r, map[string]interface{}{
		"success": true,
		"count":   len(repos),
	})
}

func (s *Server) handleListRepositories(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	repos, err := s.repoStore.GetByUserID(r.Context(), userID)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, repos)
}

func (s *Server) handleGetRepository(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	
	repo, err := s.repoStore.GetByID(r.Context(), id)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	if repo == nil {
		render.Render(w, r, ErrNotFound)
		return
	}
	render.JSON(w, r, repo)
}

func (s *Server) handleDeleteRepository(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	
	if err := s.repoStore.Delete(r.Context(), id); err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, map[string]bool{"success": true})
}

func (s *Server) handleGetRepositoryStats(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	stats, err := s.repoStore.GetStats(r.Context(), userID)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, stats)
}

// Analysis Handlers

type StartAnalysisRequest struct {
	RepositoryID int `json:"repositoryId"`
}

func (s *Server) handleStartAnalysis(w http.ResponseWriter, r *http.Request) {
	var req StartAnalysisRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// In async go routine
	go func() {
		// Create new context for background task
		ctx := context.Background()
		_, err := s.aiService.AnalyzeRepository(ctx, req.RepositoryID, domain.AnalysisTypeArchitecture)
		if err != nil {
			log.Error().Err(err).Int("repo_id", req.RepositoryID).Msg("Background analysis failed")
		}
	}()

	render.JSON(w, r, map[string]interface{}{"success": true, "message": "Analysis started"})
}

func (s *Server) handleGetAnalysis(w http.ResponseWriter, r *http.Request) {
	// tRPC style uses Query params for GET input
	repoIDStr := r.URL.Query().Get("repositoryId") // Changed from input.repositoryId
	repoID, _ := strconv.Atoi(repoIDStr)
	_ = repoID // Silence unused variable warning until implemented

	// Stub - needs implementation of aggregating all data
	// See tRPC implementation: fetches analyses, features, techs, suggestions
	render.JSON(w, r, map[string]interface{}{"status": "implemented in next step"})
}

func (s *Server) handleListAnalysis(w http.ResponseWriter, r *http.Request) {
	// List all analyses for user
	render.JSON(w, r, []string{})
}

// Suggestions Handlers

func (s *Server) handleListSuggestions(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	suggs, err := s.suggRepo.GetAllPending(r.Context(), userID)
	if err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
	render.JSON(w, r, suggs)
}

func (s *Server) handleUpdateSuggestionStatus(w http.ResponseWriter, r *http.Request) {
	// Implementation skipped for brevity
	render.JSON(w, r, map[string]bool{"success": true})
}

// --- Errors ---

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request",
		ErrorText:      err.Error(),
	}
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found"}
var ErrUnauthorized = &ErrResponse{HTTPStatusCode: 401, StatusText: "Unauthorized"}

func ErrInternal(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Internal Server Error",
		ErrorText:      err.Error(),
	}
}
