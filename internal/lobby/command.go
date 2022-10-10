package lobby

import (
	"fmt"
	i "net-cat/internal/userInterface"
	"strings"
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

type Command struct {
	command string
	name    string
	message Message
	user    i.User
}

func NewCommand(command string, name string, user i.User, message Message) Command {
	return Command{command: command, name: name, user: user, message: message}
}

func (lobby *Lobby) SendCommand(command string, user i.User) {
	temp := strings.Split(command, " ")

	switch len(temp) {
	case 1:
		lobby.cmdChannel <- NewCommand(command, "", user, Message{})
	case 2:
		lobby.cmdChannel <- NewCommand(temp[0], temp[1], user, Message{})
	case 3:
		lobby.cmdChannel <- NewCommand(temp[0], temp[1], user, NewMessage(temp[2], user))
	}
}

func (lobby *Lobby) ParseCommand(cmd Command) {
	switch cmd.command {
	case CMD_LIST:
		lobby.ListChats(cmd.user)
	case CMD_USERS:
		if name, ok := cmd.user.GetRoomName(); ok {
			lobby.GetChatroom(name).ListUsers(cmd.user)
		} else {
			lobby.ListUsers(cmd.user)
		}
	case CMD_DIRECT:
		to := lobby.users[cmd.name]
		toChatName, toOk := to.GetRoomName()
		fromChatName, fromOk := cmd.user.GetRoomName()
		if fromOk && toOk && toChatName == fromChatName {
			fmt.Fprintln(to.GetConn())
			fmt.Fprintln(to.GetConn(), "[DIRECT]"+cmd.message.prefix+cmd.message.text)
			fmt.Fprint(cmd.user.GetConn(), cmd.message.prefix)
			fmt.Fprint(to.GetConn(), cmd.message.prefix)
		}
	case CMD_LEAVE:
		name, _ := cmd.user.GetRoomName()
		lobby.GetChatroom(name).DeleteUser(cmd.user)
	case CMD_JOIN:
		if lobby.rooms[cmd.name].IsFull() {
			fmt.Fprintln(cmd.user.GetConn(), "The Chat is full, join later or create a new one")
		} else {
			lobby.rooms[cmd.name].AddUser(cmd.user)
		}
	case CMD_CREATE:
		if !lobby.CreateChatroom(cmd.name) {
			fmt.Fprintln(cmd.user.GetConn(), "The chat with a given name exists")
		} else {
			lobby.GetChatroom(cmd.name).AddUser(cmd.user)
		}

	case CMD_HELP:
		lobby.ListCommands(cmd.user)
	}
}