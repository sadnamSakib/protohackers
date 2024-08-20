package budgetchat

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

func nameResolution(conn *net.Conn, connections map[string]*net.Conn) string {
	buf := make([]byte, 1024)
	n, err := (*conn).Read(buf)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	clientName := string(buf[:n])
	clientName, _ = strings.CutSuffix(clientName, "\n")
	fmt.Println("Client : ", clientName)
	pattern := `^[a-zA-Z0-9]+$`
	r := regexp.MustCompile(pattern)
	if !r.MatchString(clientName) {
		(*conn).Write([]byte("Invalid name."))
		(*conn).Close()
		return ""
	} else {
		otherMemebers := make([]string, 0, len(connections))
		for key := range connections {
			otherMemebers = append(otherMemebers, key)
		}

		roomMembersMessage := "* The room contains: " + strings.Join(otherMemebers, ", ") + "\n"
		for _, roomMembers := range connections {
			(*roomMembers).Write([]byte(fmt.Sprintf("* %s has entered the room\n", clientName)))
		}
		fmt.Printf("Server: * %v has entered the room\n", clientName)
		connections[clientName] = conn
		(*conn).Write([]byte(roomMembersMessage))
		return clientName
	}

}

func handleRequest(conn *net.Conn, connections map[string]*net.Conn) {
	serverMessage := "Welcome to budgetchat! What shall I call you?\n"
	(*conn).Write([]byte(serverMessage))
	fmt.Println("Server : ", serverMessage)
	name := nameResolution(conn, connections)
	defer func() {
		delete(connections, name)
		for _, roomMembers := range connections {
			(*roomMembers).Write([]byte(fmt.Sprintf("* %s has left the room\n", name)))
		}
		fmt.Printf("Server: * %v has left the room\n", name)
		(*conn).Close()
	}()

	for {
		buf := make([]byte, 1024)
		n, err := (*conn).Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		clientMessage := string(buf[:n])
		fmt.Println("Client: ", clientMessage)
		fmt.Printf("[%v] %v\n", name, clientMessage)
		if clientMessage != "" {
			for otherClients, roomMembers := range connections {
				if name != otherClients {
					(*roomMembers).Write([]byte(fmt.Sprintf("[%s] %s\n", name, clientMessage)))
				}
			}
		} else {

			break
		}
	}

}

func Run() {
	listen, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("- listening on port 8000")
	connections := make(map[string]*net.Conn)
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go handleRequest(&conn, connections)
	}

}
