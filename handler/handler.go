package handler

import (
	"batcher/entity"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Batcher interface {
	Add(item entity.Request) error
}

type Handler struct {
	batcher Batcher
}

func NewHandler(batcher Batcher) *Handler {
	return &Handler{batcher: batcher}
}
func (h *Handler) AddToDatabase(c echo.Context) error {
	req := entity.Request{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := h.batcher.Add(req); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "ok!")
}
