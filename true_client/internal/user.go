package internal

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/jroimartin/gocui"
)

type User struct {
	name     string
	chatroom bool
	conn     net.Conn
	reader   *bufio.Reader
	writer   *bufio.Writer
}

func NewUser(conn net.Conn) *User {
	return &User{
		chatroom: false,
		conn:     conn,
		reader:   bufio.NewReader(conn),
		writer:   bufio.NewWriter(conn),
	}
}

func (u *User) keybindings(g *gocui.Gui) {
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, toggle); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, u.sendMsg); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("options", gocui.KeyArrowDown, gocui.ModNone, goDown); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("options", gocui.KeyArrowUp, gocui.ModNone, goUp); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("options", gocui.KeyEnter, gocui.ModNone, u.command); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("create", gocui.KeyEnter, gocui.ModNone, u.createChat); err != nil {
		log.Panicln(err)
	}
}

func (u *User) createChat(g *gocui.Gui, v *gocui.View) error {
	name := v.Buffer()
	fmt.Fprintln(u.conn, "/create "+name)
	v.Editable = false
	v.Clear()
	g.SetCurrentView("input")
	g.SetViewOnBottom("create")
	u.chatroom = true
	return nil
}

func (u *User) sendMsg(g *gocui.Gui, v *gocui.View) error {
	u.clearBuffer()
	msg := v.Buffer()
	v.Clear()
	v.SetCursor(0, 0)
	output, err := g.View("output")
	if err != nil {
		return err
	}
	if !isEmpty(msg) && !isCommand(msg) {
		if u.name != "" && u.chatroom {
			fmt.Fprint(output, "["+u.name+"]"+msg)
		}
		fmt.Fprint(u.conn, msg)
	}
	return nil
}

func (u *User) clearBuffer() {
	for i := 0; i < u.reader.Buffered(); i++ {
		u.reader.ReadByte()
	}
}

func isEmpty(message string) bool {
	res := strings.TrimSpace(message)
	return res == ""
}

func isCommand(message string) bool {
	return message[0] == '/'
}
