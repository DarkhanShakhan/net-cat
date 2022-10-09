package lobby

import (
	"fmt"
	"net-cat/internal/service"
	i "net-cat/internal/userInterface"
)

func (lobby *Lobby) ListChats(user i.User) {
	for key := range lobby.rooms {
		fmt.Fprintln(user.GetConn(), key)
	}
}

func (lobby *Lobby) ListUsers(user i.User) {
	for name, user := range lobby.users {
		fmt.Fprint(user.GetConn(), name)
		if name, ok := user.GetRoomName(); ok {
			fmt.Fprint(user.GetConn(), " -- "+lobby.GetChatroom(name).GetChatName())
		}
		fmt.Fprintln(user.GetConn(), "")
	}
}

func (lobby *Lobby) ListCommands(user i.User) {
	if _, ok := user.GetRoomName(); !ok {
		fmt.Fprintln(user.GetConn(), CMD_LIST+" -> display all chatrooms.")
		fmt.Fprintln(user.GetConn(), CMD_CREATE+" roomName -> create a new chatroom with a given room name.")
		fmt.Fprintln(user.GetConn(), CMD_JOIN+" roomName -> join the given chatroom.")
		fmt.Fprintln(user.GetConn(), CMD_USERS+" -> list all users.")
	} else {
		fmt.Fprintln(user.GetConn(), CMD_USERS+" -> list all users in the chat.")
		fmt.Fprintln(user.GetConn(), CMD_LEAVE+" -> leave the chatroom to the lobby.")
		fmt.Fprint(user.GetConn(), service.GetPrefix(user.GetName()))
	}
}
