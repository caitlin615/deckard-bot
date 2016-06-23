// Package connection creates the interface for all services that require a connection to the chatbot.
// Each connection must have a Start method that creates and returns a transmit and receive channel,
// which are used to communicate with the plugins.
//
// Messages sent to the chatbot via the connection interface will be routed through the rx channel. The chatbot will distribute messages
// on this channel to all plugins. The plugins then have the opportunity to respond through the tx channel.
//
// Messages sent from the chatbot from the plugins are sent into the tx channel. The chatbot takes messages
// from the tx channel and returns it to the connection interface.
package connection

import "github.com/handwritingio/deckard-bot/message"

// Connection interface has a Start method for creating the connection
// two basic channels for transmitting and receiving messages
type Connection interface {
	Start(chan error) (rx, tx message.BasicChannel)
}
