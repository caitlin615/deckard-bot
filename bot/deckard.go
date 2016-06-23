/*
Package bot manages the creation of the chatbot and adding plugins
to the chatbot. The bot is created and configured with a connection.
The connection is the method in which the chatbot interfaces with humans
Examples are stdin/stdout, Slack, Hipchat, etc.

It configures the plugin list and starts the bot with the
message pump so messages can be sent and received through all
plugins

Basics

The chatbot sends and receives messages through message channels.
This is managed by the message pump. The received message is held in an inbox and sent
through the message pump to each plugin. Only the text of the message is sent
through the message pump. Any additional values send into the message pump are
held and returned with the response.

Plugins

Plugins should be members of the plugin package and require
three methods. Most important, they must have logic in order to handle
each message that they receive.

Connections

A Connection is a way for the chatbot to interface with a 3rd party communication tool (Slack, etc).
*/
package bot

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/handwritingio/deckard-bot/connection"
	"github.com/handwritingio/deckard-bot/log"
	"github.com/handwritingio/deckard-bot/message"
	"github.com/handwritingio/deckard-bot/plugins"
)

// Deckard is the object that handles all communication with the plugins and connections
type Deckard struct {
	Name             string
	Plugins          []plugins.Plugin
	conn             connection.Connection
	pluginInitResult chan pluginResult
}

type pluginResult struct {
	Plugin plugins.Plugin
	Error  error
}

// buildTime and version come from a linker flag during build time.
// buildTime is the UTC time that the binary is built
// version is the commit short hash
var (
	buildTime = "No Build Time Provided"
	version   = "No Build Version Provided"
)

func init() {
	log.Printf("Version: %s, Build Time: %s", version, buildTime)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)

	go func() {
		sig := <-ch
		log.Fatalf("Received Signal: %s", sig)
	}()
}

// AddPlugin call's the plugin's OnInit() method in an anonymous goroutine
// and puts the plugin itself and the result of OnInit() into a structure
// on the pluginInitResult channel to be handled as part of the main loop.
// This method is async to support plugins that require more startup time to
// not block the main loop of the bot
func (d *Deckard) AddPlugin(p plugins.Plugin) {
	go func() {
		d.pluginInitResult <- pluginResult{
			p, p.OnInit(),
		}
	}()
}

// New creates a new Bot with a name, new connection, and plugins.
func New(name string, conn connection.Connection, p ...plugins.Plugin) *Deckard {
	d := &Deckard{
		Name:             name,
		Plugins:          make([]plugins.Plugin, 0),
		pluginInitResult: make(chan pluginResult),
	}

	// Set the connection
	d.conn = conn

	// Add plugins
	for _, plugin := range p {
		d.AddPlugin(plugin)
	}

	log.Infof("Bot named %s Created", name)
	return d
}

// Go starts the TX/RX channels and starts the message pump
// Both are goroutines and exit the bot if anything enters
// the errorChannel
func (d *Deckard) Go() {
	errorChannel := make(chan error)
	rx, tx := d.conn.Start(errorChannel)
	go d.waitForPlugins()
	go d.messagePump(rx, tx)
	var err error
	err = <-errorChannel
	log.Fatal(err)
}

func (d *Deckard) waitForPlugins() {
	for {
		select {
		case result := <-d.pluginInitResult:
			fields := log.Fields{
				"Plugin": result.Plugin.Name(),
			}
			if result.Error != nil {
				fields["Error"] = result.Error.Error()
				log.WithFields(fields).Warn("Plugin Registration Failed")
			} else {
				d.Plugins = append(d.Plugins, result.Plugin)
				log.WithFields(fields).Info("Plugin Registered")
			}
		}
	}
}

// messagePump distributes messages via the RX channel
// to each plugin's HandleMessage method and returns
// HandleMessage message response to the TX channel
func (d *Deckard) messagePump(rx, tx message.BasicChannel) {
	for {
		select {
		case in := <-rx:
			if in.Text == "" {
				continue
			}

			// Check if the message is meant for internal plugin
			// and don't send it to other plugins if it's meant for internal
			// Messages meant for internal responses should not make it to plugins
			internalResponse := d.pluginInternal(in)
			if internalResponse.Finished {
				tx <- internalResponse
				continue
			}
			for _, p := range d.Plugins {
				if !p.Regexp().MatchString(in.Text) {
					log.Debugf("Message did not match regex for plugin %s... skipping", p.Name())
					continue
				}
				log.Infof("Message matches regex for plugin %s... sending message to plugin", p.Name())
				out := p.HandleMessage(in)
				out.ID = in.ID       // copy the id from the incoming message
				out.Finished = false // we're not done til we exit this loop
				if out.Text != "" {
					log.Infof("Incoming message: %#v", in)
					log.Infof("Outgoing message: %#v", out)
					tx <- out
				}
			}
			tx <- message.Basic{ID: in.ID, Text: "", Finished: true}
		}
	}
}
