package user

import (
	"fmt"
	"net"
)

type User struct {
	name string
	room string
	conn net.Conn
}

func NewUser(name string, conn net.Conn) *User {
	greet := fmt.Sprintf("Hello, %s!\n", name)
	conn.Write([]byte(greet))
	return &User{name: name, conn: conn}
}

func (user *User) SetRoomName(chatName string) {
	user.room = chatName
}

func (user *User) GetName() string {
	return user.name
}
func (user *User) GetConn() net.Conn {
	return user.conn
}

func (user *User) GetRoomName() (string, bool) {
	if user.room != "" {
		return user.room, true
	}
	return "", false
}
