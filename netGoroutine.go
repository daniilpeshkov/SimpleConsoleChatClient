package main

import (
	"context"
	"runtime"

	simpleTcpMessage "github.com/daniilpeshkov/go-simple-tcp-message"
)

func netReaderGoroutine(ctx context.Context, conn *simpleTcpMessage.ClientConn, outChan chan<- *simpleTcpMessage.Message) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		msg, _ := conn.RecieveMessage()
		outChan <- msg
		runtime.Gosched()
	}
}

func netWriterGoroutine(ctx context.Context, conn *simpleTcpMessage.ClientConn, inChan <-chan *simpleTcpMessage.Message) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-inChan:
			conn.SendMessage(msg)
			sys, _ := msg.GetField(TagSys)
			if sys[0] == SysMessage {
				text, _ := msg.GetField(TagMessage)
				unconfirmedMsgChan <- string(text)
			}
		}
	}
}
