package handler

import (
	"batcher/entity"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/labstack/echo/v4"
)

type Batcher[T any] interface {
	Add(item T) error
}
type Handler struct {
	Batcher Batcher[entity.Request]
}

var handlerCounter int64

func NewHandler(batcher Batcher[entity.Request]) *Handler {
	return &Handler{Batcher: batcher}
}
func (h *Handler) AddToDatabase(c echo.Context) error {
	atomic.AddInt64(&handlerCounter, 1)
	log.Printf("HANDLER HIT = %d\n", atomic.LoadInt64(&handlerCounter))

	req := entity.Request{}
	if err := c.Bind(&req); err != nil {
		log.Printf("HANDELR: binding err -- > %v", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := h.Batcher.Add(req); err != nil {
		log.Printf("HANDLER : error from batcher --> %v", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "ok!")
}
