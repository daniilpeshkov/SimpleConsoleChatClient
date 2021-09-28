package main

import (
	"github.com/jroimartin/gocui"
)

var printChan = make(chan string, 10)

func ChatLayout(g *gocui.Gui) error {
	w, h := g.Size()
	v, err := g.SetView(ChatView, 2, 1, w-2, int(float32(h)*ChatHeightProcentage))

	if err != nil {
		v.Autoscroll = true
		v.Title = "Chat"
		v.Wrap = true
	}
	select {
	case msg := <-printChan:

		v.Write([]byte(msg + "\n"))
	default:
		break
	}
	return nil
}
