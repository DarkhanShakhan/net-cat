package lobby

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net-cat/internal/chatroom"
	"net-cat/internal/service"
	"net-cat/internal/user"
	i "net-cat/internal/userInterface"
	"sync"
)

var LOGO = ""

// TODO: handle errors
// TODO: unit tests
//TODO: mutexes
type Lobby struct {
	mu         sync.Mutex
	rooms      map[string]Chatroom
	users      map[string]i.User
	msgChannel chan Message
	cmdChannel chan Command
}

func NewLobby() *Lobby {
	return &Lobby{rooms: map[string]Chatroom{}, users: map[string]i.User{}, msgChannel: make(chan Message), cmdChannel: make(chan Command)}
}

func (lobby *Lobby) HandleUser(conn net.Conn) {
	defer conn.Close()
	// if LOGO == "" {
	// 	LOGO = service.ParseLogo()
	// }
	// fmt.Fprintln(conn, LOGO)
	lobby.PrintLogo(conn)
	user := user.NewUser(lobby.AskName(conn), conn)
	// lobby.users[user.GetName()] = user
	lobby.AddUser(user)
	flow := bufio.NewScanner(conn)
	for flow.Scan() {
		signal := flow.Text()
		lobby.SendSignal(signal, user)
	}
	lobby.mu.Lock()
	if name, ok := user.GetRoomName(); ok {
		lobby.GetChatroom(name).DeleteUser(user)
	}
	delete(lobby.users, user.GetName())
	lobby.mu.Unlock()
}

func (lobby *Lobby) AddUser(user i.User) {
	lobby.mu.Lock()
	lobby.users[user.GetName()] = user
	lobby.mu.Unlock()
}

func (lobby *Lobby) PrintLogo(conn net.Conn) {
	lobby.mu.Lock()
	if LOGO == "" {
		LOGO = service.ParseLogo()
	}
	fmt.Fprintln(conn, LOGO)
	lobby.mu.Unlock()
}

func (lobby *Lobby) AskName(conn net.Conn) string {
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
		return lobby.AskName(conn)
	}
	if lobby.UserExist(name) {
		fmt.Fprintln(conn, "name has been taken")
		return lobby.AskName(conn)
	}
	return name[:len(name)-1]
}

func (lobby *Lobby) UserExist(name string) bool {
	lobby.mu.Lock()
	defer lobby.mu.Unlock()
	_, ok := lobby.users[name[:len(name)-1]]
	return ok
}

func (lobby *Lobby) CreateChatroom(name string) bool {
	if _, ok := lobby.rooms[name]; ok {
		return false
	}
	lobby.rooms[name] = chatroom.NewChatroom(name)
	return true
}

func (lobby *Lobby) GetChatroom(name string) Chatroom {
	return lobby.rooms[name]
}
