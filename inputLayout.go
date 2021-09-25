package main

import (
	"strings"

	simpleTcpMessage "github.com/daniilpeshkov/go-simple-tcp-message"
	"github.com/jroimartin/gocui"
)

func InputLayout(g *gocui.Gui) error {
	w, h := g.Size()
	v, err := g.SetView(InputView, w/2-w/4, h-int(float32(h)*InputHeightProcentage), w/2+w/4, h-1)

	if err != nil {
		v.Editor = gocui.EditorFunc(InputEditor)
		v.Autoscroll = false
		v.Title = "Input"
		v.Editable = true
		v.Frame = true
		v.Overwrite = true
		v.Wrap = true
		v.FgColor = gocui.ColorWhite
	}

	return nil
}

func InputEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
	case key == gocui.KeyEnter:
		text := strings.TrimSpace(v.Buffer())
		if len(text) != 0 {
			msg := simpleTcpMessage.NewMessage()
			msg.AppendField(TypeName, []byte("User1"))
			msg.AppendField(TypeText, []byte(text))

			msgChan <- msg

			v.Clear()
			v.SetCursor(0, 0)
		}
	case key == gocui.KeyArrowDown:
		v.MoveCursor(0, 1, false)
	case key == gocui.KeyArrowUp:
		v.MoveCursor(0, -1, false)
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case key == gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
	case key == gocui.KeyCtrlSpace:
		v.EditNewLine()
	}
}
