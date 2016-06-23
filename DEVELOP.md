Developing Deckard
----------

## Development

**See [Contributing guidelines](CONTRIBUTING.md) for improvements, feature requests,
bug reports, etc.**

### Setup

Install Docker [here](https://docs.docker.com/installation/)

To develop plugins for Deckard, you can use **stdin/stdout**. This will be
helpful because you won't need to have a hook into Slack in order to develop
plugins that will interface with Deckard.

## Running with Docker for Development

Build the Docker image

	$ docker build -t deckard-bot .

Run the Docker Image

	$ docker run --rm -it deckard-bot

If you're running Deckard with Slack, you can run it in the background.

**You cannot run it in the background if you're using stdio connection**

## Building Plugins

1. Create a subpackage in the [plugins package](plugins) named after your plugin
1. Educate yourself on regex. Check out [regexr.com](http://www.regexr.com/), which is a great resource for testing out regex.
1. Your plugin must implement the [`Plugin` interface](plugins/plugin.go#L62).
**You must adhere to the instructions in the comments of the Plugin methods**
  1. `Name()` returns a string of the name of the plugin, which is used in the `!help` command.
  1. `Usage()` returns a string that holds the usage of the plugin. It's primarily used in the `!help` command.
  1. `Command()` returns a slice of strings that contain all the commands the plugin is configured for. It's primarily used in the `!help` command.
  1. `OnInit()` returns an error. This should take the place of the `init()` function. Returning an error will ensure
  the plugin isn't started and no messages will be routed to it. Use this to check configs and connections.
  1. `Regexp()` returns `*regexp.Regexp` and defines the regex of the commands that your plugin
	should receive. Your plugin will only be sent messages that match this regex. **This should match the return value
	from the `Command()` method.**
  1. `HandleMessage()` takes a `message.Basic` and returns a `message.Basic`. This is the primary method that handles the plugin's functionality. The returned `message.Basic` should be a response to the provided `message.Basic`.
1. Create tests for your plugin.

## Building Connections

Coming Soon...

## Installing Plugins

1. Build the plugin. See Examples [Here (simple)](plugins/tableflip/tableflip.go) or [Here (sample)](plugins/sample/sample.go).

1. Pick a connection, or build your own. See [stdio](connections/stdio/stdio.go) or [slack](connections/slack/slack.go).

1. **[Read the Docs!](https://godoc.org/github.com/handwritingio/deckard-bot)**

1. Create your main file and initialize your connection and plugins

	In `main.go`, add the following based on the plugin you created. Here's an example:

	```go
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
		// Using stdio (Terminal)
		conn := stdio.NewConnection()

		// Or using Slack (https://slack.com/)
		// You will need a Slack Custom Bot User: https://api.slack.com/bot-users#custom_bot_users
		/// conn := slack.NewConnection("SlackBotAPIKey")

		// 2. Create the bot using the connection and a list of plugins
		deckard := bot.New("Deckard", conn,
			&dice.Plugin{},
			&tableflip.Plugin{},
			&cats.Plugin{},
			&principles.Plugin{},
			// Any other plugins you create.
			// You can add or remove any of these
		)

		// 3. Start the bot!
		deckard.Go()
	}
	```
