package lobby

import (
	"net-cat/internal/service"
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
)

const (
	NON_EXIST_CMD  = "this command doesn't work here\n"
	ROOM_PREFIX    = "the chat with a given name "
	NON_EXIST_ROOM = ROOM_PREFIX + "doesn't exist\n"
	EXIST_ROOM     = ROOM_PREFIX + "exists\n"
	FULL_ROOM      = ROOM_PREFIX + "is full\n"
	INVALID_CMD    = "invalid command\n"
	INVALID_ARG    = "too many arguments\n"
	INVALID_NAME   = "invalid chat name\n"
)

type Command struct {
	command string
	name    string
	user    i.User
}

func NewCommand(command string, name string, user i.User) Command {
	return Command{command: command, name: name, user: user}
}

func (lobby *Lobby) SendCommand(command string, user i.User) {
	temp := strings.Split(command, " ")
	switch temp[0] {
	case CMD_HELP, CMD_LIST, CMD_USERS, CMD_LEAVE:
		if len(temp) > 1 {
			user.GetConn().Write([]byte(INVALID_ARG))
			if _, ok := user.GetRoomName(); ok {
				user.GetConn().Write([]byte(service.GetPrefix(user.GetName())))
			}
		} else {
			lobby.cmdChannel <- NewCommand(command, "", user)
		}
	case CMD_JOIN, CMD_CREATE:
		name := strings.Join(temp[1:], " ")
		if service.ValidInput(name) {
			lobby.cmdChannel <- NewCommand(temp[0], name, user)
		} else {
			user.GetConn().Write([]byte(INVALID_NAME))
		}

	default:
		user.GetConn().Write([]byte(INVALID_CMD))
		if _, ok := user.GetRoomName(); ok {
			user.GetConn().Write([]byte(service.GetPrefix(user.GetName())))
		}
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
			if lobby.GetChatroom(name).IsEmpty() {
				lobby.DeleteChatroom(name)
			}
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
