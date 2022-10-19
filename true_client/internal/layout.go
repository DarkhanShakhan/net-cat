package internal

import (
	"log"
	"net"

	"github.com/jroimartin/gocui"
)

func InitGui(conn net.Conn) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()
	g.Cursor = true
	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen

	input := NewInputWidget()
	output := NewOutputWidget()
	options := NewOptionsWidget()

	g.SetManager(input, output, options)
	g.SetCurrentView("input")
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, toggle); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func toggle(g *gocui.Gui, v *gocui.View) error {
	widgets := map[string]string{"input": "options", "options": "output", "output": "input"}
	if v != nil {
		curr := v.Name()
		next := widgets[curr]
		g.SetCurrentView(next)
	} else {
		g.SetCurrentView("input")
	}
	return nil
}
