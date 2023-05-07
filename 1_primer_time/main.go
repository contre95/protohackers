package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
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
	// buff := make([]byte, 1024)
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		req := scanner.Bytes()
		log.Printf("Client %d received -> %s", clients, req)
		resp, err := primeTime01(req)
		if err != nil {
			fmt.Println("Malformed JSON: ", err)
			conn.Write([]byte("Malformed JSON\n"))
			break
		}
		conn.Write(append(resp, byte('\n')))
	}
}

func primeTime01(buf []byte) ([]byte, error) {
	var req map[string]json.RawMessage
	err := json.Unmarshal(buf, &req)
	if err != nil {
		return nil, errors.New("error parsing JSON: " + err.Error())
	}
	methodBytes, ok := req["method"]
	if !ok {
		return nil, errors.New("method field not present")
	}
	numberBytes, ok := req["number"]
	if !ok {
		return nil, errors.New("number field not present")
	}
	var method string
	err = json.Unmarshal(methodBytes, &method)
	if err != nil {
		return nil, errors.New("error parsing method field: " + err.Error())
	}
	var number json.Number
	err = json.Unmarshal(numberBytes, &number)
	if err != nil {
		return nil, errors.New("error parsing number field: " + err.Error())
	}
	if method != "isPrime" {
		return nil, errors.New("unknown method: " + method)
	}
	var resp PrimeResp
	nu, err := number.Int64()
	if err != nil {
		resp = PrimeResp{
			Method: "isPrime",
			Prime:  false,
		}
		return json.Marshal(resp)
	}
	var req2 PrimeReq
	err = json.Unmarshal(buf, &req2)
	if err != nil {
		return nil, errors.New("wrong type" + method)
	}
	resp = PrimeResp{
		Method: "isPrime",
		Prime:  isPrime(nu),
	}

	if err != nil {
	}
	return json.Marshal(resp)
}

func isPrimeBig(num big.Int) bool {
	return num.ProbablyPrime(10)
}

func isPrime(num int64) bool {
	if num < 2 {
		fmt.Println("Number must be greater than 2.")
		return false
	}
	sq_root := int(math.Sqrt(float64(num)))
	for i := 2; i <= sq_root; i++ {
		if num%int64(i) == 0 {
			return false
		}
	}
	return true
}

type PrimeReq struct {
	Method *string `json:"method"`
	Number *int64  `json:"number"`
}

type PrimeResp struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}
