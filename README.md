# deckard-bot
[![CircleCI](https://circleci.com/gh/handwritingio/deckard-bot.svg?style=shield)](https://circleci.com/gh/handwritingio/deckard-bot)
[![GoDoc](https://godoc.org/github.com/handwritingio/deckard-bot?status.png)](https://godoc.org/github.com/handwritingio/deckard-bot)

Welcome to the home to [Handwriting.io](https://handwriting.io)'s chatbot, **Deckard**.

## About

**Deckard** is a chatbot library that can help you simplify your life. He was created as a way to help
[Handwriting.io](https://handwriting.io) developers work more efficiently and transparently.

Deckard has two connections built-in. You can talk to him through Slack or through a terminal (stdin/stdout).

## Installing

Download the source code

```
go get github.com/handwritingio/deckard-bot
```

Alternatively, `git clone` and `go build` to run from source.

And install

```
go install github.com/handwritingio/deckard-bot
```

## Getting started

See [PLUGINS.md](PLUGINS.md) for requirements for each of Deckard's built-in plugins.

### Want to run Deckard using Slack?

**First**, create a Slack bot:

* Deckard connects to Slack using the [Slack RTM (Real Time Messaging API)](https://api.slack.com/rtm).
* You will need to [Create a Slack Bot](https://api.slack.com/bot-users) and get your bot's Slack API Token.

**Then** initialize the Slack connection in your `main.go`

```go
import "github.com/handwritingio/deckard-bot/connection/slack"

func main() {
  ...

  slackConn := slack.NewConnection("mySlackBotAPIToken")

  ...
}
```

### What to run Deckard using terminal?

**First** initialize the Stdio connection in your `main.go`

```go
import "github.com/handwritingio/deckard-bot/connection/stdio"

func main() {
  ...

  stdioConn := slack.NewConnection()

  ...
}
```

That's it!


### Initializing Plugins and create the Bot

Once you create a connection, you should initialize plugins and create your bot.
We provided a few example plugins that you can use, but feel free to create your own!

```go
import (
	"github.com/handwritingio/deckard-bot/bot"

	"github.com/handwritingio/deckard-bot/plugins/cats"
  // Make sure to import your custom plugins
)

func main() {
  ...

  deckard := bot.New("Deckard", myConnection,
    &cats.Plugin{},
    // Any other plugins
  )
  ...
}
```

### Put it all together and start the bot!

Now that you've created a connection and initialized your plugins, you should create
a new bot and start the bot. Here's an example of what a main.go file can look like:

```go
package main

import (
	"github.com/handwritingio/deckard-bot/bot"

  "github.com/handwritingio/deckard-bot/plugins"
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
```
## Running Deckard

```
go build .
./deckard-bot
```

## Developing

See [DEVELOP.md](DEVELOP.md)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

## License

See [LICENSE](LICENSE)

## Version Numbers

Version numbers for this package will follow standard
[semantic versioning](http://semver.org/).

## Issues

Please open an issue on [Github](https://github.com/handwritingio/deckard-bot/issues)
and we will look into it as soon as possible
