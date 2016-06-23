/*
Package plugins contains all plugins that interface with the chatbot.

Each plugin must implement the Plugin interface.

Here's an example of an implementation of the Plugin interface

 package sample

 import (
 	"regexp"
 	"github.com/handwritingio/deckard-bot/log"
 	"github.com/handwritingio/deckard-bot/message"
 )

 // Plugin initializes the Plugin interface
 type Plugin struct{}

 // Usage prints detailed usage instructions for the plugin
 func (p *Plugin) Usage() string {
 	return "`!sample` is a sample plugin"
 }

 // Command returns a list of commands the plugin provides
 func (p *Plugin) Command() []string {
 	return []string{"!sample"}
 }

 // OnInit returns an error if the plugin could not be started
 func (p *Plugin) OnInit() error {
 	return nil
 }

 // Name returns the name of the plugin
 func (p *Plugin) Name() string {
 	return "Sample"
 }

 // Regexp returns the regexp of a message that should be handled by this plugin
 func (p Plugin) Regexp() *regexp.Regexp {
 	return regexp.MustCompile(`(?i)^!sample`)
 }

 // HandleMessage is responsible for handling the incoming message
 // and returning a response based on the message provides
 func (p *Plugin) HandleMessage(in message.Basic) (out message.Basic) {
 	log.Debug("Sample plugin matched")
 	out.Text = "Sample Plugin Output"
 	return
 }
*/
package plugins

import (
	"regexp"

	"github.com/handwritingio/deckard-bot/message"
)

// Plugin provides the interface for building a plugin
type Plugin interface {
	// Name returns a string that should be the name of the plugin.
	// This can then be used in other places, such as outputting
	// the usage for all plugins.
	Name() string

	// Usage method can be used to hold the plugin's usage syntax.
	// It returns a string.  You can call this in the HandleMessage
	// to send it in the response.
	Usage() string

	// Command lists the strings that this plugin will respond
	// to. It's good practice (but optional) to use the prefix
	// "!" before your commands if you're not making a sneaky
	// plugin. You can provide 1 or more strings.
	Command() []string

	// HandleMessage method is what handles the logic for the chatbot's
	// responses to messages. There should be logic built into this method.
	HandleMessage(message.Basic) message.Basic

	// OnInit is called by the bot when the plugin is
	// initialized.  You can also call other functions here that should be
	// run once when the plugin starts. If the plugin returns an error from
	// this method, the server will not send chat traffic through it, and
	// HandleMessage will never be called
	OnInit() error

	// Regexp returns the regexp compiled regular expression type
	// that is used to match a message to a plugin. This regex should
	// be as generic as possible for what the plugin requires.
	Regexp() *regexp.Regexp
}
