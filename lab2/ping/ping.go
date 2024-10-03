package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/trace"
	"time"
)

func foo(channel chan string) {
	for {
		fmt.Println("\nFoo is sending: ping")
		trace.Log(context.Background(), "Sent", "ping")
		channel <- "ping"

		message := <-channel
		fmt.Println("Foo has received:", message)
		trace.Log(context.Background(), "Received", message)
	}
}

func bar(channel chan string) {
	for {
		message := <-channel
		fmt.Println("Bar has received:", message)
		trace.Log(context.Background(), "Received", message)
		
		fmt.Println("Bar is sending: pong")
		trace.Log(context.Background(), "Sent", "pong")
		channel <- "pong"
	}
}

func pingPong() {
	pingPong := make(chan string)
	go foo(pingPong)
	go bar(pingPong)
	time.Sleep(500 * time.Millisecond)
}

func main() {
	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatalf("failed to create trace output file: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("failed to close trace file: %v", err)
		}
	}()

	if err := trace.Start(f); err != nil {
		log.Fatalf("failed to start trace: %v", err)
	}
	defer trace.Stop()

	pingPong()
}
