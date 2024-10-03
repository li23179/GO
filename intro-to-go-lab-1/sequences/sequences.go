package main

import "fmt"

func addOne(a int) int {
	return a + 1
}

func square(a int) int {
	return a * a
}

func double(slice []int) {
	slice = append(slice, slice...)
}

func mapSlice(f func(a int) int, slice []int) {
	for i, v := range slice{
		slice[i] = f(v)
	}
}

func mapArray(f func(a int) int, array *[3]int) {
	newArray := [3]int{}
	for i, v := range *array{
		newArray[i] = f(v)
	}
	*array = newArray
}

func main() {
	intsSlice := []int {1, 2, 3}
	mapSlice(addOne, intsSlice)
	fmt.Println(intsSlice)

	intsArrayPtr := new([3]int)
	*intsArrayPtr = [3]int {1, 2, 3}
	mapArray(addOne, intsArrayPtr)
	fmt.Println(*intsArrayPtr)
}
