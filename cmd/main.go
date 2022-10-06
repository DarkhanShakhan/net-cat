package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	CONN_TYPE     = "tcp"
	USAGE         = "[USAGE]: ./TCPChat $port"
	TIME_FORMAT   = "2006-01-02 15:04:05"
	INFO_LEAVE    = " has left the chatroom"
	INFO_JOIN     = " has joined the chatroom"
	LOGO_FILENAME = "cmd/logo.txt"
)

const (
	CMD        = "/"
	CMD_JOIN   = CMD + "join"
	CMD_CREATE = CMD + "create"
	CMD_LIST   = CMD + "list"
	CMD_LEAVE  = CMD + "leave"
	CMD_USERS  = CMD + "users"
	CMD_HELP   = CMD + "help"
	CMD_DIRECT = CMD + "direct"
)

var (
	CONN_PORT = "8989"
	LOGO      = ""
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println(USAGE)
		return
	} else if len(os.Args) == 2 {
		CONN_PORT = os.Args[1]
	}
	// creates listener
	listener, err := net.Listen(CONN_TYPE, "localhost:"+CONN_PORT)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Printf("Listening on the port %s\n", CONN_PORT)
	lobby := newLobby()
	// broadcast messages
	go lobby.parseSignal()
	// starts accepting connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go lobby.handleClient(conn)
	}
}

func (lobby *Lobby) askName(conn net.Conn) string {
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
		return lobby.askName(conn)
	}
	if _, ok := lobby.users[name[:len(name)-1]]; ok {
		fmt.Fprintln(conn, "name has been taken")
		return lobby.askName(conn)
	}
	return name[:len(name)-1]
}

// lobby

type Lobby struct {
	rooms       map[string]*Chatroom
	users       map[string]*Client
	msgChannel  chan Message
	cmdChannel  chan Command
	infoChannel chan Info
}

func newLobby() *Lobby {
	return &Lobby{rooms: map[string]*Chatroom{}, users: map[string]*Client{}, msgChannel: make(chan Message), cmdChannel: make(chan Command)}
}

func (lobby *Lobby) handleClient(conn net.Conn) {
	defer conn.Close()
	if LOGO == "" {
		LOGO = parseLogo()
	}
	fmt.Fprintln(conn, LOGO)

	client := newClient(lobby.askName(conn), conn)
	lobby.users[client.name] = client
	flow := bufio.NewScanner(client.conn)
	for flow.Scan() {
		signal := flow.Text()
		lobby.sendSignal(signal, client)
	}
	if client.chatroom != nil {
		client.chatroom.deleteClient(client)
	}
	delete(lobby.users, client.name)
}

func parseLogo() string {
	data, err := os.ReadFile(LOGO_FILENAME)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func (lobby *Lobby) parseSignal() {
	for {
		select {
		case cmd := <-lobby.cmdChannel:
			lobby.parseCommand(cmd)
		case msg := <-lobby.msgChannel:
			if msg.client.chatroom == nil {
				fmt.Fprintln(msg.client.conn, "use commands in the lobby, starting with '/'\nyou can look all commands with '/help")
			} else {
				lobby.broadcastMsg(msg)
			}
		}
	}
}

func (lobby *Lobby) listChats(client *Client) {
	for key := range lobby.rooms {
		fmt.Fprintln(client.conn, key)
	}
}

// TODO: What other commands? direct message
func (lobby *Lobby) parseCommand(cmd Command) {
	switch cmd.command {
	case CMD_LIST:
		lobby.listChats(cmd.client)
	case CMD_USERS:
		if cmd.client.chatroom != nil {
			cmd.client.chatroom.listUsers(cmd.client)
		} else {
			lobby.listUsers(cmd.client)
		}
	case CMD_DIRECT:
		to := lobby.users[cmd.name]
		if cmd.client.chatroom != nil && cmd.client.chatroom == to.chatroom {
			fmt.Fprintln(to.conn)
			fmt.Fprintln(to.conn, "[DIRECT]"+cmd.message.prefix+cmd.message.text)
			fmt.Fprint(cmd.client.conn, cmd.message.prefix)
			fmt.Fprint(to.conn, cmd.message.prefix)
		}
	case CMD_LEAVE:
		cmd.client.chatroom.deleteClient(cmd.client)
	case CMD_JOIN:
		if lobby.rooms[cmd.name].isFull() {
			fmt.Fprintln(cmd.client.conn, "The Chat is full, join later or create a new one")
		} else {
			lobby.rooms[cmd.name].addClient(cmd.client)
		}
	case CMD_CREATE:
		lobby.createChatroom(cmd.client, cmd.name)
	case CMD_HELP:
		lobby.listCommands(cmd.client)
	}
}

func (lobby *Lobby) listUsers(client *Client) {
	for name, user := range lobby.users {
		fmt.Fprint(client.conn, name)
		if user.chatroom != nil {
			fmt.Fprint(client.conn, " -- "+user.chatroom.name)
		}
		fmt.Fprintln(client.conn, "")
	}
}

// TODO: errors with command and arguments

func (lobby *Lobby) createChatroom(client *Client, name string) {
	if _, ok := lobby.rooms[name]; ok {
		fmt.Fprintln(client.conn, "The chat with a given name exists")
	} else {
		lobby.rooms[name] = newChatroom(name, client)
	}
}

func (lobby *Lobby) listCommands(client *Client) {
	if client.chatroom == nil {
		fmt.Fprintln(client.conn, CMD_LIST+" -> display all chatrooms.")
		fmt.Fprintln(client.conn, CMD_CREATE+" roomName -> create a new chatroom with a given room name.")
		fmt.Fprintln(client.conn, CMD_JOIN+" roomName -> join the given chatroom.")
		fmt.Fprintln(client.conn, CMD_USERS+" -> list all users.")
	} else {
		fmt.Fprintln(client.conn, CMD_USERS+" -> list all users in the chat.")
		fmt.Fprintln(client.conn, CMD_LEAVE+" -> leave the chatroom to the lobby.")
		fmt.Fprint(client.conn, getPrefix(client.name))
	}
}

func (lobby *Lobby) broadcastMsg(msg Message) {
	for key, otherClient := range msg.client.chatroom.clients {
		if key != msg.client.name {
			fmt.Fprintln(otherClient.conn, "")
			fmt.Fprintln(otherClient.conn, msg.prefix+msg.text)
			fmt.Fprint(otherClient.conn, getPrefix(otherClient.name))
		} else {
			fmt.Fprint(otherClient.conn, msg.prefix)
		}
	}
	msg.client.chatroom.logMessage(msg.prefix + msg.text)
}

// TODO: refactor the code
// func (lobby *Lobby) broadcastInfo(info Info) {
// 	for key, otherClient := range info.client.chatroom.clients {
// 		if key != info.client.conn.RemoteAddr().String() {
// 			fmt.Fprintln(otherClient.conn, "")
// 			fmt.Fprintln(otherClient.conn, info.client.name+info.text)
// 		}
// 	}
// }

// TODO: add mutexes
func (lobby *Lobby) sendSignal(signal string, client *Client) {
	switch {
	case strings.HasPrefix(signal, CMD):
		lobby.sendCommand(signal, client)
	default:
		if signal != "" {
			lobby.msgChannel <- newMessage(signal, client)
		}
	}
}

func (lobby *Lobby) sendCommand(command string, client *Client) {
	temp := strings.Split(command, " ")

	switch len(temp) {
	case 1:
		lobby.cmdChannel <- newCommand(command, "", client, Message{})
	case 2:
		lobby.cmdChannel <- newCommand(temp[0], temp[1], client, Message{})
	case 3:
		lobby.cmdChannel <- newCommand(temp[0], temp[1], client, newMessage(temp[2], client))
	}
}

type Chatroom struct {
	name    string
	clients map[string]*Client
	log     string
}

func newChatroom(name string, client *Client) *Chatroom {
	chatroom := &Chatroom{name: name, clients: map[string]*Client{}, log: ""}
	chatroom.clients[client.name] = client
	client.chatroom = chatroom
	fmt.Fprint(client.conn, getPrefix(client.name))
	return chatroom
}

func (room *Chatroom) addClient(client *Client) {
	room.broadcastInfo(INFO_JOIN, client.name)
	room.clients[client.name] = client
	client.joinChatroom(room)
}

func (room *Chatroom) isFull() bool {
	return len(room.clients) == 10
}

func (room *Chatroom) deleteClient(client *Client) {
	delete(room.clients, client.name)
	client.leaveChatroom()
	room.broadcastInfo(INFO_LEAVE, client.name)
}

func (room *Chatroom) broadcastInfo(info, name string) {
	for _, otherClient := range room.clients {
		fmt.Fprintln(otherClient.conn, "")
		fmt.Fprintln(otherClient.conn, name+info)
		fmt.Fprint(otherClient.conn, getPrefix(otherClient.name))
	}
}

func (room *Chatroom) listUsers(client *Client) {
	for name := range room.clients {
		fmt.Fprintln(client.conn, name)
	}
	fmt.Fprint(client.conn, getPrefix(client.name))
}

func (room *Chatroom) logMessage(message string) {
	room.log += message + "\n"
}

func (room *Chatroom) displayLog(client *Client) {
	fmt.Fprint(client.conn, room.log)
}

// client
type Client struct {
	name     string
	chatroom *Chatroom
	conn     net.Conn
}

func newClient(name string, conn net.Conn) *Client {
	return &Client{name: name, conn: conn}
}

func (client *Client) joinChatroom(chatroom *Chatroom) {
	client.chatroom = chatroom
	chatroom.displayLog(client)
	fmt.Fprint(client.conn, getPrefix(client.name))
}

func (client *Client) leaveChatroom() {
	client.chatroom = nil
}

// message
type Message struct {
	text   string
	client *Client
	prefix string
}

func newMessage(text string, client *Client) Message {
	return Message{text: text, prefix: getPrefix(client.name), client: client}
}

func getPrefix(name string) string {
	return fmt.Sprintf("[%s][%s]:", time.Now().Format(TIME_FORMAT), name)
}

// command
type Command struct {
	command string
	name    string
	message Message
	client  *Client
}

func newCommand(command string, name string, client *Client, message Message) Command {
	return Command{command: command, name: name, client: client, message: message}
}

// info
type Info struct {
	text   string
	client *Client
}

func newInfo(text string, client *Client) Info {
	return Info{text: text, client: client}
}
