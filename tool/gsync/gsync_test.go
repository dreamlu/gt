package gsync

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

var (
	p = NewPool(100)
	g sync.WaitGroup
)

func TestPool(t *testing.T) {

	read := func() {
		defer g.Done()
		fmt.Printf("go func time: %d\n", time.Now().Unix())
		time.Sleep(time.Second)
	}

	for i := 0; i < 1000; i++ {
		g.Add(1)
		p.Submit(read)
	}

	g.Wait()
	//time.Sleep(5*time.Second)
}

func TestPool_SubmitWait(t *testing.T) {

	read := func() {
		fmt.Printf("go func time: %d\n", time.Now().Unix())
		time.Sleep(time.Second)
	}

	for i := 0; i < 10; i++ {
		p.Submit(read)
		readWait := func() interface{} {
			time.Sleep(time.Second)
			return i
		}

		t.Log(p.SubmitWait(readWait))
	}
}

func TestResultFunc(t *testing.T) {
	getGoroutineMemConsume()
}

func getGoroutineMemConsume() {
	var c chan struct{}
	var wg sync.WaitGroup
	const goroutineNum = 1e4 // 1 * 10^4

	memConsumed := func() uint64 {
		runtime.GC() //GC，排除对象影响
		var memStat runtime.MemStats
		runtime.ReadMemStats(&memStat)
		return memStat.Sys
	}

	noop := func() {
		wg.Done()
		<-c
	}

	wg.Add(goroutineNum)
	before := memConsumed()
	for i := 0; i < goroutineNum; i++ {
		go noop()
	}
	wg.Wait()
	after := memConsumed()

	fmt.Printf("%.3f KB\n", float64(after-before)/goroutineNum/1000)
}
