package internal

import (
	"bufio"
	"fmt"
	"log"
	"strings"
	"time"

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

func (u *User) sendMessage(g *gocui.Gui, v *gocui.View) error {
	msg := v.Buffer()
	if msg != "" {
		u.conn.Write([]byte(msg + "\n"))
		out, _ := g.View("output")
		fmt.Fprint(out, GetPrefix(u.name)+msg)
		v.Clear()
		v.SetCursor(0, 0)
	}
	return nil
}

func (u *User) read(g *gocui.Gui, v *gocui.View) {
	scanner := bufio.NewScanner(u.conn)
	for scanner.Scan() {
		msg := scanner.Text()
		if strings.HasSuffix(msg, fmt.Sprintf("[%s]:", u.name[:len(u.name)])) {
			continue
		}
		g.Update(func(g *gocui.Gui) error { return nil })
		fmt.Fprintln(v, msg)
	}
}

const (
	TIME_FORMAT = "2006-01-02 15:04:05"
)

func GetPrefix(name string) string {
	return fmt.Sprintf("[%s][%s]:", time.Now().Format(TIME_FORMAT), name)
}
