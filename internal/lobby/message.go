package lobby

import (
	"fmt"
	"net-cat/internal/service"
	i "net-cat/internal/userInterface"
)

type Message struct {
	text   string
	user   i.User
	prefix string
}

func NewMessage(text string, user i.User) Message {
	return Message{text: text, prefix: service.GetPrefix(user.GetName()), user: user}
}

func (lobby *Lobby) BroadcastMsg(msg Message) {
	name, _ := msg.user.GetRoomName()
	chat := lobby.GetChatroom(name)
	for key, otherUser := range chat.GetUsers() {
		if key != msg.user.GetName() {
			fmt.Fprintln(otherUser.GetConn(), "")
			fmt.Fprintln(otherUser.GetConn(), msg.prefix+msg.text)
			fmt.Fprint(otherUser.GetConn(), service.GetPrefix(otherUser.GetName()))
		} else {
			fmt.Fprint(otherUser.GetConn(), msg.prefix)
		}
	}
	chat.LogMessage(msg.prefix + msg.text)
}
