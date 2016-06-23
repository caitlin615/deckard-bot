/*
Package write renders a PNG using the Handwriting.io API.

To use this plugin, you will need to create an account at Handwriting.io, create
a token, and then add the following environment variable:

 HANDWRITINGIO_API_URL=https://token:secret@api.handwriting.io
 S3_BUCKET=s3 bucket to store and link to rendered images
*/
package write

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/url"
	"regexp"
	"sort"
	"time"

	"github.com/handwritingio/deckard-bot/config"
	"github.com/handwritingio/deckard-bot/log"
	"github.com/handwritingio/deckard-bot/message"
	"github.com/handwritingio/go-client/handwritingio"

	"github.com/aws/aws-sdk-go/aws"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Plugin ...
type Plugin struct {
	HandwritingioAPIURL string
}

var (
	writeBase      = regexp.MustCompile(`(?i)^!write`)
	writeCmd       = regexp.MustCompile(`(?i)^!write (\w{12}|\s?)(.+)$`)
	handwritingIDs []string
	client         *handwritingio.Client
)

// Usage prints detailed usage instructions for the plugin
func (p Plugin) Usage() string {
	return "*Usage:* `!write <text>` to handwrite some text"
}

// Command returns a list of commands the plugin provides
func (p Plugin) Command() []string {
	return []string{"!write"}
}

// OnInit returns an error if the plugin could not be started
func (p Plugin) OnInit() error {

	s := p.HandwritingioAPIURL
	if s == "" {
		return errors.New("HandwritingAPIURL must be set to use this plugin!")
	}

	u, err := url.Parse(s)
	if err != nil {
		return err
	}

	client, err = handwritingio.NewClientURL(u)
	if err != nil {
		return err
	}

	ids, err := listHandwritings()
	if err != nil {
		return err
	}
	handwritingIDs = ids
	return nil
}

// Name is the name of the plugin
func (p Plugin) Name() string {
	return "Write"
}

// HandleMessage is responsible for handling the incoming message
// and returning a response based on the message provides
func (p Plugin) HandleMessage(in message.Basic) (out message.Basic) {
	chunks := writeCmd.FindStringSubmatch(in.Text)
	if len(chunks) != 3 { // invalid command
		out.Text = p.Usage()
		return
	}
	handwritingID := chunks[1]
	text := chunks[2]

	url, err := write(text, handwritingID)
	if err != nil {
		out.Text = "error writing your message: " + err.Error()
		return
	}
	out.Text = url

	return
}

// Regexp returns the regexp of a message that should be handled by this plugin
func (p Plugin) Regexp() *regexp.Regexp {
	return writeBase
}

// listHandwritings gets all handwritings from the API
func listHandwritings() (ids []string, err error) {
	params := handwritingio.DefaultHandwritingListParams
	params.Limit = 100
	params.OrderBy = "title"
	params.OrderDir = "asc"
	handwritings, err := client.ListHandwritings(params)
	if err != nil {
		return
	}

	for len(handwritings) > 0 {
		for _, hw := range handwritings {
			ids = append(ids, hw.ID)
		}

		params.Offset += params.Limit
		handwritings, err = client.ListHandwritings(params)
		if err != nil {
			return
		}

	}
	sort.Strings(ids)
	log.Debugf("Loaded %d handwritings", len(ids))
	return
}

// render will render a png based on the supplied text using
// the supplied handwriting ID.
func render(text string, handwritingID string) ([]byte, error) {
	params := handwritingio.DefaultRenderParamsPNG
	params.Text = text
	params.HandwritingID = handwritingID
	params.Height = "auto"

	r, err := client.RenderPNG(params)
	if err != nil {
		return []byte{}, err
	}
	defer r.Close()

	return ioutil.ReadAll(r)
}

// upload uploads the rendered image to the S3 bucket and returns
// the URL where its uploaded to.
func upload(img []byte) (string, error) {
	if config.S3Bucket == "" {
		return "", fmt.Errorf("S3_BUCKET variable not configured")
	}
	reader := bytes.NewReader(img)
	timestamp := time.Now().Format("2006-01-02_150405.000")

	uploader := s3manager.NewUploader(awssession.New(&aws.Config{Region: aws.String(config.AWSRegion)}))
	out, err := uploader.Upload(
		&s3manager.UploadInput{
			Bucket:      aws.String(config.S3Bucket),
			ACL:         aws.String("public-read"),
			ContentType: aws.String("image/png"),
			Key:         aws.String(timestamp + ".png"),
			Body:        reader,
		})

	return out.Location, err
}

// write returns the S3 url of the rendered text
func write(text, handwritingID string) (url string, err error) {
	log.Println("Writing something")
	if len(handwritingIDs) < 1 {
		return "", fmt.Errorf("no handwritings available")
	}
	handwritingIndex := sort.SearchStrings(handwritingIDs, handwritingID)
	if handwritingID == "" || handwritingIndex == len(handwritingIDs) {
		handwritingID = handwritingIDs[rand.Intn(len(handwritingIDs))]
	}
	img, err := render(text, handwritingID)
	if err != nil {
		return "", fmt.Errorf("rendering error: %s", err)
	}
	log.Printf("rendered %d bytes\n", len(img))

	return upload(img)
}
