package sample

import (
	"fmt"

	"github.com/handwritingio/deckard-bot/message"
)

func ExamplePlugin_HandleMessage() {
	p := new(Plugin)
	fmt.Println(p.HandleMessage(format("!sample")).Text)
	fmt.Println(rePluginRegexp.MatchString("!sample"))
	// Output:
	// Sample Plugin Output
	// true

}

func format(text string) message.Basic {
	return message.Basic{
		ID:       1,
		Text:     text,
		Finished: true,
	}
}
