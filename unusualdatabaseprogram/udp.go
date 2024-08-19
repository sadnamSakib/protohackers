package unusualdatabaseprogram

import (
	"fmt"
	"net"
	"strings"
)

func process(s string, sm map[string]string) string {
	s, _ = strings.CutSuffix(s, "\n")
	if s == "version" {
		return "version=1.0"
	} else if strings.Contains(s, "=") {
		split := strings.Split(s, "=")
		if split[0] == "version" {
			return "version="
		}
		(sm)[split[0]] = split[1]
		return "i"
	} else {
		value, ok := sm[s]
		if ok {
			return fmt.Sprintf("%s=%s", s, value)
		} else {
			return fmt.Sprintf("%s=", s)
		}
	}
}

func handleRequest(conn *net.UDPConn, sm map[string]string) {

	buf := make([]byte, 1024)
	n, addr, err := conn.ReadFromUDP(buf)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	ret := process(string(buf[0:n]), sm)
	if ret == "i" {
		fmt.Println("Stored")

	} else {

		_, err = conn.WriteToUDP([]byte(ret), addr)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

	}

}

func Run() {

	port := 8080
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("Connecting to UDP  on port 8080")
	defer conn.Close()
	sm := make(map[string]string)
	for {
		handleRequest(conn, sm)
	}

}
