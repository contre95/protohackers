package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
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
	db := map[int32]int32{}
	for {
		buff := make([]byte, 9)
		size, err := io.ReadFull(conn, buff)
		// size, err := conn.Read(buff) // This doesn't work :(
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

func meansToAnEnd02(buff []byte, db map[int32]int32) ([]byte, error) {
	var resp = make([]byte, 4)
	t, x, y, err := deserializeMsg(buff)
	if err != nil {
		return nil, err
	}
	switch *t {
	case INSERT:
		// fmt.Println("Inserting ", *y, " at time ", *x)
		db[*x] = *y
	case QUERY:
		// fmt.Println("Querying MIX MAX : ", *x, *y)
		total, count := int32(0), int32(0)
		for i := *x; i <= *y; i++ {
			if v, ok := db[i]; ok {
				total += v
				count++
			}
		}
		var r int32 = 0
		if count > 0 {
			r = total / count
		}
		binary.BigEndian.PutUint32(resp, uint32(r))
		fmt.Printf("Ansering: %d - %b\n", int32(binary.BigEndian.Uint32(resp)), resp)
		return resp, nil
	}
	return nil, nil
}

func deserializeMsg(msg []byte) (*string, *int32, *int32, error) {
	reqType := string(msg[0])
	if reqType != "I" && reqType != "Q" {
		return nil, nil, nil, errors.New(fmt.Sprintf("Invalid type: %s - %b ", reqType, msg[0]))
	}
	i1, i2 := int32(binary.BigEndian.Uint32(msg[1:5])), int32(binary.BigEndian.Uint32(msg[5:]))
	return &reqType, &i1, &i2, nil
}
