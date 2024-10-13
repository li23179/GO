package main

import (
	"fmt"
)

func main() {
	sum := 0

	delta := make(chan int, 1000)

	for i := 0; i < 1000; i++ {
		go func() {
			delta <- 1
		}()
	}

	for i:=0; i < 1000; i++ {
		sum += <-delta
	}

	fmt.Println(sum)
}
