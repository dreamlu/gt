package gsync

import (
	"math"
	"sync"
	"time"
)

const (
	// DefaultPoolSize is the default capacity for a default goroutine pool.
	DefaultPoolSize = math.MaxInt32

	// DefaultRunTime is the interval time to clean up goroutines.
	DefaultRunTime = time.Second
)

type Pool struct {
	c  chan struct{}
	wg *sync.WaitGroup
	// work func run time
	//t time.Duration
}

func NewPool(size uint64) *Pool {
	if size == 0 {
		size = DefaultPoolSize
	}
	return &Pool{
		c:  make(chan struct{}, size),
		wg: new(sync.WaitGroup),
		//t:  DefaultRunTime,
	}
}

func (p *Pool) Submit(f func()) {

	go func() {
		p.add(1)
		f()
		defer p.done()
	}()
}

func (p *Pool) SubmitWait(f func() any) any {

	res := make(chan any)
	go func() {
		p.add(1)
		r := f()
		defer p.done()
		res <- r
	}()
	return <-res
}

func (p *Pool) add(delta int) {
	p.wg.Add(delta)
	for i := 0; i < delta; i++ {
		p.c <- struct{}{}
	}
}

func (p *Pool) done() {
	<-p.c
	p.wg.Done()
}

func (p *Pool) wait() {
	p.wg.Wait()
}
