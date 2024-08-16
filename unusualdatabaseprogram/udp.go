package unusualdatabaseprogram

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

func process(s string, sm *sync.Map) string {
	fmt.Println("Processing ", s)
	if strings.Contains(s, "=") {
		split := strings.Split(s, "=")
		(*sm).Store(split[0], split[1])
		return "i"
	} else {
		value, ok := (*sm).Load(s)
		if ok {
			return fmt.Sprintf("%s=%s", s, value)
		} else {
			return fmt.Sprintf("%s=", s)
		}
	}
}

func handleRequest(conn *net.UDPConn, sm *sync.Map) {
	fmt.Println("Map before : ", *sm)
	buf := make([]byte, 1024)
	n, addr, err := conn.ReadFromUDP(buf)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("Received ", string(buf[0:n]), " from ", addr)
	ret := process(string(buf[0:n]), sm)
	if ret == "i" {
		fmt.Println("Stored")
	} else {
		_, err = conn.WriteToUDP([]byte(ret), addr)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Println("Sent ", ret, " to ", addr)
	}
	fmt.Println("Map after : ", *sm)
}

func Run() {

	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: 8080})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Connecting to UDP server on port 80")
	defer conn.Close()
	sm := new(sync.Map)
	for {
		handleRequest(conn, sm)
	}

}
