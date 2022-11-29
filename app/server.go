package main

import (
	"fmt"
	"log"
	"net"
	"os"
	stringUtils "Strings"
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

	for {
		conn, err1 := listener.Accept()
		if err1 != nil {
			log.Fatal(err1)
		}

		go handleIncomingTCPRequest(&conn)
	}
}

func handleIncomingTCPRequest(connection *net.Conn) {
	buf := make([]byte, 1024)

	_, readErr := (*connection).Read(buf)
	if readErr != nil {
		fmt.Println("Error occur while reading from connection: ", readErr.Error())
	}

	message := string(buf[:])
	if isPingMessage(message) {
		handlePingMessage(connection, message)
	}

	closedErr := (*connection).Close()
	if closedErr != nil {
		fmt.Println("Error occur while closing the connection: ", closedErr.Error())
	}
}

func isPingMessage(message string) bool {
	return stringUtils.EqualFold("PING", message)
}

func handlePingMessage(connection *net.Conn, message string) {
	messageArray := stringUtils.Split(message, " ")
	if len(messageArray) > 2 {
		fmt.Println("Incorrect PING message")
		return
	}

	_ = messageArray[:1]
	clientMessage := stringUtils.Join(messageArray[1:], " ")

	if len(clientMessage) == 0 {
		(*connection).Write([]byte("PONG"))
		return
	}

	(*connection).Write([]byte(clientMessage))
	return
}
