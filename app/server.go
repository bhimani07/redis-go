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

type CommandType string

const (
	ping    CommandType = "ping"
	echo                = "echo"
	unknown             = "unknown"
)

type MessageDataType string

const (
	simpleStrings MessageDataType = "+"
	errors                        = "-"
	Integer                       = ":"
	bulkStrings                   = "$"
	arrays                        = "*"
)

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
		buf := make([]byte, 2048)
		if _, readErr := connection.Read(buf); readErr != nil {
			if readErr == io.EOF {
				break
			}
		}

		fmt.Println("Command Type: ", determineCommandType(string(buf[:])))
		switch determineCommandType(string(buf[:])) {
		case ping:
			connection.Write([]byte(buildPingResponse(string(buf[:]))))
		case echo:
			connection.Write([]byte(buildEchoResponse(string(buf[:]))))
		case unknown:
			fmt.Println("Unknown cmd received, exiting...")
			os.Exit(1)
		}
	}
}

func determineCommandType(message string) CommandType {
	messageType := MessageDataType(message[0])

	if MessageDataType(messageType) == arrays {
		contentArray := stringUtils.Split(message, "\r\n")
		if CommandType(contentArray[2]) == ping {
			return ping
		} else if CommandType(contentArray[2]) == echo {
			return echo
		}
	}
	return unknown
}

func buildPingResponse(message string) string {
	contentArray := stringUtils.Split(message, "\r\n")

	if len(contentArray) > 4 && contentArray[3] != "" {
		intTotalArrayElems, _ := strconv.Atoi(stringUtils.Split(contentArray[1], "$")[1])
		return "*" + strconv.Itoa(intTotalArrayElems-1) + stringUtils.Join(contentArray[3:], "")
	} else {
		return "+PONG\r\n"
	}
}

func buildEchoResponse(message string) string {
	contentArray := stringUtils.Split(message, "\r\n")
	response := ""
	for _, elem := range contentArray[3:] {
		response += elem
		response += "\r\n"
	}
	return response
}
