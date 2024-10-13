package main

import (
	"fmt"
)

func sum(slice []int) {
	if len(slice) > 1 {
		middle := len(slice) / 2

		done := make(chan bool)

		go func() {
			sum(slice[:middle])
			done <- true
		}()

		sum(slice[middle:])

		<-done

		slice[0] = slice[0] + slice[middle]
	}
}

func main() {
	slice := make([]int, 1000)

	for i := 0; i < 1000; i++ {
		slice[i] = 1
	}

	sum(slice)

	fmt.Println(slice[0])
}
