package internal

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
)

type User struct {
	name     string
	chatroom bool
	conn     net.Conn
	reader   *bufio.Reader
	scanner  *bufio.Scanner
	cmd      string
}

func NewUser(conn net.Conn) *User {
	return &User{
		chatroom: false,
		conn:     conn,
		reader:   bufio.NewReader(conn),
		scanner:  bufio.NewScanner(conn),
	}
}

func (u *User) sendChatname(g *gocui.Gui, v *gocui.View) error {
	name := v.Buffer()
	if !isEmpty(name) {
		fmt.Fprint(u.conn, u.cmd+" "+name)
		g.SetCurrentView("input")
	} else {
		g.SetCurrentView("options")
		g.CurrentView().SetCursor(0, 0)
		g.CurrentView().Highlight = true
	}
	v.Editable = false
	v.Clear()
	g.SetViewOnBottom("create")
	return nil
}

func (u *User) sendMsg(g *gocui.Gui, v *gocui.View) error {
	msg := v.Buffer()
	v.Clear()
	v.SetCursor(0, 0)
	output, err := g.View("output")
	if err != nil {
		return err
	}
	if !isEmpty(msg) && !isCommand(msg) {
		if u.name != "" && u.chatroom {
			fmt.Fprint(output, getPrefix(u.name)+msg)
		}
		fmt.Fprint(u.conn, msg)
	}
	return nil
}

var TIME_FORMAT = "2006-01-02 15:04:05"

func getPrefix(name string) string {
	return fmt.Sprintf("[%s][%s]:", time.Now().Format(TIME_FORMAT), name)
}

func isEmpty(message string) bool {
	res := strings.TrimSpace(message)
	return res == ""
}

func isCommand(message string) bool {
	return message[0] == '/'
}
