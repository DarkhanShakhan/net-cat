package chatroom

import (
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
	room.DisplayLog(user)

	user.GetConn().Write([]byte(service.GetPrefix(user.GetName())))
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
		user.GetConn().Write([]byte("\n"))
		user.GetConn().Write([]byte(name + info + "\n"))
		user.GetConn().Write([]byte(service.GetPrefix(user.GetName())))
	}
}

func (room *Chatroom) ListUsers(user i.User) {
	user.GetConn().Write([]byte("\n"))
	for name := range room.users {
		user.GetConn().Write([]byte(name))
		user.GetConn().Write([]byte("\n"))
	}
	user.GetConn().Write([]byte(service.GetPrefix(user.GetName())))
}

func (room *Chatroom) LogMessage(message string) {
	room.log += message + "\n"
}

func (room *Chatroom) DisplayLog(user i.User) {
	user.GetConn().Write([]byte(room.log))
}
