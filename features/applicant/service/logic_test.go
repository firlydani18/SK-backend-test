package service

import (
	"context"
	"errors"
	"testing"

	core "firlydev/backend-test/features/applicant"
)

type fakeRepo struct {
	created core.Applicant
	list    []core.Applicant
	update  core.Applicant
	err     error
}

func (f *fakeRepo) Create(_ context.Context, a core.Applicant) (core.Applicant, error) {
	if f.err != nil {
		return core.Applicant{}, f.err
	}
	a.ID = 1
	f.created = a
	return a, nil
}
func (f *fakeRepo) List(_ context.Context, filter core.ListFilter) (core.PageResult, error) {
	return core.PageResult{Items: f.list, Page: int64(filter.Page), Limit: int64(filter.Limit), Total: int64(len(f.list))}, f.err
}
func (f *fakeRepo) GetByID(_ context.Context, id uint) (core.Applicant, error) {
	return core.Applicant{ID: id}, f.err
}
func (f *fakeRepo) Summary(context.Context) (core.Summary, error) {
	return core.Summary{Total: int64(len(f.list))}, f.err
}
func (f *fakeRepo) UpdateStage(_ context.Context, id uint, stage core.Stage) (core.Applicant, error) {
	if f.err != nil {
		return core.Applicant{}, f.err
	}
	f.update = core.Applicant{ID: id, Stage: stage}
	return f.update, nil
}
func TestCreateApplicantNormalizesAndDefaultsStage(t *testing.T) {
	repo := &fakeRepo{}
	svc := New(repo)
	got, err := svc.Create(context.Background(), core.CreateInput{Name: "  Firly  ", Email: " FIRLY@EXAMPLE.COM ", Position: " Backend Engineer "})
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "Firly" || got.Email != "firly@example.com" || got.Position != "Backend Engineer" || got.Stage != core.StageApplied {
		t.Fatalf("unexpected applicant: %+v", got)
	}
}
func TestCreateRejectsInvalidEmail(t *testing.T) {
	svc := New(&fakeRepo{})
	_, err := svc.Create(context.Background(), core.CreateInput{Name: "Firly", Email: "not-an-email", Position: "Backend"})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
}
func TestListAppliesPaginationDefaults(t *testing.T) {
	svc := New(&fakeRepo{})
	got, err := svc.List(context.Background(), core.ListFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if got.Page != 1 || got.Limit != 10 {
		t.Fatalf("unexpected pagination defaults: %+v", got)
	}
}
func TestListRejectsLimitAboveMaximum(t *testing.T) {
	svc := New(&fakeRepo{})
	_, err := svc.List(context.Background(), core.ListFilter{Page: 1, Limit: 101})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
}
func TestUpdateStageRejectsUnknownStage(t *testing.T) {
	svc := New(&fakeRepo{})
	_, err := svc.UpdateStage(context.Background(), 1, "UNKNOWN")
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
}
