package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
)

var INSERT string = "I"
var QUERY string = "Q"

func main() {
	err := StartTCPServer(8584, "0.0.0.0", 1000)
	if err != nil {
		fmt.Printf("Could not start server: %v", err)
	}
}

func StartTCPServer(port int, host string, maxPoolConnection int) error {
	portStr := strconv.Itoa(port)
	fmt.Println("Listening TCP on PORT " + portStr)
	ln, err := net.Listen("tcp", host+":"+portStr)
	if err != nil {
		return err
	}
	var conns uint64
	clients := 0
	// var room chan string
	room := make(chan string, 100)
	for {
		if conns <= uint64(maxPoolConnection-1) {
			// Holds inil a new connection is set
			conn, err := ln.Accept()
			if err != nil {
				log.Println("Error creating connection")
				panic(err)
			}
			fmt.Printf("Starting client %d\n", clients)
			go tcpHandler(conn, room, clients) // One thread per client connection
			clients++
		}
	}
}

// isValidUsername The first message from a client sets the user's name, which must contain at least 1 character, and must consist entirely of alphanumeric characters (uppercase, lowercase, and digits). Implementations may limit the maximum length of a name, but must allow at least 16 characters. Implementations may choose to either allow or reject duplicate names.
func isValidUsername(name string) bool {
	if len(name) > 0 {
		return true
	}
	return false
}

func tcpHandler(conn net.Conn, room chan string, clients int) {
	defer conn.Close()
	// Get the name
	conn.Write([]byte("Welcome to budgetchat! What shall I call you?\n"))
	var name string
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		name = string(scanner.Bytes())
		fmt.Printf("%s asdf", string(name))
		if isValidUsername(string(name)) {
			log.Printf("%s is a valid username ", string(name))
			room <- fmt.Sprintf("* %s has entered the room\n", string(name))
			go func() {
				for {
					conn.Write([]byte(<-room))
				}
			}()
			break
		}
	}

	for scanner.Scan() {
		room <- "[" + string(name) + "] " + string(scanner.Bytes()) + "\n"
	}
}
