package internal

import (
	"net"

	"github.com/jroimartin/gocui"
)

type InputWidget struct {
	username string
	conn     net.Conn
}

func NewInputWidget() *InputWidget {
	return &InputWidget{}
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
	// g.SetCurrentView("input")
	return nil
}
