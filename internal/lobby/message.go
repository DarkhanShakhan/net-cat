package lobby

import (
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
			otherUser.GetConn().Write([]byte("\n"))
			otherUser.GetConn().Write([]byte(msg.prefix + msg.text + "\n"))
			otherUser.GetConn().Write([]byte(service.GetPrefix(otherUser.GetName())))
		} else {
			otherUser.GetConn().Write([]byte(msg.prefix))
		}
	}
	chat.LogMessage(msg.prefix + msg.text)
}
