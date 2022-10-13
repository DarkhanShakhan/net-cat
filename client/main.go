package main

import (
	"log"
	"net"
	"net-cat/internal/service"

	"github.com/jroimartin/gocui"
)

var logo = ""

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("logo", (maxX-23)/2, 1, maxX-(maxX-23)/2, 19); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.Autoscroll = true
		if logo == "" {
			logo = service.ParseLogo()
		}
		v.Write([]byte(logo))
	}
	if v, err := g.SetView("username", (maxX-23)/2, maxY-4, maxX-(maxX-23)/2, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Enter your name"
		g.Cursor = true
		v.Editable = true
	}
	g.SetCurrentView("username")
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func main() {
	tcpaddr := &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8080}
	tcplistener, _ := net.ListenTCP("tcp", tcpaddr)
	for {
		tcpconn, _ := tcplistener.AcceptTCP()

		g, _ := gocui.NewGui(gocui.OutputNormal)
		defer g.Close()
		v, _ := g.SetView("de", 1, 2, 2, 3)
		v.Title = "Hello"
		tcpconn.ReadFrom(v)
		if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
			log.Panicln(err)
		}
		if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
			log.Panicln(err)
		}
		tcpconn.Close()
	}
}

func run() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.Cursor = true
	g.Highlight = true

	layout(g)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("username", gocui.KeyEnter, gocui.ModNone, startChat); err != nil {
		log.Panicln(err)
	}
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func startChat(g *gocui.Gui, v *gocui.View) error {
	name := v.Buffer()
	g.DeleteView("logo")
	g.DeleteView("username")
	chat(g, name)

	return nil
}

func chat(g *gocui.Gui, name string) error {
	maxX, maxY := g.Size()
	i, _ := g.SetView("input", 6, maxY-4, maxX-6, maxY-1)
	i.Editable = true
	o, _ := g.SetView("output", 6, 1, maxX-6, maxY-5)
	o.Autoscroll = true
	g.SetCurrentView("input")
	u := &User{name: name[:len(name)-1]}
	g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, u.printInput)
	return nil
}

type User struct {
	name string
}

func (u *User) printInput(g *gocui.Gui, v *gocui.View) error {
	msg := v.Buffer()
	v.Clear()
	v.SetCursor(0, 0)
	o, _ := g.View("output")
	prefix := service.GetPrefix(u.name)
	msg = prefix + msg
	o.Write([]byte(msg))
	return nil
}
