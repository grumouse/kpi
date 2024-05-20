package dripper_test

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	collection "github.com/grumouse/kpi/dripper"
	"github.com/stretchr/testify/assert"
)

func TestDripper(t *testing.T) {
	assert := assert.New(t)
	dr := collection.NewDripper(1024, 2)

	var n, rand, total int64
	var wg sync.WaitGroup
	const N = 1000
	for i := 0; i < N; i++ {
		i := i
		wg.Add(1)
		go func() {
			dr.Do(func() bool {
				cur := atomic.AddInt64(&n, 1)
				assert.LessOrEqual(cur, int64(2))
				atomic.AddInt64(&total, 1)
				cur = atomic.AddInt64(&n, -1)
				assert.GreaterOrEqual(cur, int64(0))
				ok := atomic.AddInt64(&rand, 1)&0b100 == 0
				fmt.Printf("%v: %v\n", i, ok)
				runtime.Gosched()

				return ok
			})
			wg.Done()
		}()
	}

	wg.Wait()
	dr.Wait()
	assert.Zero(n)
	assert.Equal(N*2, total)
	fmt.Println("END!")
}
