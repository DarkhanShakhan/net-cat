package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const (
	CONN_TYPE = "tcp"
	USAGE     = "[USAGE]: ./TCPChat $port"

	INFO_LEAVE = " has left the chatroom"
	INFO_JOIN  = " has joined the chatroom"
)

const (
	CMD        = "/"
	CMD_JOIN   = CMD + "join"
	CMD_CREATE = CMD + "create"
	CMD_LIST   = CMD + "list"
	CMD_LEAVE  = CMD + "leave"
	CMD_USERS  = CMD + "users"
)

var CONN_PORT = "8989"

func main() {
	if len(os.Args) > 2 {
		fmt.Println(USAGE)
		return
	} else if len(os.Args) == 2 {
		CONN_PORT = os.Args[1]
	}
	//creates listener
	listener, err := net.Listen(CONN_TYPE, "localhost:"+CONN_PORT)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Printf("Listening on the port %s\n", CONN_PORT)
	lobby := newLobby()
	//broadcast messages
	go lobby.parseSignal()
	//starts accepting connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go lobby.handleClient(conn)
	}
}

func askName(conn net.Conn) string {
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
		return askName(conn)
	}
	return name[:len(name)-1]
}

//lobby

type Lobby struct {
	rooms       map[string]*Chatroom
	msgChannel  chan Message
	cmdChannel  chan Command
	infoChannel chan Info
}

func newLobby() *Lobby {
	return &Lobby{rooms: map[string]*Chatroom{}, msgChannel: make(chan Message), cmdChannel: make(chan Command)}
}

//TODO:
// * deal with creating Chat
// * deal with joining Chat
func (lobby *Lobby) handleClient(conn net.Conn) {
	defer conn.Close()
	fmt.Fprintln(conn, "Welcome to TCP chat!!!")
	client := newClient(askName(conn), conn)
	flow := bufio.NewScanner(client.conn)
	for flow.Scan() {
		signal := flow.Text()
		lobby.sendSignal(signal, client)
	}
}

func (lobby *Lobby) parseSignal() {
	select {
	case cmd := <-lobby.cmdChannel:
		lobby.parseCommand(cmd)
	case msg := <-lobby.msgChannel:
		lobby.broadcastMsg(msg)
	}
}

func (lobby *Lobby) listChats(client *Client) {
	for key := range lobby.rooms {
		fmt.Fprintln(client.conn, key)
	}
}

//TODO: What other commands? join chatname, create chatname
func (lobby *Lobby) parseCommand(cmd Command) {
	switch cmd.command {
	case CMD_LIST:
		lobby.listChats(cmd.client)
	case CMD_USERS:
		cmd.client.chatroom.listUsers(cmd.client)
	case CMD_LEAVE:
		cmd.client.chatroom.deleteClient(cmd.client)
	}
}

//TODO: errors to deal: writing command into lobby or chatroom, writing message into lobby
func (lobby *Lobby) broadcastMsg(msg Message) {
	for key, otherClient := range msg.client.chatroom.clients {
		if key != msg.client.conn.RemoteAddr().String() {
			fmt.Fprintln(otherClient.conn, msg.text)
		}
	}
	msg.client.chatroom.logMessage(msg.text)
}

//TODO: refactor the code
func (lobby *Lobby) broadcastInfo(info Info) {
	for key, otherClient := range info.client.chatroom.clients {
		if key != info.client.conn.RemoteAddr().String() {
			fmt.Fprintln(otherClient.conn, info.client.name+info.text)
		}
	}
}

//TODO: add mutexes
func (lobby *Lobby) sendSignal(signal string, client *Client) {
	switch {
	case strings.HasPrefix(signal, CMD):
		lobby.cmdChannel <- newCommand(signal, client)
	default:
		lobby.msgChannel <- newMessage(signal, client)
	}
}

//TODO: add welcome logo, prefixes
//chatroom
type Chatroom struct {
	clients map[string]*Client
	log     string
}

func newChatroom() *Chatroom {
	return &Chatroom{}
}

func (room *Chatroom) addClient(client *Client) {
	room.broadcastInfo(INFO_JOIN, client.name)
	room.clients[client.conn.RemoteAddr().String()] = client
}

func (room *Chatroom) isFull() bool {
	return len(room.clients) == 10
}

func (room *Chatroom) deleteClient(client *Client) {
	delete(room.clients, client.conn.RemoteAddr().String())
	room.broadcastInfo(INFO_LEAVE, client.name)
}

func (room *Chatroom) broadcastInfo(info, name string) {
	for _, otherClient := range room.clients {
		fmt.Fprintln(otherClient.conn, name+info)
	}
}
func (room *Chatroom) listUsers(client *Client) {
	for key, otherClient := range room.clients {
		if key != client.conn.RemoteAddr().String() {
			fmt.Fprintln(client.conn, otherClient.name)
		}
	}
}

func (room *Chatroom) logMessage(message string) {
	room.log += message
}

func (room *Chatroom) displayLog(client *Client) {
	fmt.Fprintln(client.conn, room.log)
}

//client
type Client struct {
	name     string
	chatroom *Chatroom
	conn     net.Conn
}

func newClient(name string, conn net.Conn) *Client {
	return &Client{name: name, conn: conn}
}

//message
type Message struct {
	text   string
	client *Client
}

func newMessage(text string, client *Client) Message {
	return Message{text: text, client: client}
}

//command
type Command struct {
	command string
	name    string
	client  *Client
}

//TODO: add error to check command e.g.join chatname, create chatname etc. if more than two agruments, return error
func newCommand(command string, client *Client) Command {
	return Command{command: command, client: client}
}

//info
type Info struct {
	text   string
	client *Client
}

func newInfo(text string, client *Client) Info {
	return Info{text: text, client: client}
}
