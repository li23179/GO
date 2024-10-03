# Concurrency Lab 1

## Question 1 - Median Filter

### Question 1b

The below function is an example of a correct worker, but there are many other ways of doing it. For example instead of being given the start/end coords it could be given the size of the image and its 'position' (a number between 0 and 3) so that it can itself calculate where to apply the filter. What is key is that the closure is passed in and not a 2D slice.

```go
func worker(startY, endY, startX, endX int, data func(y, x int) uint8, out chan<- [][]uint8) {
    imagePart := medianFilter(startY, endY, startX, endX, data)
    out <- imagePart
}
```

### Question 1c

See `medianFilter.go`

### Question 1d

Your bar chart should look very similar to the one we presented in the lab-sheet. Please ask a TA for help if that's not the case. A slow solution could mean that you're applying the filter to the whole image in all workers - rather than just to the correct parts.

### Question 1e

In order to evaluate performance improvements due to parallelisation, we used Go benchmarks. We varied the number of worker threads and measured the total runtime of the Median Filter code. Each test used the `ship.png` input image and every measurement was repeated 5 times to reduce the impact of random noise. All measurements were collected on a university linux lab machine, which has a 6 core Intel i7-8700 CPU. The mean runtimes are presented above.

As the number of workers increased, the runtime decreased. With 16 worker threads, the filter was 2.8x faster than the serial implementation with 1 worker. The perfect improvement would be a runtime that halves as the number of threads is doubled. This is not the case here, because we measured the total time taken to apply the filter and output an image. The workers apply the filter in parallel, but the image input and output, which happens outside the workers, is not parallelisable.

The improvements also start to diminish, with the improvement from 8 to 16 workers being much smaller than from 1 to 2 workers. Code doesn't scale infinitely and since the machine's CPU has 6 cores, using any more than 6 worker threads offers virtually no improvement. The small improvement in runtime when using 16 rather than 8 workers is likely thanks to hyper-threading, which allows 6 physical cores to be treated as 12 logical cores by the operating system.

-------

The first paragraph states the setup of the experiment. It discusses variables, procedure and hardware used. It also shows consideration of noise in results - five runs of each benchmark ensures noise is reduced to a minimum. Note that the sentence "The mean runtimes are presented above" should read "The mean runtimes are presented in Figure 1" in an actual report. The figure should then have a short caption below it. Our bar chart would be even better if it included error bars to visualise variance in the data. With a larger number of configurations, a scatter plot with a trend line might be more suitable.

The second paragraph discusses the overall improvements and trends. It quotes the overall speedup of 2.8x (*not* 280% or 180%), and then explains the reasons behind it. It would be even better if the performance was evaluated on other images, which have different sizes. Further experiments could also be conducted to show that about 76% of CPU time (not to be confused with wall clock time) is spent in workers. The remaining 24% is mostly spent on input and output of images, which is not parallelisable. A CPU profile (Q1g) is one type of experiment which can provide such data.

The third paragraph discusses the result in context of the hardware of a lab machine. This is very important, because knowledge about the underlying hardware allows us to make reasonable assumptions about how our code *should* scale. This then allows us to compare these expectations with actual results. In more complicated problems, such as this unit's coursework, understanding of hardware will heavily influence our design decisions.

Finally, the answer features high quality of written English, short sentences, short paragraphs and appropriate use of terminology. Despite the possible improvements outlined above, this is an example of a first-class answer.

## Question 2 - Parallel Tree Reduction

### Question 2a

Note that solutions using waitgroups are also valid.

```go
func parallelMergeSort(slice []int, max int) {
	if len(slice) > 1 {
		if len(slice) <= max {
			mergeSort(slice)
		} else {
			middle := len(slice) / 2

			doneLeft := make(chan bool)
			doneRight := make(chan bool)

			go func() {
				parallelMergeSort(slice[:middle], max)
				doneLeft <- true
			}()

			go func() {
				parallelMergeSort(slice[middle:], max)
				doneRight <- true
			}()

			<-doneLeft
			<-doneRight

			merge(slice, middle)
		}
	}
}
```

### Question 2b

```go
func parallelMergeSort(slice []int, max int) {
	if len(slice) > 1 {
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
```

The performance should be better than the previous version, but still worse than sequential. The old version should be creating `O(n*log(n))` goroutines. The new version should be creating `n/2` goroutines.

### Question 2c

```go
func parallelMergeSort(slice []int, max int) {
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
```

Also see the `our analysis` section on the labsheet.
