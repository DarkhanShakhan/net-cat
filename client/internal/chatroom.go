package internal

import (
	"bufio"
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func (u *User) StartChat(g *gocui.Gui, v *gocui.View) error {
	views := []string{}
	for _, view := range g.Views() {
		views = append(views, view.Name())
	}
	for _, view := range views {
		g.DeleteView(view)
	}
	query := fmt.Sprintf("/join %s\n", v.Name())
	u.conn.Write([]byte(query))
	maxX, maxY := g.Size()
	if _, err := g.SetView("output", 2, 2, maxX-2, maxY-8); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		// go u.read(v)
	}
	out, _ := g.View("output")
	go u.read(out)
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

func (u *User) sendMessage(g *gocui.Gui, v *gocui.View) error {
	msg := v.Buffer()
	if msg != "" {
		u.conn.Write([]byte(msg + "\n"))
		// u.reader.ReadString(':')
		out, _ := g.View("output")
		fmt.Fprint(out, msg)
		v.Clear()
		v.SetCursor(0, 0)
	}
	return nil
}

func (u *User) read(v *gocui.View) {
	scanner := bufio.NewScanner(u.conn)
	for scanner.Scan() {
		msg := scanner.Text()
		if msg[len(msg)-2:] == "]:" {
			continue
		}
		fmt.Fprintln(v, msg)
	}
	// for {
	// 	msg, _ := u.reader.ReadString('\n')
	// 	fmt.Fprint(v, msg)
	// }
}
