package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Enums
type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

type AnalysisType string

const (
	AnalysisTypeArchitecture AnalysisType = "architecture"
	AnalysisTypeFeatures     AnalysisType = "features"
	AnalysisTypeDependencies AnalysisType = "dependencies"
	AnalysisTypeQuality      AnalysisType = "quality"
	AnalysisTypePatterns     AnalysisType = "patterns"
	AnalysisTypeSuggestions  AnalysisType = "suggestions"
)

type AnalysisStatus string

const (
	AnalysisStatusPending    AnalysisStatus = "pending"
	AnalysisStatusProcessing AnalysisStatus = "processing"
	AnalysisStatusCompleted  AnalysisStatus = "completed"
	AnalysisStatusFailed     AnalysisStatus = "failed"
)

type RelationType string

const (
	RelationTypeSimilar            RelationType = "similar"
	RelationTypeContinuation       RelationType = "continuation"
	RelationTypeSharedFeatures     RelationType = "shared_features"
	RelationTypeSharedDependencies RelationType = "shared_dependencies"
	RelationTypeRefactoredFrom     RelationType = "refactored_from"
)

type SuggestionType string

const (
	SuggestionTypeMergeFeatures    SuggestionType = "merge_features"
	SuggestionTypeAddFeature       SuggestionType = "add_feature"
	SuggestionTypeRefactor         SuggestionType = "refactor"
	SuggestionTypeBestPractice     SuggestionType = "best_practice"
	SuggestionTypeConsolidate      SuggestionType = "consolidate"
	SuggestionTypeUpdateDependency SuggestionType = "update_dependency"
)

type SuggestionPriority string

const (
	SuggestionPriorityLow      SuggestionPriority = "low"
	SuggestionPriorityMedium   SuggestionPriority = "medium"
	SuggestionPriorityHigh     SuggestionPriority = "high"
	SuggestionPriorityCritical SuggestionPriority = "critical"
)

type SuggestionStatus string

const (
	SuggestionStatusPending  SuggestionStatus = "pending"
	SuggestionStatusAccepted SuggestionStatus = "accepted"
	SuggestionStatusRejected SuggestionStatus = "rejected"
	SuggestionStatusApplied  SuggestionStatus = "applied"
)

type TechnologyType string

const (
	TechnologyTypeLanguage      TechnologyType = "language"
	TechnologyTypeFramework     TechnologyType = "framework"
	TechnologyTypeLibrary       TechnologyType = "library"
	TechnologyTypeTool          TechnologyType = "tool"
	TechnologyTypeDatabase      TechnologyType = "database"
	TechnologyTypePlatform      TechnologyType = "platform"
)

// User represents the users table
type User struct {
	ID           int            `json:"id" db:"id"`
	OpenID       string         `json:"openId" db:"openId"`
	Name         sql.NullString `json:"name" db:"name"`
	Email        sql.NullString `json:"email" db:"email"`
	LoginMethod  sql.NullString `json:"loginMethod" db:"loginMethod"`
	Role         UserRole       `json:"role" db:"role"`
	GithubUsername sql.NullString `json:"githubUsername" db:"githubUsername"`
	GithubID     sql.NullString `json:"githubId" db:"githubId"`
	CreatedAt    time.Time      `json:"createdAt" db:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt" db:"updatedAt"`
	LastSignedIn time.Time      `json:"lastSignedIn" db:"lastSignedIn"`
}

// Repository represents the repositories table
type Repository struct {
	ID            int            `json:"id" db:"id"`
	UserID        int            `json:"userId" db:"userId"`
	GithubID      string         `json:"githubId" db:"githubId"`
	Name          string         `json:"name" db:"name"`
	FullName      string         `json:"fullName" db:"fullName"`
	Description   sql.NullString `json:"description" db:"description"`
	URL           string         `json:"url" db:"url"`
	Language      sql.NullString `json:"language" db:"language"`
	IsPrivate     bool           `json:"isPrivate" db:"isPrivate"`
	Stars         int            `json:"stars" db:"stars"`
	Forks         int            `json:"forks" db:"forks"`
	Size          int            `json:"size" db:"size"`
	DefaultBranch string         `json:"defaultBranch" db:"defaultBranch"`
	LastCommitAt  sql.NullTime   `json:"lastCommitAt" db:"lastCommitAt"`
	LastSyncAt    sql.NullTime   `json:"lastSyncAt" db:"lastSyncAt"`
	CreatedAt     time.Time      `json:"createdAt" db:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt" db:"updatedAt"`
}

// Analysis represents the analyses table
type Analysis struct {
	ID           int            `json:"id" db:"id"`
	RepositoryID int            `json:"repositoryId" db:"repositoryId"`
	AnalysisType AnalysisType   `json:"analysisType" db:"analysisType"`
	Status       AnalysisStatus `json:"status" db:"status"`
	Result       sql.NullString `json:"result" db:"result"`
	Summary      sql.NullString `json:"summary" db:"summary"`
	Score        sql.NullInt32  `json:"score" db:"score"`
	ErrorMessage sql.NullString `json:"errorMessage" db:"errorMessage"`
	CreatedAt    time.Time      `json:"createdAt" db:"createdAt"`
	CompletedAt  sql.NullTime   `json:"completedAt" db:"completedAt"`
}

// Feature represents the features table
type Feature struct {
	ID           int            `json:"id" db:"id"`
	RepositoryID int            `json:"repositoryId" db:"repositoryId"`
	Name         string         `json:"name" db:"name"`
	Description  sql.NullString `json:"description" db:"description"`
	Category     sql.NullString `json:"category" db:"category"`
	FilePaths    sql.NullString `json:"filePaths" db:"filePaths"` // JSON encoded string likely
	CodeSnippet  sql.NullString `json:"codeSnippet" db:"codeSnippet"`
	Confidence   int            `json:"confidence" db:"confidence"`
	CreatedAt    time.Time      `json:"createdAt" db:"createdAt"`
}

// RepositoryRelation represents the repositoryRelations table
type RepositoryRelation struct {
	ID                 int            `json:"id" db:"id"`
	SourceRepositoryID int            `json:"sourceRepositoryId" db:"sourceRepositoryId"`
	TargetRepositoryID int            `json:"targetRepositoryId" db:"targetRepositoryId"`
	RelationType       RelationType   `json:"relationType" db:"relationType"`
	Similarity         int            `json:"similarity" db:"similarity"`
	Description        sql.NullString `json:"description" db:"description"`
	CreatedAt          time.Time      `json:"createdAt" db:"createdAt"`
}

// Suggestion represents the suggestions table
type Suggestion struct {
	ID                 int                `json:"id" db:"id"`
	RepositoryID       int                `json:"repositoryId" db:"repositoryId"`
	SuggestionType     SuggestionType     `json:"suggestionType" db:"suggestionType"`
	Title              string             `json:"title" db:"title"`
	Description        string             `json:"description" db:"description"`
	SourceRepositoryID sql.NullInt32      `json:"sourceRepositoryId" db:"sourceRepositoryId"`
	Priority           SuggestionPriority `json:"priority" db:"priority"`
	Status             SuggestionStatus   `json:"status" db:"status"`
	CreatedAt          time.Time          `json:"createdAt" db:"createdAt"`
	UpdatedAt          time.Time          `json:"updatedAt" db:"updatedAt"`
}

// Technology represents the technologies table
type Technology struct {
	ID             int            `json:"id" db:"id"`
	RepositoryID   int            `json:"repositoryId" db:"repositoryId"`
	Name           string         `json:"name" db:"name"`
	Version        sql.NullString `json:"version" db:"version"`
	Type           TechnologyType `json:"type" db:"type"`
	PackageManager sql.NullString `json:"packageManager" db:"packageManager"`
	CreatedAt      time.Time      `json:"createdAt" db:"createdAt"`
}

// UnificationOperation represents the unificationOperations table
type UnificationOperation struct {
	ID                   int            `json:"id" db:"id"`
	UserID               int            `json:"userId" db:"userId"`
	OperationID          uuid.UUID      `json:"operationId" db:"operationId"`
	SourceRepositoryIDs  string         `json:"sourceRepositoryIds" db:"sourceRepositoryIds"` // JSON encoded
	TargetRepositoryName string         `json:"targetRepositoryName" db:"targetRepositoryName"`
	TargetRepositoryURL  sql.NullString `json:"targetRepositoryUrl" db:"targetRepositoryUrl"`
	Visibility           string         `json:"visibility" db:"visibility"`
	Status               string         `json:"status" db:"status"`
	Progress             int            `json:"progress" db:"progress"`
	CurrentStep          sql.NullString `json:"currentStep" db:"currentStep"`
	FilesProcessed       int            `json:"filesProcessed" db:"filesProcessed"`
	TotalFiles           int            `json:"totalFiles" db:"totalFiles"`
	Errors               sql.NullString `json:"errors" db:"errors"` // JSON encoded
	CreatedAt            time.Time      `json:"createdAt" db:"createdAt"`
	UpdatedAt            time.Time      `json:"updatedAt" db:"updatedAt"`
	CompletedAt          sql.NullTime   `json:"completedAt" db:"completedAt"`
}
