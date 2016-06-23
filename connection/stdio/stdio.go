package stdio

import (
	"bufio"
	"os"
	"strings"

	"github.com/handwritingio/deckard-bot/log"
	"github.com/handwritingio/deckard-bot/message"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	colorRedBold = "\x1b[1;31m"
	colorYellow  = "\x1b[0;33m"
	colorReset   = "\x1b[0m"
)

// Connection provides an interface for storing the inbox for received messages via stdio connection type
type Connection struct {
	Inbox map[int]message.Basic
}

// NewConnection creates a new StdIO object with an inbox to keep track of messages
func NewConnection() *Connection {
	s := &Connection{
		make(map[int]message.Basic),
	}
	return s
}

// Start creates two message channels to send and receive messages.
// It will start two goroutines to listen and send on these channels
func (s *Connection) Start(errorChannel chan error) (rx, tx message.BasicChannel) {
	rx = make(message.BasicChannel)
	tx = make(message.BasicChannel)
	go s.startRX(rx, errorChannel)
	go s.startTX(tx, errorChannel)
	return rx, tx
}

// startRX will read lines off stdin and add them to the inbox and RX channel
func (s *Connection) startRX(rx message.BasicChannel, errorChannel chan error) {
	reader := bufio.NewReader(os.Stdin)
	counter := 0
	for {
		line, err := reader.ReadString('\n')
		line = strings.Trim(line, "\n")
		// log.Debug("Got line: ", line)
		if err != nil {
			errorChannel <- err
			break
		}
		msg := message.Basic{ID: counter, Text: line, Finished: false}
		s.Inbox[counter] = msg
		rx <- msg
		counter++
	}
}

// startRX will read lines off the TX channel and write it to stdout
func (s *Connection) startTX(tx message.BasicChannel, errorChannel chan error) {
	writer := bufio.NewWriter(os.Stdout)
	for {
		select {
		case msg := <-tx:
			if msg.Text != "" {
				msgTTY := "DECKARD RESPONSE: " + msg.Text
				if terminal.IsTerminal(int(os.Stdin.Fd())) {
					msgTTY = colorRedBold + "DECKARD RESPONSE: " + colorYellow + msg.Text + colorReset
				}
				_, err := writer.Write([]byte(msgTTY + "\n\n"))
				if err != nil {
					errorChannel <- err
					break
				}
				err = writer.Flush()
				if err != nil {
					errorChannel <- err
					break
				}
			}
			if msg.Finished {
				delete(s.Inbox, msg.ID)
			}
			log.Debug("inbox size: ", len(s.Inbox))
		}
	}
}
