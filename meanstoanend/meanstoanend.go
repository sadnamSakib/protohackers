package meanstoanend

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Message []byte
type Data struct {
	Timestamp int32
	Price     int32
}

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

func handleRequest(conn net.Conn, d *[]Data) {
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
		*d = append(*d, Data{int32(firstVal), int32(secondVal)})
		fmt.Println(d)
	} else {
		fmt.Println("Querying: ")
		firstVal, secondVal := m.getValues()
		mean := 0
		num := 0
		for _, data := range *d {
			if data.Timestamp >= int32(firstVal) && data.Timestamp <= int32(secondVal) {
				fmt.Println(data)
				mean += int(data.Price)
				num++
			}
		}
		if num == 0 {
			mean = 0
		} else {
			mean = mean / num
		}
		fmt.Println("Mean: ", mean)
		conn.Write([]byte{byte(mean >> 24), byte(mean >> 16), byte(mean >> 8), byte(mean)})

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
	d := make([]Data, 0)

	for {
		connection, err := listen.Accept()
		if err != nil {
			fmt.Println("Error: ", err.Error())
			return
		}
		fmt.Println("Accepted connection.")
		go handleRequest(connection, &d)
	}
}
