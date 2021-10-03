package main

import (
	"github.com/jroimartin/gocui"
)

func FileDialogLayout(g *gocui.Gui) error {

	w, h := g.Size()
	x0, y0, x1, y1 := w/2-18, h/2-1, w/2+18, h/2+1

	if y0 < 0 {
		y0 = 0
	}
	if y1 < 0 {
		y1 = 0
	}
	v, err := g.SetView(FileDialogView, x0, y0, x1, y1)
	if err != nil {
		v.Editor = gocui.EditorFunc(fileDialogEditor)
		v.Autoscroll = false
		v.Title = "Enter file path"
		v.Editable = true
		v.Frame = true
		v.Overwrite = true
		v.Wrap = false
		v.FgColor = gocui.ColorWhite
		v.BgColor = gocui.ColorRed
	}

	return nil
}

func fileDialogEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
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
		// myName = strings.TrimSpace(v.Buffer())
		// if len(myName) != 0 {
		// 	msg := simpleTcpMessage.NewMessage()
		// 	msg.AppendField(TagSys, []byte{SysLoginRequest})
		// 	msg.AppendField(TagName, []byte(myName))
		// 	msgOutChan <- msg
		// 	v.Clear()
		// 	v.SetCursor(0, 0)
		// }
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
