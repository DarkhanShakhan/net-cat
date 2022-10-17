package internal

import (
	"fmt"

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
	left, right := (maxX-15)/2, maxX-((maxX-15)/2)
	top, _ := (maxY-18)/2, maxY-((maxY-18)/2)
	if v, err := g.SetView("input", left-15, top, right+15, top+7); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		for i := 0; i < 2; i++ {
			msg, _ := u.reader.ReadString('\n')
			fmt.Fprint(v, msg)
		}

	}
	return nil
}
