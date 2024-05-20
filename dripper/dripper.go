package dripper

import (
	"runtime"
	"sync"
)

type Dripper struct {
	ch chan func() bool
	wg sync.WaitGroup
}

// NewDripper. start Dripper.
//
//	goroutines<=0: GOMAXPROCS(0)
func NewDripper(size int, goroutines int) *Dripper {
	ret := &Dripper{ch: make(chan func() bool, size)}
	if goroutines <= 0 {
		goroutines = runtime.GOMAXPROCS(0)
	}
	ret.wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			for f := range ret.ch {
				for !f() {
				}
			}
			ret.wg.Done()
		}()
	}

	return ret
}

// Wait. wait all func ok
func (r *Dripper) Wait() {
	close(r.ch)
	r.wg.Wait()
}

// Do. retry if f()!=true
func (r *Dripper) Do(f func() bool) { r.ch <- f }
