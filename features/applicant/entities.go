package applicant

import (
	"context"
	"time"
)

type Stage string

const (
	StageApplied   Stage = "APPLIED"
	StageScreening Stage = "SCREENING"
	StageInterview Stage = "INTERVIEW"
	StageOffer     Stage = "OFFER"
	StageHired     Stage = "HIRED"
	StageRejected  Stage = "REJECTED"
)

var ValidStages = map[Stage]bool{StageApplied: true, StageScreening: true, StageInterview: true, StageOffer: true, StageHired: true, StageRejected: true}

type Applicant struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Position  string    `json:"position"`
	Stage     Stage     `json:"stage"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type CreateInput struct{ Name, Email, Position string }
type ListFilter struct {
	Search, Stage, Position string
	Page, Limit             int
}
type PageResult struct {
	Items              []Applicant
	Page, Limit, Total int64
}
type StageSummary struct {
	Stage Stage `json:"stage"`
	Count int64 `json:"count"`
}
type Summary struct {
	Total   int64          `json:"total"`
	ByStage []StageSummary `json:"by_stage"`
}

type Repository interface {
	Create(context.Context, Applicant) (Applicant, error)
	List(context.Context, ListFilter) (PageResult, error)
	GetByID(context.Context, uint) (Applicant, error)
	UpdateStage(context.Context, uint, Stage) (Applicant, error)
	Summary(context.Context) (Summary, error)
}
type Service interface {
	Create(context.Context, CreateInput) (Applicant, error)
	List(context.Context, ListFilter) (PageResult, error)
	GetByID(context.Context, uint) (Applicant, error)
	UpdateStage(context.Context, uint, Stage) (Applicant, error)
	Summary(context.Context) (Summary, error)
}
