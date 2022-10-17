package internal

import (
	"bufio"
	"net"

	"github.com/jroimartin/gocui"
)

type User struct {
	name   string
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewUser(conn net.Conn) *User {
	return &User{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}
}

func (u *User) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
