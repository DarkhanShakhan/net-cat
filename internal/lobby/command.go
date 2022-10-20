package lobby

import (
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

const (
	NON_EXIST_CMD  = "this command doesn't work here\n"
	ROOM_PREFIX    = "the chat with a given name "
	NON_EXIST_ROOM = ROOM_PREFIX + "doesn't exist\n"
	EXIST_ROOM     = ROOM_PREFIX + "exists\n"
	FULL_ROOM      = ROOM_PREFIX + "is full\n"
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

// FIXME: parsing commands and replying to invalid commands
func (lobby *Lobby) SendCommand(command string, user i.User) {
	temp := strings.Split(command, " ")

	switch len(temp) {
	case 1:
		lobby.cmdChannel <- NewCommand(command, "", user, Message{})
	case 2:
		lobby.cmdChannel <- NewCommand(temp[0], temp[1], user, Message{})
	case 3:
		lobby.cmdChannel <- NewCommand(temp[0], temp[1], user, NewMessage(strings.Join(temp[2:], " "), user))
	}
}

func (lobby *Lobby) ParseCommand(cmd Command) {
	switch cmd.command {
	case CMD_LIST:
		if _, ok := cmd.user.GetRoomName(); ok {
			cmd.user.GetConn().Write([]byte(NON_EXIST_CMD))
		} else {
			lobby.ListChats(cmd.user)
		}
	case CMD_USERS:
		if name, ok := cmd.user.GetRoomName(); ok {
			lobby.GetChatroom(name).ListUsers(cmd.user)
		} else {
			lobby.ListUsers(cmd.user)
		}
	case CMD_LEAVE:
		if name, ok := cmd.user.GetRoomName(); ok {
			lobby.GetChatroom(name).DeleteUser(cmd.user)
		} else {
			cmd.user.GetConn().Write([]byte(NON_EXIST_CMD))
		}
	case CMD_JOIN:
		room, ok := lobby.rooms[cmd.name]
		if !ok {
			cmd.user.GetConn().Write([]byte(NON_EXIST_ROOM))
		} else if room.IsFull() {
			cmd.user.GetConn().Write([]byte(FULL_ROOM))
		} else {
			lobby.rooms[cmd.name].AddUser(cmd.user)
		}
	case CMD_CREATE:
		if !lobby.CreateChatroom(cmd.name) {
			cmd.user.GetConn().Write([]byte(EXIST_ROOM))
		} else {
			lobby.GetChatroom(cmd.name).AddUser(cmd.user)
		}

	case CMD_HELP:
		lobby.ListCommands(cmd.user)
	}
}
