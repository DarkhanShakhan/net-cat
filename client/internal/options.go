package internal

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
)

func roomOptions(g *gocui.Gui) error {
	out, _ := g.View("output")
	out.SetOrigin(0, 0)
	v, _ := g.View("options")
	v.Clear()
	fmt.Fprintln(v, "Users")
	fmt.Fprintln(v, "Leave the chat")
	return nil
}

func lobbyOptions(g *gocui.Gui) error {
	out, _ := g.View("output")
	out.SetOrigin(0, 0)
	v, _ := g.View("options")
	v.Clear()
	fmt.Fprintln(v, "Create a chat")
	fmt.Fprintln(v, "Join a chat")
	fmt.Fprintln(v, "Users")
	fmt.Fprintln(v, "Chats")
	fmt.Fprintln(v, "Quit")
	return nil
}

func goDown(g *gocui.Gui, v *gocui.View) error {
	x, y := v.Cursor()
	if word, _ := v.Word(x, y); word != "" {
		if word, _ := v.Word(x, y+1); word != "" {
			v.SetCursor(x, y+1)
		}
	}
	return nil
}

func goUp(g *gocui.Gui, v *gocui.View) error {
	x, y := v.Cursor()
	if word, _ := v.Word(x, y); word != "" {
		if word, _ := v.Word(x, y-1); word != "" {
			v.SetCursor(x, y-1)
		}
	}
	return nil
}

func (u *User) command(g *gocui.Gui, v *gocui.View) error {
	word, _ := v.Word(v.Cursor())
	output, _ := g.View("output")
	switch word {
	case "Users":
		fmt.Fprintln(u.conn, "/users")
	case "Chats":
		fmt.Fprintln(u.conn, "/list")
	case "Quit":
		return gocui.ErrQuit
	case "Create", "Join":
		output.Clear()
		v.Highlight = false
		u.cmd = "/" + strings.ToLower(word)
		g.Update(enterChat)
	case "Leave":
		output.Clear()
		fmt.Fprintln(u.conn, "/leave")
	}
	return nil
}

func enterChat(g *gocui.Gui) error {
	v, _ := g.View("create")
	v.SetCursor(0, 0)
	g.SetCurrentView("create")
	g.SetViewOnTop("create")
	v.Frame = true
	v.Editable = true
	v.Title = "Enter chat name"
	return nil
}
