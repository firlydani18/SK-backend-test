package handler

import core "firlydev/backend-test/features/applicant"

type ApplicantResponse struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Position  string     `json:"position"`
	Stage     core.Stage `json:"stage"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
}
type PaginationMeta struct {
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}
type ListResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    []ApplicantResponse `json:"data"`
	Meta    PaginationMeta      `json:"meta"`
}

func toResponse(a core.Applicant) ApplicantResponse {
	return ApplicantResponse{ID: a.ID, Name: a.Name, Email: a.Email, Position: a.Position, Stage: a.Stage, CreatedAt: a.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"), UpdatedAt: a.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z07:00")}
}
