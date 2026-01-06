package domain

import (
	"time"
)

// Helper functions to create sql.Null* types
// ... existing helpers ...

// RepositoryAnalysisResponse needed for AI Client interface
// Moved/Copied here because ports depends on domain, and AI Client needs to return this struct.
// It was defined in adapters/ai/gemini.go but needs to be accessible by ports.
type RepositoryAnalysisResponse struct {
	Architecture string `json:"architecture"`
	Features     []struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Category    string   `json:"category"`
		Confidence  int      `json:"confidence"`
		FilePaths   []string `json:"filePaths"`
	} `json:"features"`
	Technologies []struct {
		Name    string `json:"name"`
		Type    string `json:"type"`
		Version string `json:"version"`
	} `json:"technologies"`
	Patterns []string `json:"patterns"`
	Quality  struct {
		Score     int      `json:"score"`
		Issues    []string `json:"issues"`
		Strengths []string `json:"strengths"`
	} `json:"quality"`
	Suggestions []struct {
		Type        string `json:"type"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Priority    string `json:"priority"`
	} `json:"suggestions"`
}

func (r *RepositoryAnalysisResponse) ToDomain(repoID int) (*Analysis, []Feature, []Technology, []Suggestion) {
	// ... Logic moved from gemini.go ...
	// Since the logic is domain mapping, it belongs here anyway.
	
	// Analysis
	analysis := &Analysis{
		RepositoryID: repoID,
		AnalysisType: AnalysisTypeArchitecture,
		Status:       AnalysisStatusCompleted,
		Result:       SQLNullString(""), // Placeholder, usually JSON marshalling
		Summary:      SQLNullString(r.Architecture),
		Score:        SQLNullInt32(r.Quality.Score),
		CreatedAt:    time.Now(),
	}

	// Features
	var features []Feature
	for _, f := range r.Features {
		features = append(features, Feature{
			RepositoryID: repoID,
			Name:         f.Name,
			Description:  SQLNullString(f.Description),
			Category:     SQLNullString(f.Category),
			Confidence:   f.Confidence,
			CreatedAt:    time.Now(),
		})
	}

	// Technologies
	var techs []Technology
	for _, t := range r.Technologies {
		techs = append(techs, Technology{
			RepositoryID:   repoID,
			Name:           t.Name,
			Version:        SQLNullString(t.Version),
			Type:           TechnologyType(t.Type), 
			CreatedAt:      time.Now(),
		})
	}

	// Suggestions
	var suggestions []Suggestion
	for _, s := range r.Suggestions {
		suggestions = append(suggestions, Suggestion{
			RepositoryID:   repoID,
			SuggestionType: SuggestionType(s.Type),
			Title:          s.Title,
			Description:    s.Description,
			Priority:       SuggestionPriority(s.Priority),
			Status:         SuggestionStatusPending,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		})
	}

	return analysis, features, techs, suggestions
}
