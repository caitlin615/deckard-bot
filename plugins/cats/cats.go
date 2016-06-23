// Package cats is a plugin that queries some cat API's and returns responses.
// It queries a cat fact API and a cat image/gif API
package cats

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/handwritingio/deckard-bot/log"
	"github.com/handwritingio/deckard-bot/message"
)

// Plugin ...
type Plugin struct{}

var (
	// reCats is the regexp variables for logic in HandleMessage
	reCats     = regexp.MustCompile("(?i)^!cat")
	reCatsType = regexp.MustCompile("(?i)^!cat (\\w+)$")

	catImgURL  = "http://thecatapi.com/api/images/get?format=src&size=med&type="
	catFactURL = "http://catfacts-api.appspot.com/api/facts?number=1"
)

// Command returns a list of commands the plugin provides
func (p Plugin) Command() []string {
	return []string{"!cat"}
}

// Usage prints detailed usage instructions for the plugin
func (p Plugin) Usage() string {
	return "`!cat image` to generate new cat photo\n" +
		"`!cat gif` to generate new cat gif\n" +
		"`!cat fact` to generate new cat fact"
}

// Regexp returns the regexp of a message that should be handled by this plugin
func (p Plugin) Regexp() *regexp.Regexp {
	return reCats
}

// HandleMessage is responsible for handling the incoming message
// and returning a response based on the message provides
func (p Plugin) HandleMessage(in message.Basic) (out message.Basic) {
	chunks := reCatsType.FindStringSubmatch(in.Text)
	log.Debug("Arguments: ", len(chunks))
	// Define the option and argument strings, if supplied
	var cmd string
	if len(chunks) > 1 {
		cmd = strings.ToLower(chunks[1])

		log.Debugf("Cat cmd: %s\n", cmd)
	}

	switch cmd {
	case "gif":
		out.Text = getCatImage("gif")
	case "image":
		out.Text = getCatImage("jpg")
	case "fact":
		out.Text = getCatFact()
	default:
		out.Text = p.Usage()
	}

	return
}

// OnInit returns an error if the plugin could not be started
func (p Plugin) OnInit() error {
	return nil
}

// Name is the name of the plugin
func (p Plugin) Name() string {
	return "Cats"
}

// getCatFact returns a random cat fact from the catFactURL
func getCatFact() string {
	resp, err := http.Get(catFactURL)
	if err != nil {
		log.Debugf("Error getting cat fact response: %s", err)
	}
	defer resp.Body.Close()
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Debugf("Error reading cat fact response body: %s", err)
	}
	facts := struct {
		Facts   []string `json:"facts"`
		Success string   `json:"success"`
	}{}
	err = json.Unmarshal(raw, &facts)
	if err != nil {
		log.Debugf("Error unmarshaling cat json: %s", err)
	}
	if facts.Success != "true" {
		log.Debug("Cat facts json did not return success")
	}

	var s []string
	for _, k := range facts.Facts {
		s = append(s, k)
	}
	return strings.Join(s, "\n")
}

// getCatImage returns a random cat git or png from the catImageURL
func getCatImage(t string) string {
	req, err := http.NewRequest("GET", catImgURL+t, nil)
	if err != nil {
		log.Debugf("Error getting cat photo request: %s", err)
	}
	transport := http.Transport{}
	resp, err := transport.RoundTrip(req)
	if err != nil {
		log.Debugf("error with cat response: %s", err)
	}

	if resp.StatusCode != 200 && resp.StatusCode != 302 {
		log.Debug("Failed with status: ", resp.Status)
	}

	return resp.Header.Get("Location")
}
