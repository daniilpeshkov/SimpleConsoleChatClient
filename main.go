package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"syscall"
	"time"

	simpleTcpMessage "github.com/daniilpeshkov/go-simple-tcp-message"
	"github.com/jroimartin/gocui"
)

const (
	ChatView       = "chat"
	InputView      = "input"
	LoginView      = "login"
	FileDialogView = "file"
)

var (
	curState            = LoginView
	msgInChan           = make(chan *simpleTcpMessage.Message, 10)
	msgOutChan          = make(chan *simpleTcpMessage.Message, 10)
	myName              string
	unconfirmedMsgChan  = make(chan string, 10)
	unconfirmedFileChan = make(chan string, 10)
	clientConn          *simpleTcpMessage.ClientConn
	globalCtx           = context.Background()
	cancelFunc          context.CancelFunc
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
		gocui.ManagerFunc(LoginLayout),
		gocui.ManagerFunc(FileDialogLayout))

	g.ASCII = false
	// g.Highlight = true
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyF1, gocui.ModAlt, changeStyle); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlF, gocui.ModNone, openFileDialog); err != nil {
		log.Panicln(err)
	}
	g.Cursor = true
	//
	// conn, err := net.Dial("tcp", "127.0.0.1:25565")
	// conn, err := net.Dial("tcp", "185.24.53.156:25565")
	conn, err := net.Dial("tcp", "192.168.1.45:25565")
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

			timeB, _ := msg.GetField(TagTime)
			msgTime := new(time.Time)
			msgTime.UnmarshalBinary(timeB)
			localMsgTime := msgTime.Local()
			timeFmt := Cyan + localMsgTime.Format(time.RFC822) + ": "
			sysMsg, _ := msg.GetField(TagSys)
			nameB, _ := msg.GetField(TagName)
			name := string(nameB)
			switch {
			case sysMsg[0] == SysLoginRequest && sysMsg[1] == LOGIN_OK:
				curState = InputView
				printChan <- timeFmt + Green + "<connected>"
				g.Update(SetFocus)
			case sysMsg[0] == SysUserLoginNotiffication:
				switch sysMsg[1] {
				case USER_CONNECTED:
					printChan <- timeFmt + Green + fmt.Sprintf("<%s connected>", name)
				case USER_DISCONECTED:
					printChan <- timeFmt + Red + fmt.Sprintf("<%s disconnected>", name)
				}
				g.Update(ChatLayout)
			case sysMsg[0] == SysMessage:
				if len(sysMsg) == 1 { //other message
					text, _ := msg.GetField(TagMessage)
					printChan <- timeFmt + fmt.Sprintf("%s[%s]: %s%s", Yellow, string(name), White, text)
					g.Update(ChatLayout)
				} else { //confirmed message
					if sysMsg[1] == MESSAGE_SENT {
						text := <-unconfirmedMsgChan
						printChan <- timeFmt + fmt.Sprintf("%s[%s]: %s%s", Yellow, string(name), White, text)
						g.Update(ChatLayout)
					}
				}
			case sysMsg[0] == SysFile:
				if len(sysMsg) == 1 {
					fileNameB, _ := msg.GetField(TagFileName)
					fileContB, _ := msg.GetField(TagFile)

					fileName := string(fileNameB)
					os.WriteFile("./"+fileName, fileContB, syscall.O_CREAT|syscall.O_EXCL)
					printChan <- timeFmt + fmt.Sprintf("%s[%s]: %ssent file %s", Yellow, string(myName), Blue, fileName)
					g.Update(ChatLayout)
				} else { //confirmed message

					if sysMsg[1] == FILE_SENT {
						fileName := <-unconfirmedFileChan
						printChan <- fmt.Sprintf("%s%s%s sent", timeFmt, Yellow, fileName)
						g.Update(ChatLayout)
					}
				}
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
	if curState == FileDialogView {
		g.SetViewOnTop(FileDialogView)
	} else {
		g.SetViewOnBottom(FileDialogView)
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
func openFileDialog(g *gocui.Gui, v *gocui.View) error {
	if curState == InputView {
		curState = FileDialogView
		g.Update(SetFocus)
		// g.Update(FileDialogLayout)
	} else if curState == FileDialogView {
		curState = InputView
		g.Update(SetFocus)
	}
	return nil
}
