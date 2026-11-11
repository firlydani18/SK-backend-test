package handler

import (
	"errors"
	"strconv"

	core "firlydev/backend-test/features/applicant"
	"firlydev/backend-test/features/applicant/data"
	"firlydev/backend-test/features/applicant/service"
	"firlydev/backend-test/utils/responses"

	"github.com/labstack/echo/v4"
)

type Handler struct{ service core.Service }

func New(s core.Service) *Handler { return &Handler{service: s} }
func (h *Handler) Create(c echo.Context) error {
	var req CreateRequest
	if err := c.Bind(&req); err != nil {
		return responses.Failure(c, 400, "invalid request body", err.Error())
	}
	item, err := h.service.Create(c.Request().Context(), core.CreateInput{Name: req.Name, Email: req.Email, Position: req.Position})
	if err != nil {
		return h.handleError(c, err)
	}
	return responses.Success(c, 201, "applicant created", toResponse(item))
}
func (h *Handler) List(c echo.Context) error {
	page, err := queryInt(c, "page", 1)
	if err != nil {
		return responses.Failure(c, 400, "invalid query parameter", "page must be a positive integer")
	}
	limit, err := queryInt(c, "limit", 10)
	if err != nil {
		return responses.Failure(c, 400, "invalid query parameter", "limit must be a positive integer")
	}
	result, err := h.service.List(c.Request().Context(), core.ListFilter{Search: c.QueryParam("search"), Stage: c.QueryParam("stage"), Position: c.QueryParam("position"), Page: page, Limit: limit})
	if err != nil {
		return h.handleError(c, err)
	}
	items := make([]ApplicantResponse, 0, len(result.Items))
	for _, a := range result.Items {
		items = append(items, toResponse(a))
	}
	meta := PaginationMeta{Page: result.Page, Limit: result.Limit, Total: result.Total, TotalPages: (result.Total + result.Limit - 1) / result.Limit}
	return c.JSON(200, ListResponse{Success: true, Message: "applicants retrieved", Data: items, Meta: meta})
}
func (h *Handler) GetByID(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return responses.Failure(c, 400, "invalid applicant id", "id must be a positive integer")
	}
	item, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return h.handleError(c, err)
	}
	return responses.Success(c, 200, "applicant retrieved", toResponse(item))
}
func (h *Handler) UpdateStage(c echo.Context) error {
	id, err := parseID(c)
	if err != nil {
		return responses.Failure(c, 400, "invalid applicant id", "id must be a positive integer")
	}
	var req UpdateStageRequest
	if err := c.Bind(&req); err != nil {
		return responses.Failure(c, 400, "invalid request body", err.Error())
	}
	item, err := h.service.UpdateStage(c.Request().Context(), id, core.Stage(req.Stage))
	if err != nil {
		return h.handleError(c, err)
	}
	return responses.Success(c, 200, "applicant stage updated", toResponse(item))
}
func (h *Handler) Summary(c echo.Context) error {
	result, err := h.service.Summary(c.Request().Context())
	if err != nil {
		return h.handleError(c, err)
	}
	return responses.Success(c, 200, "applicant summary retrieved", result)
}
func (h *Handler) handleError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, service.ErrValidation):
		return responses.Failure(c, 400, "validation failed", err.Error())
	case errors.Is(err, service.ErrDuplicateEmail):
		return responses.Failure(c, 409, "email already exists", err.Error())
	case errors.Is(err, data.ErrNotFound):
		return responses.Failure(c, 404, "applicant not found", err.Error())
	default:
		return responses.Failure(c, 500, "internal server error", "please try again later")
	}
}
func parseID(c echo.Context) (uint, error) {
	v, e := strconv.ParseUint(c.Param("id"), 10, 64)
	if e != nil || v == 0 {
		return 0, errors.New("invalid id")
	}
	return uint(v), nil
}
func queryInt(c echo.Context, key string, def int) (int, error) {
	raw := c.QueryParam(key)
	if raw == "" {
		return def, nil
	}
	v, e := strconv.Atoi(raw)
	if e != nil || v < 1 {
		return 0, errors.New("invalid integer")
	}
	return v, nil
}
