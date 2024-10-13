package main

import (
	"flag"
	// "net"
	"net/rpc"

	"bufio"
	"fmt"
	"os"
	"uk.ac.bris.cs/distributed2/secretstrings/stubs"
)

// Stage 3: Update your client to read words from the wordlist file and reverse them all

func main() {
	server := flag.String("server", "127.0.0.1:8030", "IP:port string to connect to as server")
	flag.Parse()
	fmt.Println("Server: ", *server)
	//TODO: connect to the RPC server and send the request(s)
	client, _ := rpc.Dial("tcp", *server)
	defer client.Close()

	wordList, _ := os.Open("../wordlist")
	scanner := bufio.NewScanner(wordList)
	
	request := stubs.Request{Message: ""}
	response := new(stubs.Response)

	for scanner.Scan(){
		text := scanner.Text()
		request.Message = text
		client.Call(stubs.PremiumReverseHandler, request, response)
		fmt.Println("Response : " + response.Message)
	}

}
