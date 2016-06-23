/*
This is a sample main.go file which includes setup for a stdio connection
and plugins that do not require any configurations. If you'd like to build
your own plugins or connections, see our DEVELOP.md documentation
*/
package main

import (
	"github.com/handwritingio/deckard-bot/bot"

	"github.com/handwritingio/deckard-bot/plugins/cats"
	"github.com/handwritingio/deckard-bot/plugins/dice"
	"github.com/handwritingio/deckard-bot/plugins/principles"
	"github.com/handwritingio/deckard-bot/plugins/tableflip"

	"github.com/handwritingio/deckard-bot/connection/stdio"
)

func main() {
	// 1. Setup a new connection
	conn := stdio.NewConnection()

	// 2. Create the bot using the connection and a list of plugins
	deckard := bot.New("Deckard", conn,
		&dice.Plugin{},
		&tableflip.Plugin{},
		&cats.Plugin{},
		&principles.Plugin{},
	)

	// 3. Start the bot!
	deckard.Go()
}
