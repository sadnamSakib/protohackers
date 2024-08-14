package primetime

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
)

type Request struct {
	Method *string  `json:"method"`
	Number *float64 `json:"number"`
}
type Response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}
type BadResponse struct {
	Bad string `json:"bad"`
}

const BadResponseJSON = `{"bad":"bad request"}`

func checkPrime(req Request) bool {
	number := int(*req.Number)
	if float64(number) != *req.Number {
		return false
	}
	if number <= 1 {
		return false
	}
	if number == 2 {
		return true
	}
	for i := 2; i*i <= number; i++ {
		if number%i == 0 {
			return false
		}
	}
	return true
}

func checkJSON(data []byte) (Request, bool) {
	var req Request
	err := json.Unmarshal(data, &req)
	if err != nil {
		return req, false
	}
	if req.Method == nil || req.Number == nil {
		return req, false
	}
	if *req.Method != "isPrime" {
		return req, false
	}

	return req, true
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		reqBytes, err := reader.ReadBytes('\n')
		fmt.Println("Raw : ", string(reqBytes))
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			return
		}

		reqBytes = bytes.TrimRight(reqBytes, "\n")
		req, ok := checkJSON(reqBytes)
		fmt.Println("Request:", req)
		if !ok {
			_, err = conn.Write([]byte(BadResponseJSON + "\n"))
			if err != nil {
				fmt.Println("Error writing:", err.Error())
			}
			fmt.Println("Bad request")
			return
		} else {
			var res Response
			if checkPrime(req) {
				res = Response{"isPrime", true}
			} else {
				res = Response{"isPrime", false}
			}
			fmt.Println("Response:", res)
			resJSON, err := json.Marshal(res)
			if err != nil {
				fmt.Println("Error marshalling response:", err.Error())
				return
			}
			_, err = conn.Write(append(resJSON, '\n'))
			if err != nil {
				fmt.Println("Error writing:", err.Error())
				return
			}
		}
	}
}

func Run() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Listening on port 8080")

	for {
		connection, err := listen.Accept()
		if err != nil {
			fmt.Println("Error: ", err.Error())
			return
		}
		fmt.Println("Accepted connection.")
		go handleRequest(connection)
	}
}
