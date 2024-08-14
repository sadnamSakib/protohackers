package smoketest

import (
	"fmt"
	"net"
)

func handleRequest(conn net.Conn) {
	defer conn.Close()

	for {
		buf := make([]byte, 1024)
		reqLen, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			return
		}
		fmt.Println("Received data:", string(buf[:reqLen]))
		conn.Write([]byte(buf[:reqLen]))
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
