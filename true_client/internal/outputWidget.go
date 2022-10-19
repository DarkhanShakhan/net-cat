package internal

import (
	"net"

	"github.com/jroimartin/gocui"
)

type OutputWidget struct {
	username string
	conn     net.Conn
}

func NewOutputWidget() *OutputWidget {
	return &OutputWidget{}
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
	return nil
}
