package main

import (
	"bufio"
	"log"
	"net"

	"github.com/jroimartin/gocui"
)

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (u *User) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("logo", (maxX-23)/2, 1, maxX-(maxX-23)/2, 19); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.Autoscroll = true
		logo := u.ParseLogo()
		v.Write([]byte(logo))
	}
	if v, err := g.SetView("username", (maxX-23)/2, maxY-4, maxX-(maxX-23)/2, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		u.reader.ReadString(':')
		v.Title = "Enter your name"
		v.SetCursor(0, 0)
		v.Editable = true
	}
	g.SetCurrentView("username")
	return nil
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		log.Fatal(err)
	}
	u := &User{conn: conn, reader: bufio.NewReader(conn)}
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()
	g.Cursor = true
	g.Mouse = true
	u.layout(g)
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("username", gocui.KeyEnter, gocui.ModNone, u.getName); err != nil {
		log.Panicln(err)
	}
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (u *User) getName(g *gocui.Gui, v *gocui.View) error {
	u.name = v.Buffer()
	g.DeleteView("username")
	g.DeleteView("logo")
	g.Cursor = false
	u.conn.Write([]byte(u.name))
	u.menu(g)
	return nil
}

func (u *User) menu(g *gocui.Gui) {
	welcome, _ := g.SetView("welcome", 5, 2, 30, 5)
	welcome.BgColor = gocui.ColorCyan
	welcome.Frame = true
	welcome.Wrap = true
	welcome.Write([]byte("Welcome, " + u.name))

	create, _ := g.SetView("create", 5, 8, 20, 12)
	// create.SelBgColor = gocui.ColorWhite
	// create.Highlight = true
	create.Frame = true
	create.Write([]byte("Create a chat"))
	g.SetKeybinding("create", gocui.MouseLeft, gocui.ModNone, u.createchat)
}

func (u *User) createchat(g *gocui.Gui, v *gocui.View) error {
	chat, _ := g.SetView("createchat", 5, 9, 18, 15)
	g.SetViewOnTop("createchat")
	g.SetCurrentView("createchat")
	chat.Editable = true
	chat.SetCursor(0, 0)
	g.Cursor = true
	g.Mouse = false
	g.SetKeybinding("createchat", gocui.KeyEnter, gocui.ModNone, u.ChatCreated)
	return nil
}

func (u *User) ChatCreated(g *gocui.Gui, v *gocui.View) error {
	cr, _ := g.View("createchat")
	u.conn.Write([]byte("/create " + cr.Buffer()))
	g.SetManagerFunc(u.chatlayout)
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	return nil
}

func (u *User) chatlayout(g *gocui.Gui) error {
	g.Cursor = true
	in, _ := g.SetView("input", 4, 25, 45, 28)
	in.Editable = true
	out, _ := g.SetView("output", 4, 2, 45, 23)
	out.Autoscroll = true
	g.SetCurrentView("input")
	g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, u.SendMessage)
	return nil
}

func (u *User) SendMessage(g *gocui.Gui, v *gocui.View) error {
	msg := v.Buffer()
	u.conn.Write([]byte(msg))
	out, _ := g.View("output")
	out.Write([]byte(msg))
	v.Clear()
	v.SetCursor(0, 0)
	return nil
}

type User struct {
	conn   net.Conn
	reader *bufio.Reader
	name   string
}

func (u *User) ParseLogo() string {
	logo := ""
	for i := 0; i < 17; i++ {
		line, _ := u.reader.ReadString('\n')
		logo += line
	}
	return logo[:len(logo)-1]
}
