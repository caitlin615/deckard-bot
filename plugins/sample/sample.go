// Package sample is a sample implementation of a plugin.
// Use this as a place to start to build a plugin
package sample

import (
	"regexp"

	"github.com/handwritingio/deckard-bot/log"
	"github.com/handwritingio/deckard-bot/message"
)

// Plugin initializes the interface
type Plugin struct{}

var rePluginRegexp = regexp.MustCompile(`(?i)^!sample`)

// Usage prints detailed usage instructions
func (p *Plugin) Usage() string {
	return "`!sample` is a sample plugin"
}

// Command returns a list of commands the plugin provides
func (p *Plugin) Command() []string {
	return []string{"!sample"}
}

// OnInit returns an error if the plugin could not be started
// OnInit doesn't do anything here as nothing needs to be initialized for this plugin
func (p *Plugin) OnInit() error {
	return nil
}

// Name returns the name of the plugin
func (p *Plugin) Name() string {
	return "Sample"
}

// Regexp returns the regexp of a message that should be handled by this plugin
func (p Plugin) Regexp() *regexp.Regexp {
	return rePluginRegexp
}

// HandleMessage is responsible for handling the incoming message
// and returning a response based on the message provided
func (p *Plugin) HandleMessage(in message.Basic) (out message.Basic) {
	log.Debug("Sample plugin matched")
	out.Text = "Sample Plugin Output"
	return
}
