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
	fmt.Fprintln(user.GetConn(), "Joining "+room.name)
	fmt.Fprintln(user.GetConn())
	room.DisplayLog(user)

	user.GetConn().Write([]byte(service.GetPrefix(user.GetName())))
}

func (room *Chatroom) IsFull() bool {
	return len(room.users) == 10
}

func (room *Chatroom) IsEmpty() bool {
	return len(room.users) == 0
}

func (room *Chatroom) DeleteUser(user i.User) {
	delete(room.users, user.GetName())
	user.SetRoomName("")
	fmt.Fprintln(user.GetConn())
	fmt.Fprintln(user.GetConn(), "Leaving "+room.name)
	room.broadcastInfo(INFO_LEAVE, user.GetName())
}

func (room *Chatroom) broadcastInfo(info, name string) {
	for _, user := range room.users {
		user.GetConn().Write([]byte("\n"))
		user.GetConn().Write([]byte(name + info + "\n"))
		user.GetConn().Write([]byte(service.GetPrefix(user.GetName())))
	}
}

func (room *Chatroom) ListUsers(user i.User) {
	fmt.Fprintln(user.GetConn())
	info := fmt.Sprintf("%d user(s) in the chat\n", len(room.users))
	user.GetConn().Write([]byte(info))
	for name := range room.users {
		user.GetConn().Write([]byte(name))
		user.GetConn().Write([]byte("\n"))
	}
	fmt.Fprintln(user.GetConn())
	user.GetConn().Write([]byte(service.GetPrefix(user.GetName())))
}

func (room *Chatroom) LogMessage(message string) {
	room.log += message + "\n"
}

func (room *Chatroom) DisplayLog(user i.User) {
	user.GetConn().Write([]byte(room.log))
}
