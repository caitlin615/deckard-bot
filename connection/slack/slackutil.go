package slack

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/handwritingio/deckard-bot/config"
	"github.com/handwritingio/deckard-bot/log"

	"golang.org/x/net/websocket"
)

var reSlackFormat = regexp.MustCompile(`<https?:\/\/(\S+)\|(\S+)>`)

// formatSlackMsg remove automatic formatting done by Slack before it gets to the rx channel
func formatSlackMsg(msg string) string {
	formatMatch := reSlackFormat.MatchString(msg)
	if !formatMatch {
		return msg
	}

	log.Debug("Slack formatting detected: ", formatMatch)
	log.Debugf("Original message: %s", msg)
	text := reSlackFormat.FindStringSubmatch(msg)
	if len(text) != 3 {
		// Something went wrong, the regex has 2 capture groups, so it should only
		// return a slice of length 3
		log.Debug("Error parsing link, not enough arguments")
	}
	fixed := reSlackFormat.ReplaceAllString(msg, text[2])
	log.Debugf("Fixed message: %s", fixed)
	return fixed
}

// messageIDGen creates a channel for generating the messageId needed
// to send back a message
func messageIDGen(start int, step int) <-chan int {
	c := make(chan int)

	go func() {
		counter := start
		for {
			c <- counter
			counter += step
		}
	}()
	return c
}

// keepalive is responsible for pinging the Slack websocket connection
// in order to keep it alive. It will respond with `pong` if the ping was successful.
func keepalive(ws *websocket.Conn, msgChan <-chan int) error {

	for {
		msgID := <-msgChan
		err := websocket.JSON.Send(ws, struct {
			ID   int    `json:"id"`
			Type string `json:"type"`
		}{ID: msgID, Type: "ping"})
		if err != nil {
			log.Error(err)
			return err
		}
		// wait 15 second
		time.Sleep(15 * time.Second)
	}
}

// getWSSUrl returns the websocket url from a json payload
// based on the supplied API auth token
func getWSSUrl(token string) (string, error) {
	resp, err := http.Get(config.SlackAPIURL + "/rtm.start?token=" + token)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// no need to check status code, errors returned inside json
	// https://api.slack.com/methods/rtm.start

	var rtm struct {
		URL   string `json:"url"`
		Ok    bool   `json:"ok"`
		Error string `json:"error"`
	}

	err = json.Unmarshal(raw, &rtm)
	if err != nil {
		return "", err
	}

	// Error reponses based on an ok: false
	if !rtm.Ok {
		switch rtm.Error {
		case "migration_in_progress":
			err = errors.New("Team is being migrated between servers.")
		case "not_authed":
			err = errors.New("No authentication token provided.")
		case "invalid_auth":
			err = errors.New("Invalid authentication token.")
		case "account_inactive":
			err = errors.New("Authentication token is for a deleted user or team.")
		default:
			err = errors.New("Something else went wrong. rtm.start status not ok. See https://api.slack.com/methods/rtm.start")
		}
	}

	log.WithFields(log.Fields{
		"url":      rtm.URL,
		"response": rtm.Ok,
	}).Debug("Returned json from rtm.start call")

	return rtm.URL, err
}

// apiTokenAuthTest tests the connection based on the method endpoint of the Slack API
// This is mostly used to get the username and id for the account associated with the API token
func apiTokenAuthTest(token string) (botID string, err error) {

	resp, err := http.Get(config.SlackAPIURL + "/auth.test?token=" + token)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var authTestResp struct {
		Ok     bool   `json:"ok"`
		User   string `json:"user"`
		UserID string `json:"user_id"`
		Error  string `json:"error"`
	}
	err = json.Unmarshal(raw, &authTestResp)
	if err != nil {
		return
	}

	// Error reponses based on an ok: false
	if !authTestResp.Ok {
		switch authTestResp.Error {
		case "not_authed":
			err = errors.New("No authentication token provided.")
		case "invalid_auth":
			err = errors.New("Invalid authentication token.")
		case "account_inactive":
			err = errors.New("Authentication token is for a deleted user or team.")
		default:
			err = errors.New("Something else went wrong. auth.test status not ok. See https://api.slack.com/methods/auth.test")
		}
	}
	log.Printf("Bot Info: %#v", authTestResp)
	botID = authTestResp.UserID
	return
}

type userInfo struct {
	Members []member
	Ok      bool   `json:"ok"`
	Error   string `json:"error"`
}

type member struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	IsBot bool   `json:"is_bot"`
}

type channelInfo struct {
	Channels []channel
	Ok       bool   `json:"ok"`
	Error    string `json:"error"`
}

type channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// getUserById returns the Slack user Id based on the supplied username.
// This is helpful because the RTM only sends the user Id, and sometimes
// you'll need to convert that to the username.
func (s *Connection) getUserByID(userName string) (string, error) {
	resp, err := http.Get(config.SlackAPIURL + "/users.list?token=" + s.Token)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var userResp userInfo
	err = json.Unmarshal(raw, &userResp)
	if err != nil {
		return "", err
	}

	// Error reponses based on an ok: false
	if !userResp.Ok {
		switch userResp.Error {
		case "user_not_found":
			err = errors.New("Value passed for user was invalid.")
		case "user_not_visible":
			err = errors.New("The requested user is not visible to the calling user")
		case "not_authed":
			err = errors.New("No authentication token provided.")
		case "invalid_auth":
			err = errors.New("Invalid authentication token.")
		case "account_inactive":
			err = errors.New("Authentication token is for a deleted user or team.")
		default:
			err = errors.New("Something else went wrong. users.info status not ok. See https://api.slack.com/methods/users.info")
		}
	}
	// loop through all members and return the user Id if the
	// associated username matches the supplied username
	var retValue string
	user := userResp.Members
	for i := 0; i < len(user); i++ {
		item := user[i]
		if item.Name == strings.ToLower(userName) {
			retValue = item.ID
		}
	}
	if retValue == "" {
		err = errors.New("User not found!")
	}
	return retValue, err
}

// getChannelByName returns the Slack channel Id based on the supplied channel name.
// This is helpful because the RTM can only send messages to a channel based on the
// channel ID.
func (s *Connection) getChannelByName(channel string) (string, error) {

	resp, err := http.Get(config.SlackAPIURL + "/channels.list?token=" + s.Token)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var c channelInfo
	err = json.Unmarshal(raw, &c)
	if err != nil {
		return "", err
	}

	// Error reponses based on an ok: false
	if !c.Ok {
		switch c.Error {
		case "not_authed":
			err = errors.New("No authentication token provided.")
		case "invalid_auth":
			err = errors.New("Invalid authentication token.")
		case "account_inactive":
			err = errors.New("Authentication token is for a deleted user or team.")
		default:
			err = errors.New("Something else went wrong. users.info status not ok. See https://api.slack.com/methods/users.info")
		}
	}
	// loop through all members and return the user Id if the
	// associated username matches the supplied username
	var retValue string
	channelList := c.Channels
	for i := 0; i < len(channelList); i++ {
		item := channelList[i]
		if item.Name == strings.ToLower(channel) {
			retValue = item.ID
		}
	}
	if retValue == "" {
		err = errors.New("Channel not found!")
	}
	return retValue, err
}
