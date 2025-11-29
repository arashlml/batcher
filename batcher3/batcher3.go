package batcherV3

import (
	"context"
	"errors"
	"log"
	"time"
)

type Batcher[T any] struct {
	itemsCh    chan T
	batch      []T
	maxSize    int
	function   func(context.Context, []T) error
	interval   time.Duration
	quit       chan struct{}
	cancelTime time.Duration
	outChan    chan []T
}

func NewBatcher[T any](maxSize int, interval time.Duration, function func(context.Context, []T) error, cancelTime time.Duration) *Batcher[T] {
	log.Println("BATCHER: MAKING THE BATCHER")
	b := &Batcher[T]{
		itemsCh:    make(chan T, 0, 2*maxSize),
		batch:      make([]T, 0, maxSize),
		maxSize:    maxSize,
		interval:   interval,
		quit:       make(chan struct{}),
		cancelTime: cancelTime,
	}
	go b.run()
	return b
}
func (b *Batcher[T]) OutChan() <-chan []T {
	return b.outChan
}

func (b *Batcher[T]) Add(item T) error {
	select {
	case b.itemsCh <- item:
		return nil
	default:
		return errors.New("BATCHER: BUFFER IS FULL")
	}
}

func (b *Batcher[T]) run() {
	var ticker *time.Ticker
	var tickerCh <-chan time.Time
	for {
		select {
		case item := <-b.itemsCh:
			b.batch = append(b.batch, item)
			if len(b.batch) == 1 {
				log.Println("BATCHER: TICK TOCK  ( ﾉ ﾟｰﾟ)ﾉ")
				ticker = time.NewTicker(b.interval)
				tickerCh = ticker.C
			}
			if len(b.batch) >= b.maxSize {
				log.Println("BATCHER: BATCH IS FULL, I HAVE TO FLUSH! 🚽")
				b.flush()
				ticker.Stop()
			}
		case <-tickerCh:
			log.Println("BATCHER: FLUSHING TIME HAS COME 🚽")
			b.flush()
			if ticker != nil {
				log.Println("BATCHER: TICKER IS NI")
				ticker.Stop()
				ticker = nil
				tickerCh = nil
			}
		case <-b.quit:
			log.Println("BATCHER: QUITING THE BATCHER")
			b.flush()
			close(b.itemsCh)
			close(b.quit)
			return
		}
	}
}

func (b *Batcher[T]) flush() {
	batchToInsert := b.batch
	b.outChan <- batchToInsert
}

func (b *Batcher[T]) Close() {
	close(b.quit)
}
