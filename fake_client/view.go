package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net-cat/internal/service"
	"strconv"
	"strings"

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
	u.name = v.Buffer()[:len(v.Buffer())-1]
	g.DeleteView("username")
	g.DeleteView("logo")
	g.Cursor = false
	u.conn.Write([]byte(u.name + "\n"))
	u.menu(g)
	return nil
}

func (u *User) menu(g *gocui.Gui) {
	welcome, _ := g.SetView("welcome", 5, 2, 30, 5)
	welcome.BgColor = gocui.ColorCyan
	welcome.Frame = true
	welcome.Wrap = true
	welcome.Write([]byte("Welcome, " + u.name))
	// g.CurrentView().BgColor = gocui.ColorGreen
	create, _ := g.SetView("create", 5, 8, 20, 11)
	// create.SelBgColor = gocui.ColorWhite
	create.FgColor = gocui.ColorGreen
	create.Frame = true
	g.SetCurrentView("create")
	g.CurrentView().BgColor = gocui.ColorGreen
	create.Write([]byte("Create a chat"))
	g.SetKeybinding("create", gocui.KeyEnter, gocui.ModNone, u.createchat)

	join, _ := g.SetView("join", 5, 12, 20, 15)
	join.Frame = true
	join.Write([]byte("Join a chat"))
	join.FgColor = gocui.ColorGreen
	g.SetKeybinding("join", gocui.KeyEnter, gocui.ModNone, u.showChats)

	users, _ := g.SetView("users", 5, 16, 20, 19)
	users.Frame = true
	users.Write([]byte("Show all users"))
	users.FgColor = gocui.ColorGreen
	g.SetKeybinding("users", gocui.KeyEnter, gocui.ModNone, u.showUsers)

	exit, _ := g.SetView("exit", 5, 20, 20, 23)
	exit.Frame = true
	exit.Write([]byte("Exit"))
	exit.FgColor = gocui.ColorGreen
	g.SetKeybinding("exit", gocui.KeyEnter, gocui.ModNone, quit)

	for _, option := range options {
		g.SetKeybinding(option, gocui.KeyArrowUp, gocui.ModNone, u.goUp)
		g.SetKeybinding(option, gocui.KeyArrowDown, gocui.ModNone, u.goDown)
	}
}

func (u *User) showChats(g *gocui.Gui, v *gocui.View) error {
	v.BgColor = gocui.ColorDefault
	u.conn.Write([]byte("/list\n"))
	sh, _ := g.SetView("showChats", 30, 15, 75, 20)
	g.SetViewOnTop("showChats")
	g.SetCurrentView("showChats")
	sh.Frame = true
	g.SetKeybinding("showChats", gocui.KeyEnter, gocui.ModNone, u.CloseShowChats)

	fmt.Fprintln(sh, "List of chats here:")
	d, _ := u.reader.ReadString('\n')
	t := strings.Split(d, " ")
	nbr, _ := strconv.Atoi(t[0])
	for i := 0; i < nbr; i++ {
		r, _ := u.reader.ReadString('\n')
		fmt.Fprint(sh, r)
	}

	return nil
}
func (u *User) CloseShowChats(g *gocui.Gui, v *gocui.View) error {
	g.DeleteView("showChats")
	g.SetCurrentView("join")
	g.CurrentView().BgColor = gocui.ColorGreen
	return nil
}

func (u *User) showUsers(g *gocui.Gui, v *gocui.View) error {
	v.BgColor = gocui.ColorDefault
	u.conn.Write([]byte("/users\n"))
	sh, _ := g.SetView("showUsers", 30, 15, 75, 20)
	g.SetViewOnTop("showUsers")
	g.SetCurrentView("showUsers")
	sh.Frame = true
	sh.Write([]byte("List of users here: "))
	g.SetKeybinding("showUsers", gocui.KeyEnter, gocui.ModNone, u.CloseShowUsers)

	d, _ := u.reader.ReadString('\n')
	t := strings.Split(d, " ")
	nbr, _ := strconv.Atoi(t[0])
	for i := 0; i < nbr; i++ {
		r, _ := u.reader.ReadString('\n')
		fmt.Fprint(sh, r)
	}
	return nil
}

func (u *User) CloseShowUsers(g *gocui.Gui, v *gocui.View) error {
	g.DeleteView("showUsers")
	g.SetCurrentView("users")
	g.CurrentView().BgColor = gocui.ColorGreen
	return nil
}

var options []string = []string{"create", "join", "users", "exit"}

func (u *User) goDown(g *gocui.Gui, v *gocui.View) error {
	from := v.Name()
	var to string
	for i, option := range options {
		if from == option {
			if i < len(options)-1 {
				to = options[i+1]
			} else {
				to = options[0]
			}
			break
		}
	}
	g.SetCurrentView(to)
	v.BgColor = gocui.ColorDefault
	g.CurrentView().BgColor = gocui.ColorGreen
	return nil
}

func (u *User) goUp(g *gocui.Gui, v *gocui.View) error {
	from := v.Name()
	var to string
	for i, option := range options {
		if from == option {
			if i != 0 {
				to = options[i-1]
			} else {
				to = options[len(options)-1]
			}
			break
		}
	}
	g.SetCurrentView(to)
	v.BgColor = gocui.ColorDefault
	g.CurrentView().BgColor = gocui.ColorGreen
	return nil
}

func (u *User) createchat(g *gocui.Gui, v *gocui.View) error {
	chat, _ := g.SetView("createchat", 5, 9, 18, 15)
	g.SetViewOnTop("createchat")
	g.SetCurrentView("createchat")
	chat.Editable = true
	chat.SetCursor(0, 0)
	g.Cursor = true
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
	if msg != "" {
		prefix := service.GetPrefix(u.name)
		u.conn.Write([]byte(prefix + msg))
		out, _ := g.View("output")
		out.Write([]byte(msg))
		v.Clear()
		v.SetCursor(0, 0)
	}
	// prefix := service.GetPrefix(u.name)
	// u.conn.Write([]byte(prefix + msg))
	// out, _ := g.View("output")
	// out.Write([]byte(msg))
	// v.Clear()
	// v.SetCursor(0, 0)
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
