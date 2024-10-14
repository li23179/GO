package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"time"
)

// Each process needs to accept an ip:port string for
// the 'next' buddy who will follow on from them in the song.
// You'll have to configure them in a loop.

// You don't want clients to try connect to each other straight away,
// or you won't have time to set the final process running so that the first can connect.

// When you set up the processes, you'll also need some way to indicate 
// which of them should start the song (I suggest allowing any n bottles of beer, for testing purposes).
// Only the last process you set up should need to be told the n to count down from.

var initialised = false
var nextAddr string
var nextPerson *rpc.Client

type Token struct{
	Bottles int
}

type Round struct{}

var turnOperation = "Round.NextRound"

func (r *Round) NextRound(req Token, res *Token){
	bottles := req.Bottles
	Sing(bottles)
	if bottles > 0{
		Config(bottles - 1)
	}
}

func Config(bottles int){
	if !initialised{
		nextPerson, _ = rpc.Dial("tcp", nextAddr)
		initialised = true
	}

	req := Token{Bottles: bottles}
	res := new(Token)

	nextPerson.Go(turnOperation, req, res, nil)
}

func Sing(bottles int){
	time.Sleep(1 * time.Second)
	if bottles > 0 {
		fmt.Printf("%v bottles of beer on the wall, %v bottles of beer. Take one down, pass it around...\n", 
			bottles, bottles)
	} else {
		fmt.Println("NO MORE BEERS!!!")
	}
}

func main() {
	thisPort := flag.String("this", "8030", "Port for this process to listen on")
	flag.StringVar(&nextAddr, "next", "localhost:8040", "IP:Port string for next member of the round.")
	bottles := flag.Int("n", 0, "Bottles of Beer (launches song if not 0)")
	flag.Parse()

	//TODO: Up to you from here! Remember, you'll need to both listen for
	//RPC calls and make your own.

	listener, _ := net.Listen("tcp", ":" + *thisPort)
	defer listener.Close()

	if *bottles > 0{
		Sing(*bottles)
		go Config(*bottles - 1)
	}
	rpc.Accept(listener)

}
