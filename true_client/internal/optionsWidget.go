package internal

import (
	"net"

	"github.com/jroimartin/gocui"
)

type OptionsWidget struct {
	username string
	conn     net.Conn
}

func NewOptionsWidget() *OptionsWidget {
	return &OptionsWidget{}
}

func (opt *OptionsWidget) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	v, err := g.SetView("options", maxX-32, 2, maxX-2, maxY-6)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	v.Autoscroll = true
	v.Wrap = true
	return nil
}
