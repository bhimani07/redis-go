package main

import (
	"fmt"
	"net"
	"os"
	"log"
)

const NETWORK_TYPE = "tcp"
const HOST = "localhost"
const PORT = "6379"

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	listener, err := net.Listen(NETWORK_TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	_, err1 := listener.Accept()
	if err1 != nil {
		log.Fatal(err1)
		os.Exit(1)
	}

}
