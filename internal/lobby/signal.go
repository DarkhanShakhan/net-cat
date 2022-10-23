package lobby

import (
	"net-cat/internal/service"
	i "net-cat/internal/userInterface"
	"strings"
)

const INVALID_LOBBY_INPUT = "use commands in the lobby, starting with '/'\nyou can look all commands with '/help"

func (lobby *Lobby) ParseSignal() {
	for {
		select {
		case cmd := <-lobby.cmdChannel:
			lobby.ParseCommand(cmd)
		case msg := <-lobby.msgChannel:
			if _, ok := msg.user.GetRoomName(); !ok {
				msg.user.GetConn().Write([]byte(INVALID_LOBBY_INPUT + "\n"))
			} else {
				lobby.BroadcastMsg(msg)
			}
		}
	}
}

func (lobby *Lobby) SendSignal(signal string, user i.User) {
	switch {
	case strings.HasPrefix(signal, CMD):
		lobby.SendCommand(signal, user)
	default:
		if service.ValidInput(signal) {
			lobby.msgChannel <- NewMessage(signal, user)
		} else {
			user.GetConn().Write([]byte(service.GetPrefix(user.GetName())))
		}
	}
}
