package main

import (
	"strings"

	simpleTcpMessage "github.com/daniilpeshkov/go-simple-tcp-message"
	"github.com/jroimartin/gocui"
)

func InputLayout(g *gocui.Gui) error {
	w, h := g.Size()
	x0, y0, x1, y1 := 2, h-int(float32(h)*InputHeightProcentage), w-2, h-1

	if y0 < 0 {
		y0 = 0
	}
	if y1 < 0 {
		y1 = 0
	}
	v, err := g.SetView(InputView, x0, y0, x1, y1)
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
			msg.AppendField(TagText, []byte(text))
			msgOutChan <- msg
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
	case key == gocui.KeyEnter && mod == gocui.ModAlt:
		v.EditNewLine()
	case key == gocui.KeyCtrlSpace:
		v.EditNewLine()
	}
}
