package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func ChatLayout(g *gocui.Gui) error {
	w, h := g.Size()
	v, err := g.SetView(ChatView, w/2-w/4, 0, w/2+w/4, int(float32(h)*ChatHeightProcentage))

	if err != nil {
		v.Autoscroll = true
		v.Title = "Chat"
		v.Wrap = true
	}
	select {
	case msg := <-msgChan:
		tmp, _ := msg.GetField(TypeName)
		name := string(tmp)
		tmp, _ = msg.GetField(TypeText)
		text := string(tmp)

		v.Write([]byte(Blue + fmt.Sprintf("[%s]: ", name) + White + text + "\n"))
	default:
		break
	}
	return nil
}
