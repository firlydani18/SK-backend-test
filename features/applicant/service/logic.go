package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	core "firlydev/backend-test/features/applicant"
)

var (
	ErrValidation     = errors.New("validation failed")
	ErrDuplicateEmail = errors.New("email already exists")
)

type Service struct{ repo core.Repository }

func New(repo core.Repository) *Service { return &Service{repo: repo} }

func (s *Service) Create(ctx context.Context, input core.CreateInput) (core.Applicant, error) {
	input.Name, input.Email, input.Position = strings.TrimSpace(input.Name), strings.ToLower(strings.TrimSpace(input.Email)), strings.TrimSpace(input.Position)
	if input.Name == "" {
		return core.Applicant{}, fmt.Errorf("%w: name is required", ErrValidation)
	}
	if len(input.Name) > 120 {
		return core.Applicant{}, fmt.Errorf("%w: name must be at most 120 characters", ErrValidation)
	}
	if input.Email == "" {
		return core.Applicant{}, fmt.Errorf("%w: email is required", ErrValidation)
	}
	if len(input.Email) > 160 {
		return core.Applicant{}, fmt.Errorf("%w: email must be at most 160 characters", ErrValidation)
	}
	if _, err := mail.ParseAddress(input.Email); err != nil {
		return core.Applicant{}, fmt.Errorf("%w: email must be valid", ErrValidation)
	}
	if input.Position == "" {
		return core.Applicant{}, fmt.Errorf("%w: position is required", ErrValidation)
	}
	if len(input.Position) > 120 {
		return core.Applicant{}, fmt.Errorf("%w: position must be at most 120 characters", ErrValidation)
	}
	created, err := s.repo.Create(ctx, core.Applicant{Name: input.Name, Email: input.Email, Position: input.Position, Stage: core.StageApplied})
	if err != nil {
		// SQLite and GORM report uniqueness violations differently across drivers.
		if strings.Contains(strings.ToLower(err.Error()), "unique constraint") || strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return core.Applicant{}, ErrDuplicateEmail
		}
		return core.Applicant{}, err
	}
	return created, nil
}
func (s *Service) List(ctx context.Context, filter core.ListFilter) (core.PageResult, error) {
	filter.Search = strings.TrimSpace(filter.Search)
	filter.Position = strings.TrimSpace(filter.Position)
	filter.Stage = strings.ToUpper(strings.TrimSpace(filter.Stage))
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}
	if filter.Limit > 100 {
		return core.PageResult{}, fmt.Errorf("%w: limit must not exceed 100", ErrValidation)
	}
	if filter.Stage != "" && !core.ValidStages[core.Stage(filter.Stage)] {
		return core.PageResult{}, fmt.Errorf("%w: invalid stage filter", ErrValidation)
	}
	return s.repo.List(ctx, filter)
}
func (s *Service) GetByID(ctx context.Context, id uint) (core.Applicant, error) {
	if id == 0 {
		return core.Applicant{}, fmt.Errorf("%w: id must be greater than zero", ErrValidation)
	}
	return s.repo.GetByID(ctx, id)
}
func (s *Service) Summary(ctx context.Context) (core.Summary, error) { return s.repo.Summary(ctx) }
func (s *Service) UpdateStage(ctx context.Context, id uint, stage core.Stage) (core.Applicant, error) {
	if id == 0 {
		return core.Applicant{}, fmt.Errorf("%w: id must be greater than zero", ErrValidation)
	}
	stage = core.Stage(strings.ToUpper(strings.TrimSpace(string(stage))))
	if !core.ValidStages[stage] {
		return core.Applicant{}, fmt.Errorf("%w: stage must be one of APPLIED, SCREENING, INTERVIEW, OFFER, HIRED, REJECTED", ErrValidation)
	}
	return s.repo.UpdateStage(ctx, id, stage)
}

var _ core.Service = (*Service)(nil)
