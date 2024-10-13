package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	sum := int32(0)

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			atomic.AddInt32(&sum, 1)
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println(sum)
}
