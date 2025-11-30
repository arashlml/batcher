package handler

import (
	"batcher/entity"
	"context"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

type Batcher[T any] interface {
	Add(item T) error
	OutChan() <-chan []T
}
type Repo interface {
	BatchInsert(ctx context.Context, items []entity.Request) error
}

type Handler struct {
	handlerCounter int64
	batcher        Batcher[entity.Request]
	repo           Repo
}

func NewHandler(batcher Batcher[entity.Request], repo Repo) *Handler {
	h := &Handler{
		batcher: batcher,
		repo:    repo,
	}
	go h.Consume()
	return h
}

func (h *Handler) AddToDatabase(c echo.Context) error {
	atomic.AddInt64(&h.handlerCounter, 1)
	log.Printf("HANDLER HIT = %d\n", atomic.LoadInt64(&h.handlerCounter))

	req := entity.Request{}
	if err := c.Bind(&req); err != nil {
		log.Printf("HANDELR: binding err -- > %v", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := h.batcher.Add(req); err != nil {
		log.Printf("HANDLER : error from batcher --> %v", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "ok!")
}
func (h *Handler) Consume() {
	channel := h.batcher.OutChan()

	for msg := range channel {
		log.Printf("HANDLER: received batch with %d items", len(msg))

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		err := h.repo.BatchInsert(ctx, msg)

		cancel()

		if err != nil {
			log.Printf("HANDLER : batch insert err -- > %v", err)
		}
	}

	log.Println("HANDLER: OutChan closed, consumer stopped")
}
