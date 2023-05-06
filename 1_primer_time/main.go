package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
)

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
	for {
		if conns <= uint64(maxPoolConnection-1) {
			// Holds inil a new connection is set
			conn, err := ln.Accept()
			if err != nil {
				log.Println("Error creating connection")
				panic(err)
			}
			go tcpHandler(conn) // One thread per client connection
		}
	}
}

func tcpHandler(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Reached EOF, closing connection")
				fmt.Println(err)
				break
			}
			break
		}
		// buf, err = smokeTest00(buf)
		buf, err = primeTime01(buf, n)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println("Responding:", string(append(buf, byte('\n'))))
		_, err = conn.Write(append(buf, byte('\n')))
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func smokeTest00(buf []byte, n int) ([]byte, error) {
	return buf[:n], nil // Echo
}

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

type PrimeReq struct {
	Method string `json:"method"`
	Number int    `json:"number"`
}
type PrimeResp struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func primeTime01(buf []byte, n int) ([]byte, error) {
	str := string(bytes.Trim(buf[:n], "\x00"))
	fmt.Println("Received string:", str)
	var req PrimeReq
	err := json.Unmarshal([]byte(str), &req)
	if err != nil {
		return nil, errors.New("Error parsing JSON: " + err.Error())
	}
	resp := PrimeResp{
		Method: "isPrime",
		Prime:  isPrime(req.Number),
	}
	return json.Marshal(resp)
}
