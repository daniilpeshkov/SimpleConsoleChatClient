package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"runtime"

	simpleTcpMessage "github.com/daniilpeshkov/go-simple-tcp-message"
	"github.com/jroimartin/gocui"
)

const (
	ChatView  = "chat"
	InputView = "input"
	LoginView = "login"
)

var (
	curState   = LoginView
	msgInChan  = make(chan *simpleTcpMessage.Message, 10)
	msgOutChan = make(chan *simpleTcpMessage.Message, 10)
	clientConn *simpleTcpMessage.ClientConn
	globalCtx  = context.Background()
	cancelFunc context.CancelFunc
)

func main() {

	//defer profile.Start().Stop()

	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Panicln(err)
	}

	defer g.Close()
	g.FgColor = gocui.ColorYellow

	g.SetManager(gocui.ManagerFunc(ChatLayout),
		gocui.ManagerFunc(InputLayout),
		gocui.ManagerFunc(SetFocus),
		gocui.ManagerFunc(LoginLayout))

	g.ASCII = false
	// g.Highlight = true
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyF1, gocui.ModAlt, changeStyle); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(ChatView, gocui.MouseWheelUp, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	g.Cursor = true

	conn, err := net.Dial("tcp", "185.24.53.156:25565")
	if err != nil {
		log.Panicln("Server inaccessible")
	}
	clientConn = simpleTcpMessage.NewClientConn(conn)

	var ctx context.Context
	ctx, cancelFunc = context.WithCancel(globalCtx)

	go netReaderGoroutine(ctx, clientConn, msgInChan)
	go netWriterGoroutine(ctx, clientConn, msgOutChan)
	go func() {
		for {

			msg := <-msgInChan
			if sysMsg, ok := msg.GetField(TagSys); ok {
				switch {
				case sysMsg[0] == SysLoginResponse && sysMsg[1] == LOGIN_OK:
					curState = InputView
					g.Update(SetFocus)
				case sysMsg[0] == SysUserLoginNotiffication:
					name := string(sysMsg[2:])
					switch sysMsg[1] {
					case USER_CONNECTED:
						printChan <- Green + fmt.Sprintf("<%s connected>", name)
					case USER_DISCONECTED:
						printChan <- Red + fmt.Sprintf("<%s disconnected>", name)
					}
					g.Update(ChatLayout)
				}
			} else {

				name, _ := msg.GetField(TagName)
				text, _ := msg.GetField(TagText)
				printChan <- fmt.Sprintf("%s[%s]: %s%s", Yellow, string(name), White, text)
				g.Update(ChatLayout)
			}

			runtime.Gosched()
		}

	}()
	g.MainLoop()
}

func SetFocus(g *gocui.Gui) error {
	if curState == LoginView {
		g.SetViewOnTop(LoginView)
	} else {
		g.SetViewOnBottom(LoginView)
	}
	if curState == ChatView {
		g.SetViewOnTop(ChatView)
	}
	g.SetCurrentView(curState)
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	cancelFunc()
	return gocui.ErrQuit
}
func changeStyle(g *gocui.Gui, v *gocui.View) error {
	g.ASCII = !g.ASCII
	return nil
}
