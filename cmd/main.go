package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	i "net-cat/internal"
	"os"
)

const (
	CONN_TYPE = "tcp"
	USAGE     = "[USAGE]: ./TCPChat $port"
)

var CONN_PORT = "8989"

func main() {
	if len(os.Args) > 2 {
		fmt.Println(USAGE)
		return
	} else if len(os.Args) == 2 {
		CONN_PORT = os.Args[1]
	}
	listener, err := net.Listen(CONN_TYPE, "localhost:"+CONN_PORT)
	if err != nil {
		log.Fatal(err)
	}
	// defer listener.Close()
	fmt.Printf("Listening on the port %s\n", CONN_PORT)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go welcomeNewClient(conn)
	}
}

func welcomeNewClient(conn net.Conn) {
	defer conn.Close()
	name := askClientName(conn)
	NewClient(name, conn)
}

//TODO: what to with client?
//give interface to chat
func askClientName(conn net.Conn) string {
	_, err := fmt.Fprint(conn, "Enter your name: ")
	if err != nil {
		log.Fatal(err)
	}
	name, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	if name[:len(name)-1] == "" {
		fmt.Fprintln(conn, "name cannot be empty")
		return askClientName(conn)
	}
	return name[:len(name)-1]
}

type Chatroom struct {
}

type Client struct {
	name     string
	chatroom *Chatroom
	conn     net.Conn
	outgoing chan *Message
}

func NewClient(name string, conn net.Conn) *Client {
	return &Client{name: name, conn: conn}
}

type Message struct {
	text   string
	client *Client
}

var clients = make(map[string]i.Client)
