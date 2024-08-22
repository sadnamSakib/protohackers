package budgetchat

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

type Client struct {
	name string
	conn net.Conn
}
type Message struct {
	Name string
	Text string
}

func nameResolution(conn net.Conn, connections map[string]Client) (string, bool) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return "", false
	}
	clientName := string(buf[:n])
	clientName = strings.TrimSpace(clientName)

	if len(clientName) == 0 || !isValidName(clientName) {
		conn.Write([]byte("Invalid name.\n"))
		conn.Close()
		return "", false
	}

	otherMembers := make([]string, 0, len(connections))
	for key := range connections {
		otherMembers = append(otherMembers, key)
	}

	roomMembersMessage := "* The room contains: " + strings.Join(otherMembers, ", ") + "\n"
	conn.Write([]byte(roomMembersMessage))

	return clientName, true
}

func isValidName(name string) bool {
	pattern := `^[[:alnum:]]{1,16}$`
	r := regexp.MustCompile(pattern)
	return r.MatchString(name)
}

func handleClient(client Client, join chan<- Client, leave chan<- Client, broadcast chan<- Message, clients map[string]Client) {
	defer func() {
		leave <- client
		client.conn.Close()
	}()

	client.conn.Write([]byte("Welcome to budgetchat! What shall I call you?\n"))

	name, ok := nameResolution(client.conn, clients)
	if !ok {
		return
	}

	client.name = name
	join <- client

	for {
		buf := make([]byte, 1024)
		n, err := client.conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		clientMessage := strings.TrimSpace(string(buf[:n]))
		if clientMessage != "" {
			broadcast <- Message{
				Name: client.name,
				Text: fmt.Sprintf("[%s] %s\n", client.name, clientMessage),
			}
		} else {
			break
		}
	}
}

func manageClients(join <-chan Client, leave <-chan Client, broadcast <-chan Message, clients map[string]Client) {

	for {
		select {
		case client := <-join:
			clients[client.name] = client
			message := fmt.Sprintf("* %s has entered the room\n", client.name)
			for _, c := range clients {
				if c.name != client.name {
					c.conn.Write([]byte(message))
				}
			}
			fmt.Print(message)

		case client := <-leave:
			delete(clients, client.name)
			message := fmt.Sprintf("* %s has left the room\n", client.name)
			for _, c := range clients {
				c.conn.Write([]byte(message))
			}
			fmt.Print(message)

		case message := <-broadcast:
			for _, client := range clients {
				if client.name != message.Name {
					client.conn.Write([]byte(message.Text))
				}
			}
		}
	}
}

func Run() {
	listen, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listen.Close()

	fmt.Println("- listening on port 8000")

	join := make(chan Client)
	leave := make(chan Client)
	broadcast := make(chan Message)
	clients := make(map[string]Client)
	go manageClients(join, leave, broadcast, clients)

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		client := Client{conn: conn}
		go handleClient(client, join, leave, broadcast, clients)
	}
}
