package internal

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func (u *User) createChat(g *gocui.Gui, v *gocui.View) error {
	views := []string{}
	for _, view := range g.Views() {
		views = append(views, view.Name())
	}
	for _, view := range views {
		g.DeleteView(view)
	}
	maxX, maxY := g.Size()
	if in, err := g.SetView("chatname", maxX/2-7, maxY/2-1, maxX/2+7, maxY/2+1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		in.Editable = true
		in.Title = "Enter chat name"
	}
	g.SetCurrentView("chatname")
	if err := g.SetKeybinding("chatname", gocui.KeyEnter, gocui.ModNone, u.createdChat); err != nil {
		log.Panicln(err)
	}
	return nil
}

func (u *User) createdChat(g *gocui.Gui, v *gocui.View) error {
	views := []string{}
	for _, view := range g.Views() {
		views = append(views, view.Name())
	}
	for _, view := range views {
		g.DeleteView(view)
	}
	chatname := v.Buffer()[:len(v.Buffer())-1]
	query := fmt.Sprintf("/create %s\n", chatname)
	u.conn.Write([]byte(query))
	maxX, maxY := g.Size()
	if _, err := g.SetView("output", 2, 2, maxX-2, maxY-8); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	out, _ := g.View("output")
	go u.read(g, out)
	if tr, err := g.SetView("input", 2, maxY-7, maxX-2, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		tr.Editable = true
	}
	if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, u.sendMessage); err != nil {
		log.Panicln(err)
	}
	g.SetCurrentView("input")
	return nil
}
