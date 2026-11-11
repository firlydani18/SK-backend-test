package responses

import "github.com/labstack/echo/v4"

type Envelope struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func Success(c echo.Context, status int, message string, data interface{}) error {
	return c.JSON(status, Envelope{Success: true, Message: message, Data: data})
}
func Failure(c echo.Context, status int, message string, detail interface{}) error {
	return c.JSON(status, Envelope{Success: false, Message: message, Error: detail})
}
