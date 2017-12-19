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
	catFactURL = "https://catfact.ninja/fact"
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
		fact := getCatFact()
		if len(fact) == 0 {
			out.Text = "Sorry, I was unable to retrieve a cat fact for you :crying_cat_face:."
			return
		}
		out.Text = fact
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
		log.Errorf("Error getting cat fact response: %s", err)
		return ""
	}
	defer resp.Body.Close()
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading cat fact response body: %s", err)
		return ""
	}
	factResp := struct {
		Fact string `json:"fact"`
	}{}
	err = json.Unmarshal(raw, &factResp)
	if err != nil {
		log.Errorf("Error unmarshaling cat json: %s", err)
		return ""
	}
	return factResp.Fact
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
