// Package tableflip returns an ASCII-based emoticon depicting a
// person flipping a table out of rage or unflipping the table once
// thing have cooled down
package tableflip

import (
	"regexp"

	"github.com/handwritingio/deckard-bot/message"
)

// Plugin ...
type Plugin struct{}

var (
	// reTableFlip is the regexp variables for logic in HandleMessage
	reTableFlipChill = regexp.MustCompile(`(?i)^!table(flip|chill)$`)
)

// Usage prints detailed usage instructions for the plugin
func (p Plugin) Usage() string {
	return "`!tableflip` to get some table flipping action!\n" +
		"`!tablechill` to calm things down"
}

// HandleMessage is responsible for handling the incoming message
// and returning a response based on the message provides
func (p Plugin) HandleMessage(in message.Basic) (out message.Basic) {
	flipOrChill := reTableFlipChill.FindStringSubmatch(in.Text)[1]
	if flipOrChill == "flip" {
		out.Text = "(╯°□°）╯︵ ┻━┻"
	}
	if flipOrChill == "chill" {
		out.Text = "┬─┬ノ( º _ ºノ)"
	}
	return
}

// Command returns a list of commands the plugin provides
func (p Plugin) Command() []string {
	return []string{"!tableflip", "!tablechill"}
}

// OnInit returns an error if the plugin could not be started
func (p Plugin) OnInit() error {
	return nil
}

// Name is the name of the plugin
func (p Plugin) Name() string {
	return "TableFlip"
}

// Regexp returns the regexp of a message that should be handled by this plugin
func (p Plugin) Regexp() *regexp.Regexp {
	return reTableFlipChill
}
