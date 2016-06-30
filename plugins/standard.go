package plugins

import (
	"github.com/handwritingio/deckard-bot/plugins/cats"
	"github.com/handwritingio/deckard-bot/plugins/dice"
	"github.com/handwritingio/deckard-bot/plugins/principles"
	"github.com/handwritingio/deckard-bot/plugins/tableflip"
)

// Standard returns all the pre-installed standard plugins that require no initialization
func Standard() []Plugin {
	return []Plugin{
		&dice.Plugin{},
		&principles.Plugin{},
		&tableflip.Plugin{},
		&cats.Plugin{},
	}
}
