package internal

import (
	"fmt"
	"net"

	"github.com/jroimartin/gocui"
)

type InputWidget struct {
	username string
	conn     net.Conn
}

func NewInputWidget(conn net.Conn) *InputWidget {
	return &InputWidget{conn: conn}
}

func (in *InputWidget) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	v, err := g.SetView("input", 2, maxY-5, maxX-2, maxY-2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	v.Editable = true
	v.Wrap = true
	g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, in.sendMsg)
	return nil
}

func (in *InputWidget) sendMsg(g *gocui.Gui, v *gocui.View) error {
	msg := v.Buffer()
	v.Clear()
	v.SetCursor(0, 0)
	fmt.Fprintln(in.conn, msg)
	return nil
}
