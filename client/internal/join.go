package internal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jroimartin/gocui"
)

func (u *User) getChats(g *gocui.Gui, v *gocui.View) error {
	//
	views := []string{}
	for _, view := range g.Views() {
		views = append(views, view.Name())
	}
	for _, view := range views {
		g.DeleteView(view)
	}
	u.conn.Write([]byte("/list\n"))
	title, err := u.reader.ReadString('\n')
	if err != nil {
		return err
	}
	temp := strings.Split(title, " ")
	nbr, err := strconv.Atoi(temp[0])
	if err != nil {
		return err
	}
	rooms := make([]string, nbr)
	for i := 0; i < nbr; i++ {
		room, err := u.reader.ReadString('\n')
		if err != nil {
			return err
		}
		rooms[i] = room
	}

	if err = u.displayChats(g, rooms); err != nil {
		return err
	}
	return nil
}

func (u *User) displayChats(g *gocui.Gui, rooms []string) error {
	maxX, maxY := g.Size()
	left, right := (maxX-20)/2, maxX-((maxX-20)/2)
	top := (maxY - len(rooms) - 1) / 2
	d, err := g.SetView("display", left, top, right, top+2)
	if err != gocui.ErrUnknownView {
		return err
	}
	d.Frame = false
	title := fmt.Sprintf("%d chat(s) available", len(rooms))
	fmt.Fprint(d, title)
	from := top + 2
	for i, room := range rooms {
		if v, err := g.SetView(room, left, from+(i*2), right, from+(i+1)*2); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			g.SetKeybinding(room, gocui.KeyEnter, gocui.ModNone, u.joinChat)
			fmt.Fprint(v, room)
		}
	}
	chMenu := &chatMenu{rooms: rooms}
	for _, room := range chMenu.rooms {
		g.SetKeybinding(room, gocui.KeyArrowUp, gocui.ModNone, chMenu.goUpChat)
		g.SetKeybinding(room, gocui.KeyArrowDown, gocui.ModNone, chMenu.goDownChat)
		g.SetKeybinding(room, gocui.KeyEnter, gocui.ModNone, u.StartChat)
	}
	g.SetCurrentView(rooms[0])
	g.CurrentView().BgColor = gocui.ColorGreen
	return nil
}

type chatMenu struct {
	rooms []string
}

func (c *chatMenu) goDownChat(g *gocui.Gui, v *gocui.View) error {
	from := v.Name()
	var to string
	for i, chat := range c.rooms {
		if from == chat {
			if i < len(c.rooms)-1 {
				to = c.rooms[i+1]
			} else {
				to = c.rooms[0]
			}
			break
		}
	}
	g.SetCurrentView(to)
	v.BgColor = gocui.ColorDefault
	g.CurrentView().BgColor = gocui.ColorGreen
	return nil
}

func (c *chatMenu) goUpChat(g *gocui.Gui, v *gocui.View) error {
	from := v.Name()
	var to string
	for i, room := range c.rooms {
		if from == room {
			if i != 0 {
				to = c.rooms[i-1]
			} else {
				to = c.rooms[len(c.rooms)-1]
			}
			break
		}
	}
	g.SetCurrentView(to)
	v.BgColor = gocui.ColorDefault
	g.CurrentView().BgColor = gocui.ColorGreen
	return nil
}

func (u *User) joinChat(g *gocui.Gui, v *gocui.View) error {
	return nil
}
