package internal

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func HandleIncomingClient(conn net.Conn, clients map[string]Client) {
	// mu := sync.Mutex{}
	// mu.Lock()
	// if len(clients) == 2 {
	// 	fmt.Fprintln(conn, "Chatroom is full, connect later")
	// 	// mu.Unlock()
	// 	// conn.Close()
	// 	return
	// }
	// mu.Unlock()
	// creates new client
	newClient := Client{address: conn.RemoteAddr().String(), conn: conn}
	fmt.Fprint(conn, "[ENTER YOUR NAME]: ")
	name, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	newClient.name = name[:len(name)-1]
	clients[newClient.address] = newClient
	// send message to inform about join
	messages <- newMessage(newClient.name+" has joined our chat...", newClient.address)
	input := bufio.NewScanner(conn)
	newClient.PrintPrefix(conn)
	for input.Scan() {
		newClient.PrintPrefix(conn)
		messages <- newMessage(getPrefix(newClient)+input.Text(), newClient.address)
	}

	delete(clients, conn.RemoteAddr().String())
	conn.Close()
	leaving <- newMessage(newClient.name+" has left our chat...", newClient.address)
}

func (c *Client) PrintPrefix(conn net.Conn) {
	fmt.Fprintf(conn, getPrefix(*c))
}

func newMessage(text, from string) message {
	return message{
		text: text,
		from: from,
	}
}
