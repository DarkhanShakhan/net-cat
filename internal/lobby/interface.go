package lobby

import (
	i "net-cat/internal/userInterface"
)

type Chatroom interface {
	DeleteUser(i.User)
	ListUsers(i.User)
	DisplayLog(i.User)
	IsFull() bool
	AddUser(i.User)
	GetChatName() string
	GetUsers() map[string]i.User
	LogMessage(string)
}
