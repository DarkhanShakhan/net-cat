package internal

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jroimartin/gocui"
)

var options []string = []string{"create", "join", "users", "exit"}

func (u *User) menuLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	left, right := (maxX-15)/2, maxX-((maxX-15)/2)
	top, _ := (maxY-18)/2, maxY-((maxY-18)/2)
	g.Cursor = false
	// menu title
	if v, err := g.SetView("title", left, top, right, top+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		fmt.Fprint(v, "Menu")
	}
	if v, err := g.SetView("create", left, top+3, right, top+5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprint(v, "Create a chat")
	}
	if v, err := g.SetView("join", left, top+5, right, top+7); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprint(v, "Join a chat")
	}
	if v, err := g.SetView("users", left, top+7, right, top+9); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprint(v, "Users")
	}
	if v, err := g.SetView("exit", left, top+9, right, top+11); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprint(v, "Exit")
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, u.quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("join", gocui.KeyEnter, gocui.ModNone, u.getChats); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("exit", gocui.KeyEnter, gocui.ModNone, u.quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("users", gocui.KeyEnter, gocui.ModNone, u.getUsers); err != nil {
		log.Panicln(err)
	}
	g.SetCurrentView("create")
	g.CurrentView().BgColor = gocui.ColorGreen
	for _, option := range options {
		g.SetKeybinding(option, gocui.KeyArrowUp, gocui.ModNone, u.goUp)
		g.SetKeybinding(option, gocui.KeyArrowDown, gocui.ModNone, u.goDown)
	}
	return nil
}

func (u *User) goDown(g *gocui.Gui, v *gocui.View) error {
	from := v.Name()
	var to string
	for i, option := range options {
		if from == option {
			if i < len(options)-1 {
				to = options[i+1]
			} else {
				to = options[0]
			}
			break
		}
	}
	g.SetCurrentView(to)
	v.BgColor = gocui.ColorDefault
	g.CurrentView().BgColor = gocui.ColorGreen
	return nil
}

func (u *User) goUp(g *gocui.Gui, v *gocui.View) error {
	from := v.Name()
	var to string
	for i, option := range options {
		if from == option {
			if i != 0 {
				to = options[i-1]
			} else {
				to = options[len(options)-1]
			}
			break
		}
	}
	g.SetCurrentView(to)
	v.BgColor = gocui.ColorDefault
	g.CurrentView().BgColor = gocui.ColorGreen
	return nil
}

func (u *User) getUsers(g *gocui.Gui, v *gocui.View) error {
	v.BgColor = gocui.ColorDefault
	u.conn.Write([]byte("/users\n"))
	title, err := u.reader.ReadString('\n')
	if err != nil {
		return err
	}
	temp := strings.Split(title, " ")
	nbr, err := strconv.Atoi(temp[0])
	if err != nil {
		return err
	}

	maxX, maxY := g.Size()
	left, right := (maxX-20)/2, maxX-((maxX-20)/2)
	top := (maxY - nbr - 1) / 2
	d, err := g.SetView("display", left, top, right, top+nbr*2+1)
	if err != gocui.ErrUnknownView {
		return err
	}
	fmt.Fprint(d, title)
	d.Frame = true
	for i := 0; i < nbr; i++ {
		name, err := u.reader.ReadString('\n')
		if err != nil {
			return err
		}
		fmt.Fprint(d, "-- "+name)
	}
	g.SetViewOnTop("display")
	g.SetCurrentView("display")
	g.SetKeybinding("display", gocui.KeyEnter, gocui.ModNone, u.closeDisplay)
	return nil
}

// func (u *User) getChats(g *gocui.Gui, v *gocui.View) error {
// 	v.BgColor = gocui.ColorDefault
// 	u.conn.Write([]byte("/list\n"))
// 	title, err := u.reader.ReadString('\n')
// 	if err != nil {
// 		return err
// 	}
// 	temp := strings.Split(title, " ")
// 	nbr, err := strconv.Atoi(temp[0])
// 	if err != nil {
// 		return err
// 	}
// 	maxX, maxY := g.Size()
// 	left, right := (maxX-20)/2, maxX-((maxX-20)/2)
// 	top := (maxY - nbr - 1) / 2
// 	d, err := g.SetView("display", left, top, right, top+2)
// 	if err != gocui.ErrUnknownView {
// 		return err
// 	}
// 	fmt.Fprint(d, title)
// 	d.Frame = true
// 	rooms := make([]string, nbr)
// 	for i := 0; i < nbr; i++ {
// 		room, err := u.reader.ReadString('\n')
// 		if err != nil {
// 			return err
// 		}
// 		rooms[i] = room
// 		// fmt.Fprint(d, "-- "+room)
// 	}
// 	// fmt.Fprint(d, len(rooms))
// 	from := top + 2
// 	for i, room := range rooms {
// 		if v, err := g.SetView(room, left, from+(i*2), right, from+(i+1)*2); err != nil {
// 			if err != gocui.ErrUnknownView {
// 				return err
// 			}
// 			g.SetKeybinding(room, gocui.KeyEnter, gocui.ModNone, u.closeDisplay)
// 			fmt.Fprint(v, room)
// 		}
// 	}
// 	if _, err := g.View(rooms[0]); err == nil {
// 		g.SetCurrentView(rooms[0])
// 		g.CurrentView().BgColor = gocui.ColorGreen
// 	} else {
// 		g.SetCurrentView("display")
// 	}
// 	// g.SetViewOnTop("display")
// 	// g.SetCurrentView("display")
// 	g.SetKeybinding("display", gocui.KeyEnter, gocui.ModNone, u.closeDisplay)
// 	return nil
// }

func (u *User) closeDisplay(g *gocui.Gui, v *gocui.View) error {
	g.DeleteView(v.Name())
	g.SetCurrentView("create")
	g.CurrentView().BgColor = gocui.ColorGreen
	return nil
}