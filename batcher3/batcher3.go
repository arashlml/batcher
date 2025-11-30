package batcherV3

import (
	"context"
	"errors"
	"log"
	"sync/atomic"
	"time"
)

type Batcher[T any] struct {
	itemsCh      chan T
	batch        []T
	maxSize      int
	function     func(context.Context, []T) error
	interval     time.Duration
	quit         chan struct{}
	cancelTime   time.Duration
	outChan      chan []T
	addCounter   int64
	flushCounter int64
}

func NewBatcher[T any](maxSize int, interval time.Duration) *Batcher[T] {
	log.Println("BATCHER: MAKING THE BATCHER")
	b := &Batcher[T]{
		itemsCh:  make(chan T, 1000),
		batch:    make([]T, 0, maxSize),
		maxSize:  maxSize,
		interval: interval,
		quit:     make(chan struct{}),
		outChan:  make(chan []T, 100),
	}
	log.Printf("🚀 Batcher Initialized Successfully\n\nReady to collect, batch, and flush your data with:\n- ⏱ Time-based batching\n- 📦 Size-based flushing\n- ⚡ High-throughput processing\n\nLet the batching begin.\n")
	go b.run()
	return b
}
func (b *Batcher[T]) OutChan() <-chan []T {
	return b.outChan
}

func (b *Batcher[T]) Add(item T) error {
	select {
	case b.itemsCh <- item:
		atomic.AddInt64(&b.addCounter, 1)
		log.Printf("ADD COUNTER : %d", atomic.LoadInt64(&b.addCounter))
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
				log.Println("BATCHER: STOPING THE TICKER")
				ticker.Stop()
				ticker = nil
				tickerCh = nil
			}
		case <-b.quit:
			log.Println("BATCHER: QUITING THE BATCHER")
			b.flush()
			close(b.itemsCh)
			close(b.quit)
			close(b.outChan)
			return
		}
	}
}

func (b *Batcher[T]) flush() {
	batchToInsert := make([]T, len(b.batch))
	log.Printf("BATCHER: LENGH OF BATCH IS: %d \n", len(b.batch))
	copy(batchToInsert, b.batch)
	b.batch = make([]T, 0, b.maxSize)

	if len(batchToInsert) != 0 {
		atomic.AddInt64(&b.flushCounter, 1)
		log.Printf("BATCHER: FLUSH CALLED : %d \n", atomic.LoadInt64(&b.flushCounter))
	}
	b.outChan <- batchToInsert
}

func (b *Batcher[T]) Close() {
	close(b.quit)
}
