/*
Package slack is a Connection to the Slack Real Time Messaging API
(https://api.slack.com/rtm). To use this connection, you will need
to create a custom bot user. See https://api.slack.com/bot-users#custom_bot_users.

Once you've created a bot user, you will need to initialize the Slack connection
using the API key for this bot user.

 slackConnection := NewConnection("MySlackBotAPIKey")
*/
package slack

import (
	"encoding/json"
	"errors"

	"github.com/handwritingio/deckard-bot/log"
	"github.com/handwritingio/deckard-bot/message"

	"golang.org/x/net/websocket"
)

// Connection provides an interface for storing the Slack API key and the inbox for storing received messages
type Connection struct {
	Token string
	Inbox map[int]Message
}

// Message provides the interface for all Slack messages
type Message struct {
	message.Basic
	Type      string `json:"type"`
	Channel   string `json:"channel"`
	User      string `json:"user"`
	Timestamp string `json:"ts"`
}

// NewConnection returns a new Connection to Slack
func NewConnection(slackAPIKey string) *Connection {
	return &Connection{
		slackAPIKey,
		make(map[int]Message),
	}
}

// Start creates the connection for the Slack RTM and creates the transmit and receive
// goroutines that listen and send messages through the tx and rx channels
func (s *Connection) Start(errorChannel chan error) (rx, tx message.BasicChannel) {
	rx = make(message.BasicChannel)
	tx = make(message.BasicChannel)
	wssurl, err := getWSSUrl(s.Token)
	if err != nil {
		errorChannel <- err
		panic(err) // TODO
	}
	// Connect to the websocket
	ws, err := websocket.Dial(wssurl, "", "http://localhost/")
	if err != nil {
		log.Fatal(err)
	}
	// defer ws.Close() // TODO

	// start message Id generator
	msgChan := messageIDGen(0, 1)

	// run keepalive to keep the websocket connection running
	go keepalive(ws, msgChan)

	// start the RX and TX methods
	go s.startRX(ws, rx, errorChannel)
	go s.startTX(ws, tx, msgChan, errorChannel)
	return rx, tx
}

// startRX listens to all Slack messages. It adds all messages of type 'message' to the inbox and
// adds the message with text only (m.Basic) to the rx channel. Messages send into the rx channel
// are sent to the messagePump, which sends the message to each plugin
func (s *Connection) startRX(ws *websocket.Conn, rx message.BasicChannel, errorChannel chan error) {
	// get info of bot
	BotID, err := apiTokenAuthTest(s.Token)
	if err != nil {
		errorChannel <- err
	}

	// run infinite loop for receiving messages
	counter := 0
	for {
		var raw json.RawMessage
		err := websocket.JSON.Receive(ws, &raw)
		if err != nil {
			errorChannel <- err
		}

		var event struct {
			Type    string          `json:"type"`
			Error   json.RawMessage `json:"error"`
			ReplyTo int             `json:"reply_to"`
		}
		err = json.Unmarshal(raw, &event)
		if err != nil {
			errorChannel <- err
		}

		switch event.Type {
		case "":
			log.Debug("Acknowledge message: ", event.ReplyTo)
		case "pong", "presence_change", "user_typing", "reconnect_url":
			continue
		case "hello":
			// Send response to hello straight into websocket without going through messagePump
			log.Debug("Hello Event: ", event.Type)
		case "message":
			var m Message
			err = json.Unmarshal(raw, &m)
			if err != nil {
				errorChannel <- err
			}
			m.Basic.Text = formatSlackMsg(m.Basic.Text)
			log.Debugf("Full msg: %v\n", m)

			// if the message is not from the configured Bot
			// we don't want the bot responding to its own messages
			if m.User != BotID {
				// returns response string
				m.Basic.ID = counter
				s.Inbox[counter] = m
				rx <- m.Basic
				counter++
			}
		}
	}
}

// startTX is responsible for listening on the tx channel and sending all non-blank messages back through the
// websocket connection. The outgoing message is reassembled from the text from the tx channel and the rest of
// the original message attributes. Since this is the Slack startTX, it add a mention before the text to alert
// user that sent the original message that the bot has responded
func (s *Connection) startTX(ws *websocket.Conn, tx message.BasicChannel, msgChan <-chan int, errorChannel chan error) {
	for {
		select {
		case msg := <-tx:
			// handle everything except blank messages
			if msg.Text != "" {

				out, ok := s.Inbox[msg.ID]
				if ok != true {
					errorChannel <- errors.New("unknown id")
				}
				// get the UserId from the message sent
				// Add it to the beginning of the text
				msgUser := "<@" + out.User + ">: "
				out.Text = msgUser + msg.Text
				id := <-msgChan
				out.ID = id

				// send response struct
				err := websocket.JSON.Send(ws, &out)
				if err != nil {
					errorChannel <- err
				}
				if msg.Finished {
					delete(s.Inbox, msg.ID)
				}
				log.Debug("inbox size: ", len(s.Inbox))
			}
		}
	}
}
