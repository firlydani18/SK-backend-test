package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	core "firlydev/backend-test/features/applicant"

	"github.com/labstack/echo/v4"
)

type testService struct {
	created core.Applicant
	items   []core.Applicant
	updated core.Applicant
}

func (s *testService) Create(_ interfaceContext, _ core.CreateInput) (core.Applicant, error) {
	return s.created, nil
}

// interfaceContext is an alias-compatible test helper for context.Context.
type interfaceContext = context.Context

func (s *testService) List(_ context.Context, f core.ListFilter) (core.PageResult, error) {
	return core.PageResult{Items: s.items, Page: int64(f.Page), Limit: int64(f.Limit), Total: int64(len(s.items))}, nil
}
func (s *testService) GetByID(_ context.Context, id uint) (core.Applicant, error) {
	return core.Applicant{ID: id}, nil
}
func (s *testService) Summary(context.Context) (core.Summary, error) {
	return core.Summary{Total: int64(len(s.items))}, nil
}
func (s *testService) UpdateStage(_ context.Context, id uint, stage core.Stage) (core.Applicant, error) {
	return core.Applicant{ID: id, Stage: stage}, nil
}

func TestCreateApplicantBadRequestForInvalidPayload(t *testing.T) {
	e := echo.New()
	h := New(&testService{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/applicants", strings.NewReader("{"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if err := h.Create(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestCreateApplicantReturnsCreated(t *testing.T) {
	e := echo.New()
	h := New(&testService{created: core.Applicant{ID: 7, Name: "Alya", Email: "alya@example.com", Position: "Backend", Stage: core.StageApplied}})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/applicants", strings.NewReader(`{"name":"Alya","email":"alya@example.com","position":"Backend"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if err := h.Create(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestUpdateStageReturnsUpdatedApplicant(t *testing.T) {
	e := echo.New()
	h := New(&testService{})
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/applicants/7/stage", strings.NewReader(`{"stage":"INTERVIEW"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("7")
	if err := h.UpdateStage(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}
