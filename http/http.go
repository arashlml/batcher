package http

import (
	"batcher/handler"

	"github.com/labstack/echo/v4"
)

func RegisterAPI(e *echo.Echo, h *handler.Handler) {
	e.POST("/batcher", h.AddToDatabase)
}
