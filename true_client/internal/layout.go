package internal

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/jroimartin/gocui"
)

func InitGui(conn net.Conn) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()
	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen
	user := NewUser(conn)
	g.SetManagerFunc(user.layout)
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	user.keybindings(g)
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (u *User) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("input", 2, maxY-5, maxX-2, maxY-2); err != nil {
		if err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
		}
		v.Editable = true
		v.Wrap = true
		v.Title = "Enter your name"
		g.SetCurrentView("input")
	}
	if v, err := g.SetView("options", maxX-32, 2, maxX-2, maxY-6); err != nil {
		if err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
		}
		v.Autoscroll = true
		v.Wrap = true
		v.Title = "Options"
	}
	if v, err := g.SetView("output", 2, 2, maxX-34, maxY-6); err != nil {
		if err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
		}
		v.Autoscroll = true
		v.Wrap = true
		v.Title = "Output"
	}
	if v, err := g.SetView("create", maxX/2-15, maxY/2-1, maxX/2+15, maxY/2+1); err != nil {
		if err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
		}
		v.Frame = false
		g.SetViewOnBottom("create")

	}
	go u.Read(g)
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func toggle(g *gocui.Gui, v *gocui.View) error {
	widgets := map[string]string{"input": "options", "options": "output", "output": "input"}
	curr := v.Name()
	next := widgets[curr]
	c, _ := g.SetCurrentView(next)
	if next == "options" {
		c.Highlight = true
		c.SetCursor(0, 0)
	}
	if curr == "options" {
		v.Highlight = false
	}
	g.SetViewOnTop(next)
	return nil
}

func (u *User) Read(g *gocui.Gui) {
	scanner := bufio.NewScanner(u.conn)
	v, _ := g.View("output")
	for scanner.Scan() {
		msg := scanner.Text()
		if u.name == "" {
			msg = strings.TrimPrefix(msg, "Enter your name:")
			if msg[:7] == "Hello, " {
				u.name = msg[7 : len(msg)-1]
				v.Clear()
				input, _ := g.View("input")
				input.Title = "Input"
				u.menu(g)
			}
		}
		if msg == "use commands in the lobby, starting with '/'" {
			continue
		}
		if msg == "you can look all commands with '/help" {
			v.Clear()
			fmt.Fprintln(v, "go to options and choose one of the given")
			g.Update(func(*gocui.Gui) error { return nil })
			continue
		}
		if !strings.HasSuffix(msg, "]:") {
			fmt.Fprintln(v, msg)
			g.Update(func(*gocui.Gui) error { return nil })
		}
	}
}

func (u *User) menu(g *gocui.Gui) {
	v, _ := g.View("options")
	fmt.Fprintln(v, "Create a chat")
	fmt.Fprintln(v, "Join a chat")
	fmt.Fprintln(v, "Display users")
	fmt.Fprintln(v, "Quit")
}

func goDown(g *gocui.Gui, v *gocui.View) error {
	x, y := v.Cursor()
	if word, _ := v.Word(x, y); word != "" {
		if word, _ := v.Word(x, y+1); word != "" {
			v.SetCursor(x, y+1)
		}
	}
	return nil
}

func goUp(g *gocui.Gui, v *gocui.View) error {
	x, y := v.Cursor()
	if word, _ := v.Word(x, y); word != "" {
		if word, _ := v.Word(x, y-1); word != "" {
			v.SetCursor(x, y-1)
		}
	}
	return nil
}

func (u *User) command(g *gocui.Gui, v *gocui.View) error {
	word, _ := v.Word(v.Cursor())
	output, _ := g.View("output")
	output.Clear()
	switch word {
	case "Display":
		fmt.Fprintln(u.conn, "/users")
	case "Quit":
		return gocui.ErrQuit
	case "Create":
		v.Highlight = false
		// createChat(g)
		g.Update(createChat)
	}
	return nil
}

func createChat(g *gocui.Gui) error {
	v, _ := g.View("create")
	g.SetCurrentView("create")
	g.SetViewOnTop("create")
	v.Frame = true
	v.Editable = true
	v.Title = "Enter chat name"
	return nil
	// g.Update()
	// g.SetCurrentView("display")
	// g.SetViewOnTop("display")
}
