package bot

import (
	"regexp"
	"strings"

	"github.com/handwritingio/deckard-bot/message"
)

var (
	// Regex for messages that should only be answered by Deckard
	// These messages won't be sent to plugins
	reDeckardHelp = regexp.MustCompile("(?i)^!help\\s*(\\S*)")
	reDeckardWho  = regexp.MustCompile("(?i)^!who$")
)

func (d *Deckard) pluginInternal(in message.Basic) message.Basic {
	switch {
	case reDeckardHelp.MatchString(in.Text):
		cmd := reDeckardHelp.FindStringSubmatch(in.Text)
		plugin := cmd[1]
		help := strings.Join(d.pluginHelp(plugin), "\n")
		return message.Basic{ID: in.ID, Text: help, Finished: true}

	case reDeckardWho.MatchString(in.Text):
		who := "Hello, I Am " + d.Name
		return message.Basic{ID: in.ID, Text: who, Finished: true}

	}
	return in
}

func (d *Deckard) pluginHelp(plugin string) (s []string) {
	if plugin != "" {
		// Return the specified plugin's usage
		for _, p := range d.Plugins {
			if strings.ToLower(plugin) == strings.ToLower(p.Name()) {
				s = append(s, "**Usage for `"+p.Name()+"` Plugin**")
				s = append(s, p.Usage())
				return
			}
		}
	} else {
		// Return the list of plugins and commands
		s = append(s, "*Here's a list of all known commands:*")
		for _, r := range d.Plugins {
			command := formatCommands(r.Command())
			s = append(s, "• Plugin *"+r.Name()+"* -- "+command)
		}
	}
	return
}
func formatCommands(cmd []string) string {
	s := []string{}
	for _, v := range cmd {
		s = append(s, "`"+v+"`")
	}
	return strings.Join(s, " ")
}
