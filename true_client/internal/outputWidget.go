package internal

import (
	"bufio"
	"fmt"
	"net"

	"github.com/jroimartin/gocui"
)

type OutputWidget struct {
	username string
	conn     net.Conn
}

func NewOutputWidget(conn net.Conn) *OutputWidget {
	return &OutputWidget{conn: conn}
}

func (out *OutputWidget) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	v, err := g.SetView("output", 2, 2, maxX-34, maxY-6)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	v.Autoscroll = true
	v.Wrap = true
	go out.Read(g)
	return nil
}

func (out *OutputWidget) Read(g *gocui.Gui) {
	scanner := bufio.NewScanner(out.conn)
	v, _ := g.View("output")
	for scanner.Scan() {
		msg := scanner.Text()
		fmt.Fprintln(v, msg)
		g.Update(func(*gocui.Gui) error { return nil })
	}
}
