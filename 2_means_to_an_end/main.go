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

func tcpHandler(conn net.Conn, clients int) {
	defer conn.Close()
	db := map[uint32]uint32{}
	for {
		buff := make([]byte, 9)
		size, err := conn.Read(buff)
		if err != nil {
			fmt.Println("Can't read from connection: ", err)
			break
		}
		fmt.Printf("Client %d received -> %s (%b) of size: %d\n", clients, string(buff[0]), buff, size)
		resp, err := meansToAnEnd02(buff, db)
		if err != nil {
			fmt.Println(err)
			conn.Write([]byte("Invalid request\n"))
			continue
		}
		conn.Write(resp)
	}
}

func meansToAnEnd02(buff []byte, db map[uint32]uint32) ([]byte, error) {
	var resp = make([]byte, 4)
	t, x, y, err := deserializeMsg(buff)
	if err != nil {
		return nil, err
	}
	switch *t {
	case INSERT:
		fmt.Println("Inserting ", *y, " at time ", *x)
		timestamp := *x
		price := *y
		db[timestamp] = price
	case QUERY:
		fmt.Println("Querying MIX MAX : ", *x, *y)
		total, count := uint32(0), uint32(0)
		for i := *x; i <= *y; i++ {
			if v, ok := db[i]; ok {
				total += v
				count++
			}
		}
		var r uint32 = 0
		if count > 0 {
			r = total / count
		}
		binary.BigEndian.PutUint32(resp, r)
	}
	fmt.Printf("Ansering: %d - %b\n", binary.BigEndian.Uint16(resp), resp)
	return resp, nil
}

func deserializeMsg(msg []byte) (*string, *uint32, *uint32, error) {
	reqType := string(msg[0])
	if reqType != "I" && reqType != "Q" {
		return nil, nil, nil, errors.New(fmt.Sprintln("Invalid type: ", reqType, msg[0]))
	}
	i1, i2 := binary.BigEndian.Uint32(msg[1:5]), binary.BigEndian.Uint32(msg[5:])
	return &reqType, &i1, &i2, nil
}
