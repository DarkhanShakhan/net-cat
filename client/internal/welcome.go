package internal

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func (u *User) InitGui() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()
	g.Cursor = true
	g.ASCII = true
	u.welcomeLayout(g)

	// quit key binding
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, u.quit); err != nil {
		log.Panicln(err)
	}
	// main loop
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (u *User) welcomeLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	inputY := 20
	if maxY < inputY+3 {
		inputY = 0
	}
	// logo display
	if v, err := g.SetView("logo", (maxX-23)/2, 1, maxX-(maxX-23)/2, 19); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		// v.Autoscroll = true
		logo := u.parseLogo()
		v.Write([]byte(logo))
	}
	// get username
	if v, err := g.SetView("username", (maxX-23)/2, ((maxY-inputY)/2+inputY)-1, maxX-(maxX-23)/2, ((maxY-inputY)/2+inputY)+1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		title := ""
		for u.reader.Buffered() > 0 {
			b, _ := u.reader.ReadByte()
			title += string(b)
		}
		v.Title = title[:len(title)-1]
		v.SetCursor(0, 0)
		v.Editable = true
	}
	// error display
	if v, err := g.SetView("error", (maxX-23)/2, ((maxY-inputY)/2+inputY)-3, maxX-(maxX-23)/2, ((maxY-inputY)/2+inputY)-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.FgColor = gocui.ColorRed
	}
	g.SetCurrentView("username")
	// enter your name keybinding
	if err := g.SetKeybinding("username", gocui.KeyEnter, gocui.ModNone, u.getName); err != nil {
		log.Panicln(err)
	}
	return nil
}

func (u *User) getName(g *gocui.Gui, v *gocui.View) error {
	u.name = v.Buffer()
	errLog, err := g.View("error")
	if err != nil {
		fmt.Println("here")
		log.Fatal(err)
	}
	errLog.Clear()
	if u.name == "" {
		// FIXME:have to check for existence
		fmt.Fprint(errLog, "empty name")
	} else {
		g.DeleteView("username")
		g.DeleteView("logo")
		g.Cursor = false
		u.conn.Write([]byte(u.name))
		u.menuLayout(g)
	}
	return nil
}

func (u *User) parseLogo() string {
	logo := ""
	for i := 0; i < 17; i++ {
		line, _ := u.reader.ReadString('\n')
		logo += line
	}
	return logo[:len(logo)-1]
}
