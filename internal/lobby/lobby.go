package lobby

import (
	"bufio"
	"log"
	"net"
	"net-cat/internal/chatroom"
	"net-cat/internal/service"
	"net-cat/internal/user"
	i "net-cat/internal/userInterface"
	"sync"
)

var LOGO = ""

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
	lobby.PrintLogo(conn)
	username := lobby.AskName(conn)
	user := user.NewUser(username, conn)
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
	conn.Write([]byte(LOGO + "\n"))
	lobby.mu.Unlock()
}

func (lobby *Lobby) AskName(conn net.Conn) string {
	conn.Write([]byte("Enter your name:"))
	name, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	if name[:len(name)-1] == "" {
		conn.Write([]byte("name cannot be empty\n"))
		return lobby.AskName(conn)
	}
	if lobby.UserExist(name) {
		conn.Write([]byte("name has been taken\n"))
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

func (lobby *Lobby) DeleteChatroom(name string) {
	delete(lobby.rooms, name)
}

func (lobby *Lobby) GetChatroom(name string) Chatroom {
	return lobby.rooms[name]
}
