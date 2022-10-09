package userInterface

import "net"

type User interface {
	SetRoomName(string)
	GetName() string
	GetConn() net.Conn
	GetRoomName() (string, bool)
}
