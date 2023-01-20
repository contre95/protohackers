package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

func main() {
	err := StartTCPServer(8584, "0.0.0.0", 150)
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
	for {
		if conns <= uint64(maxPoolConnection-1) {
			// Holds inil a new connection is set
			conn, err := ln.Accept()
			atomic.AddUint64(&conns, 1)
			if err != nil {
				log.Println("Error creating connection")
				panic(err)
			}
			go tcpHandler(&conn)
		}
	}
}

func tcpHandler(conn *net.Conn) {}
