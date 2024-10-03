package main

import (
	"math/rand"
)

func merge(slice []int32, middle int) {
	sliceClone := make([]int32, len(slice))
	copy(sliceClone, slice)
	a := sliceClone[middle:]
	b := sliceClone[:middle]
	i := 0
	j := 0
	for k := 0; k < len(slice); k++ {
		if i >= len(a) {
			slice[k] = b[j]
			j++
		} else if j >= len(b) {
			slice[k] = a[i]
			i++
		} else if a[i] > b[j] {
			slice[k] = b[j]
			j++
		} else {
			slice[k] = a[i]
			i++
		}
	}
}

func mergeSort(slice []int32) {
	if len(slice) > 1 {
		middle := len(slice) / 2
		mergeSort(slice[:middle])
		mergeSort(slice[middle:])
		merge(slice, middle)
	}
}

func parallelMergeSort(slice []int32, max int) {
	if len(slice) > 1 {
		if len(slice) <= max {
			mergeSort(slice)
		} else {
			middle := len(slice) / 2
			done := make(chan bool)

			go func() {
				parallelMergeSort(slice[:middle], max)
				done <- true
			}()

			parallelMergeSort(slice[middle:], max)
			<-done
			merge(slice, middle)
		}
	}
}

func main() {
	//f, err := os.Create("trace.out")
	//if err != nil {
	//	log.Fatalf("failed to create trace output file: %v", err)
	//}
	//defer func() {
	//	if err := f.Close(); err != nil {
	//		log.Fatalf("failed to close trace file: %v", err)
	//	}
	//}()
	//
	//if err := trace.Start(f); err != nil {
	//	log.Fatalf("failed to start trace: %v", err)
	//}
	//defer trace.Stop()

	s := make([]int32, 100)
	for i := 0; i < len(s); i++ {
		s[i] = rand.Int31()
	}
	parallelMergeSort(s, 0)
}
