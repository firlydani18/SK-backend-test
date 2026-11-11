package routers

import (
	core "firlydev/backend-test/features/applicant"
	"firlydev/backend-test/features/applicant/handler"

	"github.com/labstack/echo/v4"
)

func Register(e *echo.Echo, s core.Service) {
	h := handler.New(s)
	e.GET("/health", func(c echo.Context) error { return c.JSON(200, map[string]string{"status": "ok"}) })
	v1 := e.Group("/api/v1")
	v1.POST("/applicants", h.Create)
	v1.GET("/applicants", h.List)
	v1.GET("/applicants/summary", h.Summary)
	v1.GET("/applicants/:id", h.GetByID)
	v1.PATCH("/applicants/:id/stage", h.UpdateStage)
}
