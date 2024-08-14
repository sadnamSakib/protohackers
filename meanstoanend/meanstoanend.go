package meanstoanend

import (
	"fmt"
	"net"
)

type Message []byte

const (
	// Message types
	INSERT = iota
	QUERY
)

func (m Message) getMethod() int {
	if rune(m[0]) == 'I' {
		return INSERT
	} else {
		return QUERY
	}
}
func (m Message) getValues() (int, int) {
	firstVal := 0
	for i := 5; i >= 1; i-- {
		firstVal = firstVal*10 + int(m[i]-'0')
	}
	return 0, 0
}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	m := make(Message, 9)
	_, err := conn.Read(m)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return
	}
	fmt.Println("Received message: ", m)

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
