package main

import (
	"fmt"
	"io"
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
		go handleIncomingTCPRequest(conn)
	}
}

func handleIncomingTCPRequest(connection net.Conn) {
	defer connection.Close()

	for {
		buf := make([]byte, 1024)
		if _, readErr := connection.Read(buf); readErr != nil {
			if readErr == io.EOF {
				break
			}
		}

		message := ifPingThenReturnMessage(string(buf[:]))
		if message != "" {
			connection.Write([]byte(message))
		}
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
	messageType := MessageType(message[0])

	if MessageType(messageType) == arrays {
		contentArray := stringUtils.Split(message, "\r\n")
		if len(contentArray) >= 4 && contentArray[2] == "ping" {
			if len(contentArray) > 4 && contentArray[3] != "" {
				intTotalArrayElems, _ := strconv.Atoi(stringUtils.Split(contentArray[1], "$")[1])
				return "*" + strconv.Itoa(intTotalArrayElems-1) + stringUtils.Join(contentArray[3:], "")
			} else {
				return "+PONG\r\n"
			}
		}
	}
	return ""
}
