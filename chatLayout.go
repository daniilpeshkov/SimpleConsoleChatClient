package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

func ChatLayout(g *gocui.Gui) error {
	w, h := g.Size()
	v, err := g.SetView(ChatView, 2, 1, w-2, int(float32(h)*ChatHeightProcentage))

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
		tmp, _ = msg.GetField(TypeTime)
		time := string(tmp)
		v.Write([]byte(Red + time + " " + Blue + fmt.Sprintf("[%s]: ", name) + White + text + "\n"))
	default:
		break
	}
	return nil
}
