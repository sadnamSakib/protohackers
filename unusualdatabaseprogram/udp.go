package unusualdatabaseprogram

import (
	"fmt"
	"net"
	"strings"
)

func process(s string, sm map[string]string) string {
	s, _ = strings.CutSuffix(s, "\n")
	if strings.Contains(s, "=") {
		split := strings.Split(s, "=")
		if split[0] != "version" {

			(sm)[split[0]] = split[1]
		} else {
			return "ignore"
		}
		return "i"
	} else {
		return fmt.Sprintf("%s=%s", s, sm[s])
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
		fmt.Println("Inserted")
		return
	} else if ret == "ignore" {
		fmt.Println("Ignored")
		return
	} else {

		_, err = conn.WriteToUDP([]byte(ret), addr)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Println("Response : ", string([]byte(ret)))

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
	sm := map[string]string{
		"version": "Ken's Key-Value Store 1.0",
	}
	for {
		handleRequest(conn, sm)
	}

}
