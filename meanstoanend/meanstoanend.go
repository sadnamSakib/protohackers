package meanstoanend

import (
	"encoding/binary"
	"fmt"
	"io"
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
func (m Message) getValues() (int32, int32) {

	firstVal := int32(binary.BigEndian.Uint32(m[1:5]))
	secondVal := int32(binary.BigEndian.Uint32(m[5:]))
	return firstVal, secondVal

}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	d := make([]Data, 0)
	for {
		m := make(Message, 9)
		if _, err := io.ReadFull(conn, m); err != nil {
			break
		}
		if m.getMethod() == INSERT {

			firstVal, secondVal := m.getValues()
			d = append(d, Data{(firstVal), (secondVal)})

		} else if m.getMethod() == QUERY {

			firstVal, secondVal := m.getValues()
			mean := 0
			num := 0
			for _, data := range d {
				if data.Timestamp >= firstVal && data.Timestamp <= int32(secondVal) {
					mean += int(data.Price)
					num += 1
					fmt.Println("Data:", data.Timestamp, data.Price)
					fmt.Println("Mean:", mean, num)
				}
			}
			fmt.Printf("Final result: %d/%d", mean, num)
			if num <= 0 {
				mean = 0
			} else {
				mean = mean / num
			}

			b := make([]byte, 4)
			binary.BigEndian.PutUint32(b, uint32(mean))
			if _, err := conn.Write(b); err != nil {
				break
			}

		} else {
			return
		}

	}

}

func Run() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		connection, err := listen.Accept()
		if err != nil {
			fmt.Println("Error: ", err.Error())
			return
		}
		go handleRequest(connection)
	}
}
