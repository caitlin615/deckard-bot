// Package dice is a plugin that will roll dice based on the number of dice
// you specify and the number of sides on each dice (ex: roll 7 6-sized dice)
package dice

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/handwritingio/deckard-bot/log"
	"github.com/handwritingio/deckard-bot/message"
)

// Plugin ...
type Plugin struct{}

var (
	// Dice Regexp variables for logic in HandleMessage
	reDice = regexp.MustCompile(`(?i)^!dice`)
	reRoll = regexp.MustCompile(`(?i)^!dice (\d{1,5})d(\d{1,5})$`)
	// Seeded random number generator
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// Regexp returns the regexp of a message that should be handled by this plugin
func (p Plugin) Regexp() *regexp.Regexp {
	return reDice
}

// Command returns a list of commands the plugin provides
func (p Plugin) Command() []string {
	return []string{"!dice"}
}

// Usage prints detailed usage instructions for the plugin
func (p Plugin) Usage() string {
	return "`!dice nDm` where n is the number of dice and m is the number of sides " +
		"(e.g. `!dice 2d6` means roll 2 6-sided dice)"
}

// HandleMessage is responsible for handling the incoming message
// and returning a response based on the message provides
func (p Plugin) HandleMessage(in message.Basic) (out message.Basic) {
	log.Debug("dice matched...")
	chunks := reRoll.FindStringSubmatch(in.Text)
	if len(chunks) != 3 { // invalid command
		out.Text = p.Usage()
		return
	}

	// We don't need to worry about the strings not being ints because the regex only
	// matches ints, and not negative ints or other chars
	nDice, _ := strconv.Atoi(chunks[1])
	nSides, _ := strconv.Atoi(chunks[2])
	log.Debugln("dice:", nDice, "sides:", nSides)

	// We do, however have to worry about zeros:
	if nSides == 0 {
		out.Text = "I can't roll a 0-sided die!"
		return
	}
	if nDice == 0 {
		out.Text = "I can't roll 0 dice!"
		return
	}
	out.Text = fmt.Sprintf("you rolled `%d`", roll(nDice, nSides))

	return
}

// OnInit returns an error if the plugin could not be started
func (p Plugin) OnInit() error {
	return nil
}

// Name is the name of the plugin
func (p Plugin) Name() string {
	return "Dice"
}

func roll(dice, sides int) (total int) {
	total = 0
	for i := 0; i < dice; i++ {
		// Intn returns [0, n), so we add 1 to get [1, n]
		total += rng.Intn(sides) + 1
	}
	return total
}
