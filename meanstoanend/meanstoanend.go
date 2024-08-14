package meanstoanend

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Message []byte

const (
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

	firstVal := int(binary.BigEndian.Uint32(m[1:5]))
	secondVal := int(binary.BigEndian.Uint32(m[5:9]))
	return firstVal, secondVal

}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	m := make(Message, 9)
	for i := range m {
		rawData := make([]byte, 1)
		_, err := conn.Read(rawData)
		if err != nil {
			fmt.Println("Error: ", err.Error())
			return
		}
		fmt.Println("Received data: ", rawData)
		fmt.Println("Received data: ", string(rawData))
		m[i] = rawData[0]
	}

	if m.getMethod() == INSERT {
		fmt.Println("Inserting: ")
		firstVal, secondVal := m.getValues()
		fmt.Println("First value: ", firstVal)
		fmt.Println("Second value: ", secondVal)
	} else {
		fmt.Println("Querying: ")
		firstVal, secondVal := m.getValues()
		fmt.Println("First value: ", firstVal)
		fmt.Println("Second value: ", secondVal)
	}

	fmt.Println("Received message: ", string(m))

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
