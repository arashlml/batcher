package batcher

import (
	"context"
	"errors"
	"log"
	"sync/atomic"
	"time"
)

type Batcher[T any] struct {
	itemsCh  chan T
	batch    []T
	maxSize  int
	Function func(context.Context, []T) error
	interval time.Duration
	quit     chan struct{}
}

func NewBatcher[T any](maxSize int, interval time.Duration, function func(context.Context, []T) error) *Batcher[T] {
	log.Println("BATCHER: MAKING THE BATCHER")
	b := &Batcher[T]{
		itemsCh:  make(chan T, 1000),
		batch:    make([]T, 0, maxSize),
		maxSize:  maxSize,
		interval: interval,
		quit:     make(chan struct{}),
		Function: function,
	}
	go b.run()
	return b
}

var addCounter int64

func (b *Batcher[T]) Add(item T) error {
	select {
	case b.itemsCh <- item:
		atomic.AddInt64(&addCounter, 1)
		log.Printf("BATCHER: BATCHER ADD = %d\n", atomic.LoadInt64(&addCounter))
		return nil
	default:
		return errors.New("BATCHER: BUFFER IS FULL")
	}
}

func (b *Batcher[T]) run() {
	ticker := time.NewTicker(b.interval)
	defer ticker.Stop()
	for {
		select {
		case item := <-b.itemsCh:
			b.batch = append(b.batch, item)
			if len(b.batch) >= b.maxSize {
				log.Println("BATCHER: BATCH IS FULL, I HAVE TO FLUSH! 🚽")
				b.flush()
			}
		case <-ticker.C:
			log.Println("BATCHER: FLUSHING TIME HAS COME 🚽")
			b.flush()
		case <-b.quit:
			log.Println("BATCHER: QUITING THE BATCHER")
			b.flush()
			return
		}
	}
}

var flushCounter int64

func (b *Batcher[T]) flush() {
	batchToInsert := b.batch
	b.batch = make([]T, 0, b.maxSize)
	if err := b.Function(context.Background(), batchToInsert); err != nil {
		log.Printf("BATCHER: ERROR SENDING THE BATCH --> %v ", err)
	} else {
		atomic.AddInt64(&flushCounter, 1)
		log.Printf("BATCHER: FLUSH CALLED : %d \n", atomic.LoadInt64(&flushCounter))
	}
}

func (b *Batcher[T]) Close() {
	close(b.quit)
}
