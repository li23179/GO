package main

import (
	"fmt"
	"sync"
)

func main() {
	sum := int32(0)

	var mutex = &sync.Mutex{}

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			mutex.Lock()
			sum= sum + 1
			mutex.Unlock()
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println(sum)
}
