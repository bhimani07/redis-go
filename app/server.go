package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	stringUtils "strings"
)

const NETWORK_TYPE = "tcp"
const HOST = "localhost"
const PORT = "6379"

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	listener, err := net.Listen(NETWORK_TYPE, HOST+":"+PORT)
	defer listener.Close()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	for {
		conn, listenerError := listener.Accept()
		if listenerError != nil {
			log.Fatal("error while accepting connection: ", listenerError)
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

	message := ifPingThenReturnMessage(string(buf[:]))
	if message != "" {
		(*connection).Write([]byte(message))
	}

	closedErr := (*connection).Close()
	if closedErr != nil {
		fmt.Println("Error occur while closing the connection: ", closedErr.Error())
	}
}

type MessageType string

const (
	simpleStrings MessageType = "+"
	errors                    = "-"
	Integer                   = ":"
	bulkStrings               = "$"
	arrays                    = "*"
)

func ifPingThenReturnMessage(message string) string {
	var messageType = MessageType(message[0])

	if MessageType(messageType) == arrays {
		contentArray := stringUtils.Split(message, "\r\n")
		if len(contentArray) >= 4 && contentArray[2] == "ping" {
			if len(contentArray) > 4 && contentArray[3] != "" {
				intTotalArrayElem, _ := strconv.Atoi(stringUtils.Split(contentArray[1], "$")[1])
				return "*" + strconv.Itoa(intTotalArrayElem-1) + stringUtils.Join(contentArray[3:], "")
			} else {
				return "+PONG\r\n"
			}
		}
	}
	return ""
}
