// Package message contains the structure for basic messages that
// are allowed to be sent through the RX and TX channels.
package message

// Basic implements the message structure that is added to the inbox and
// sent to the plugins
type Basic struct {
	ID       int    `json:"id"`
	Text     string `json:"text"`
	Finished bool
}

// BasicChannel is a channel that accepts Basic messages.
// All transmit and receive channels must fit this type
type BasicChannel chan Basic
