# Intro to Go Lab 2

## Question 1 - Messages

### Question 1a

3 messages are sent and received. The two goroutines have to sync on send/receive.

### Question 1b

The first 2 messages have been received. Then `main()` exits and therefore so does the whole program. The third message is therefore lost. You may see the the line `sendMessages is sending: pinggg` but that 3rd message is not actually sent as there is no corresponding receive operation.

### Question 1c

`main()` is stuck trying to receive a 4th message that does not exist. The `sendMessages()` goroutine has already terminated and no other goroutine exists to be able to send the 4th message. Therefore the Go language runtime reports a deadlock.

### Question 1d

All 3 messages are quickly sent and queued up on the buffer, and then `main()` reads them off the buffer one by one, 1s apart. The important distinction to the original solution is that the two goroutines no longer need to sync on send/receive and therefore `sendMessages` goroutine terminates as soon as it's queued the messages up -before `main()` had a chance to actually receive anything.

### Question 2a

```go
func foo(channel chan string) {
    fmt.Println("\nFoo is sending: ping")
    channel <- "ping"

    message := <-channel
    fmt.Println("Foo has received:", message)
}

func bar(channel chan string) {
    message := <-channel
    fmt.Println("Bar has received:", message)

    fmt.Println("Bar is sending: pong")
    channel <- "pong"
}

func pingPong() {
    pingPong := make(chan string)
    go foo(pingPong)
    go bar(pingPong)
    time.Sleep(500 * time.Millisecond)
}
```

### Question 2b

See `ping.go`

### Question 2c

In summary, each block is a running goroutine and each arrow is a ping or a pong. The shape of the trace directly shows how the messages are "ping-ponging" and it also illustrates that exactly one of `foo` and `bar` are blocked at any one time. Read the hint to the original question and watch the video on tracing if anything isn't clear.

![Flow enable](content/scheduler.png)

In some cases you may have observed that for a few ms the ping ponging pauses and either `bar()` or `foo()` seems to be taking a very long time. This is not a bug. It simply shows that another process was scheduled to run by the OS for those few ms and our ping pong program was paused. The lesson here is that in order to get acurate traces and benchmarks we have to ensure that as few other processes are running at the same time as possible.



## Question 3 - for-select

### Question 3b

```go
func fasterSender(c chan<- []int) {
	for {
		time.Sleep(200 * time.Millisecond)
		c <- []int{1, 2, 3}
	}
}


...


for {
		select {
		case s := <-strings:
			fmt.Println("Received a string", s)
		case i := <-ints:
			fmt.Println("Received an int", i)
		case s := <-slices:
			fmt.Println("Received a slice", s)
		}
	}

```

### Question 3c / 3d



```go
	for {
		select {
		case s := <-strings:
			fmt.Println("Received a string", s)
		case i := <-ints:
			fmt.Println("Received an int", i)
		case s := <-slices:
			fmt.Println("Received a slice", s)
		default:
			fmt.Println("--- Nothing to receive, sleeping for 3s...")
			time.Sleep(3 * time.Second)
		}
	}
```

With no buffer, once every 3 seconds, 3 messages are received. 1 each from strings, ints and slices. This is because main and the sender have to synchronise.

This is not the case when using buffered channels. When a buffered channel is in use the senders will keep sending messages until the buffer is full. In other words, they can send messages even if the main thread is blocked/sleeping/deadlocked.

## Question 4 - Quiz

### Question 4a

See `quizA`.

### Question 4b

See `quizB`.


