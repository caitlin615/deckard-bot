# Contributing

Report issues via [Github Issues](https://github.com/handwritingio/deckard-bot/issues).
The quickest way to resolve issues is to open your own Pull Request.

If you'd like to make changes to this code, submit a Pull Request.

## Running Tests

The commands to run tests, vet and lint are:

    go test ./...
    go vet ./...
    golint ./...

If you're using docker-compose, you may run:

    make test vet lint
