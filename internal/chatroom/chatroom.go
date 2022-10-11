package chatroom

import (
	"fmt"
	"net-cat/internal/service"
	i "net-cat/internal/userInterface"
)

const (
	INFO_LEAVE = " has left the chatroom"
	INFO_JOIN  = " has joined the chatroom"
)

// TODO:unit test
type Chatroom struct {
	name  string
	users map[string]i.User
	log   string
}

func NewChatroom(name string) *Chatroom {
	chatroom := &Chatroom{name: name, users: map[string]i.User{}, log: ""}
	return chatroom
}

func (room *Chatroom) GetChatName() string {
	return room.name
}

func (room *Chatroom) GetUsers() map[string]i.User {
	return room.users
}

func (room *Chatroom) AddUser(user i.User) {
	room.broadcastInfo(INFO_JOIN, user.GetName())
	room.users[user.GetName()] = user
	user.SetRoomName(room.name)
	fmt.Fprint(user.GetConn(), service.GetPrefix(user.GetName()))
}

func (room *Chatroom) IsFull() bool {
	return len(room.users) == 10
}

func (room *Chatroom) DeleteUser(user i.User) {
	delete(room.users, user.GetName())
	user.SetRoomName("")
	room.broadcastInfo(INFO_LEAVE, user.GetName())
}

func (room *Chatroom) broadcastInfo(info, name string) {
	for _, user := range room.users {
		fmt.Fprintln(user.GetConn(), "")
		fmt.Fprintln(user.GetConn(), name+info)
		fmt.Fprint(user.GetConn(), service.GetPrefix(user.GetName()))
	}
}

func (room *Chatroom) ListUsers(user i.User) {
	for name := range room.users {
		fmt.Fprintln(user.GetConn(), name)
	}
	fmt.Fprint(user.GetConn(), service.GetPrefix(user.GetName()))
}

func (room *Chatroom) LogMessage(message string) {
	room.log += message + "\n"
}

func (room *Chatroom) DisplayLog(user i.User) {
	fmt.Fprint(user.GetConn(), room.log)
}
