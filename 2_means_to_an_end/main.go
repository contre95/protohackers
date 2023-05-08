package main

import (
	"encoding/binary"
	"errors"
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
	for {
		if conns <= uint64(maxPoolConnection-1) {
			// Holds inil a new connection is set
			conn, err := ln.Accept()
			if err != nil {
				log.Println("Error creating connection")
				panic(err)
			}
			fmt.Printf("Starting client %d\n", clients)
			go tcpHandler(conn, clients) // One thread per client connection
			clients++
		}
	}
}
func meansToAnEnd02(buff []byte, size int) ([]byte, error) {
	// fmt.Println("Received: ", buff)
	db := map[uint32]uint32{}
	var resp []byte
	t, x, y, err := deserializeMsg(buff)
	if err != nil {
		log.Println("Error deserializing message: ", err)
		return nil, err
	}
	switch t {
	case &INSERT:
		timestamp := *x
		price := *y
		db[timestamp] = price
	case &QUERY:
		total, count := uint32(0), uint32(0)
		for i := *x; i < *y; i++ {
			if v, ok := db[i]; ok {
				total += v
				count++
			}
		}
		binary.BigEndian.PutUint32(resp, total/count)
	}
	return resp, nil
}

func deserializeMsg(msg []byte) (*string, *uint32, *uint32, error) {
	reqType := string(msg[0])
	if reqType != "I" && reqType != "O" {
		return nil, nil, nil, errors.New("Wrong type")
	}
	i1, i2 := binary.BigEndian.Uint32(msg[1:5]), binary.BigEndian.Uint32(msg[5:])
	return &reqType, &i1, &i2, nil

}

func tcpHandler(conn net.Conn, clients int) {
	defer conn.Close()
	for {
		buff := make([]byte, 9)
		size, err := conn.Read(buff)
		if err != nil {
			log.Println("Error while reading from connection: ", err)
			break
		}
		resp, err := meansToAnEnd02(buff, size)
		if err != nil {
			fmt.Println("Malformed JSON: ", err)
			conn.Write([]byte("Malformed JSON\n"))
			break
		}
		conn.Write(append(resp, byte('\n')))
	}
}
