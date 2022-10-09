package lobby

import (
	"fmt"
	i "net-cat/internal/userInterface"
	"strings"
)

func (lobby *Lobby) ParseSignal() {
	for {
		select {
		case cmd := <-lobby.cmdChannel:
			lobby.ParseCommand(cmd)
		case msg := <-lobby.msgChannel:
			if _, ok := msg.user.GetRoomName(); !ok {
				fmt.Fprintln(msg.user.GetConn(), "use commands in the lobby, starting with '/'\nyou can look all commands with '/help")
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
		if signal != "" {
			lobby.msgChannel <- NewMessage(signal, user)
		}
	}
}
