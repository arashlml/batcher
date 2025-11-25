package batcher

import (
	"batcher/entity"
	"context"
	"errors"
	"log"
	"time"
)

type Repository interface {
	BatchInsert(ctx context.Context, items []entity.Request) error
}

type Batcher struct {
	itemsCh  chan entity.Request
	batch    []entity.Request
	maxSize  int
	repo     Repository
	interval time.Duration
	quit     chan struct{}
}

func NewBatcher(maxSize int, interval time.Duration, repo Repository) *Batcher {
	log.Println("Making the batcher")
	b := &Batcher{
		itemsCh:  make(chan entity.Request, maxSize*4),
		batch:    make([]entity.Request, 0, maxSize),
		maxSize:  maxSize,
		interval: interval,
		quit:     make(chan struct{}),
		repo:     repo,
	}

	go b.run()

	return b
}

func (b *Batcher) Add(item entity.Request) error {
	select {
	case b.itemsCh <- item:
		return nil
	default:
		return errors.New("buffer is full")
	}
}

func (b *Batcher) run() {
	ticker := time.NewTicker(b.interval)
	defer ticker.Stop()

	for {
		select {
		case item := <-b.itemsCh:
			b.batch = append(b.batch, item)
			if len(b.batch) >= b.maxSize {
				log.Println("batch is full. let's flush")
				b.flush()
			}
		case <-ticker.C:
			if len(b.batch) == 0 {
				continue
			}
			log.Println("Flushing time has come")
			b.flush()
		case <-b.quit:
			log.Println("quitting")
			b.flush()
			return
		}
	}
}
func (b *Batcher) flush() {
	if len(b.batch) == 0 {
		log.Println("Nothing to flush")
		return
	}

	batchToInsert := b.batch
	b.batch = make([]entity.Request, 0, b.maxSize)
	log.Println("i've flushed the batch")
	if err := b.repo.BatchInsert(context.Background(), batchToInsert); err != nil {
		log.Printf("error inserting the batch --> %v ", err)
	}
}

func (b *Batcher) Close() {
	close(b.quit)
}
