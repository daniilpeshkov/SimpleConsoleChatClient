package main

import (
	"log"

	simpleTcpMessage "github.com/daniilpeshkov/go-simple-tcp-message"
	"github.com/jroimartin/gocui"
)

const (
	ChatView  = "chat"
	InputView = "input"
)

var msgChan = make(chan *simpleTcpMessage.Message, 100)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}

	defer g.Close()
	g.FgColor = gocui.ColorCyan

	g.SetManager(gocui.ManagerFunc(ChatLayout), gocui.ManagerFunc(InputLayout), gocui.ManagerFunc(SetFocus(InputView)))

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(ChatView, gocui.MouseWheelUp, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	g.SetViewOnTop(InputView)
	g.Cursor = true
	g.MainLoop()
}

func SetFocus(name string) func(g *gocui.Gui) error {
	return func(g *gocui.Gui) error {
		_, err := g.SetCurrentView(name)
		return err
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
