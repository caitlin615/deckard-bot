// Package github is a wrapper around the go Github client and API
package github

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/handwritingio/deckard-bot/log"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Client is a wrapper for the github Client
type Client struct {
	client *github.Client
}

const archiveFormat = github.Tarball

// NewClient creates a new Client including authentication
func NewClient(apiKey string) *Client {
	if apiKey == "" {
		// return a non-authenticated client if an API key isn't set,
		// (so client can still access public resources)
		return &Client{client: github.NewClient(nil)}
	}
	// return an authenticated client
	// https://github.com/google/go-github#authentication
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: apiKey},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	return &Client{
		client: github.NewClient(tc),
	}
}

// GetFile returns the contents of a file and the download URL of the file
// from a file within a github repository. A repository and path to a file must be supplied.
func (c *Client) GetFile(org, repo, path string) ([]byte, string, error) {
	opt := &github.RepositoryContentGetOptions{}
	content, _, resp, err := c.client.Repositories.GetContents(org, repo, path, opt)
	if resp.StatusCode != 200 {
		return nil, "", errors.New("Bad response from Github: " + resp.Status)
	}
	if err != nil {
		return nil, "", err
	}
	decoded, err := content.Decode()
	if err != nil {
		return nil, "", err
	}
	return decoded, github.Stringify(content.DownloadURL), nil
}

// CheckGithubRateLimit returns the API Rate limit to the debug console
// https://github.com/google/go-github/blob/master/examples/repos/main.go
func (c *Client) CheckGithubRateLimit() {
	rate, _, err := c.client.RateLimit()
	if err != nil {
		log.Debugf("Error fetching Github rate limit: %#v\n", err)
	} else {
		log.Debugf("Github API Rate Limit: %#v\n", rate)
	}
}

// checkGithubRepo takes a repo as a string you'd like to check
// and confirms whether or not the repo exists and the Client has access to it
func (c *Client) checkGithubRepo(org, repo string) bool {
	repos, _, _ := c.client.Repositories.ListByOrg(org, nil)
	for _, r := range repos {
		if *r.Name == repo {
			return true
		}
	}
	return false
}

// GetArchive returns an Archive based on the repo and branch supplied
func (c *Client) GetArchive(org, repo, branch string) (*url.URL, string, error) {
	err := c.checkRepoAndBranch(org, repo, branch)
	if err != nil {
		return nil, "", err
	}
	archiveURL, commit, err := c.getArchive(org, repo, branch)
	return archiveURL, commit, nil
}

func (c *Client) getArchive(org, repo, branch string) (*url.URL, string, error) {
	opts := github.RepositoryContentGetOptions{
		Ref: branch,
	}
	archiveURL, _, err := c.client.Repositories.GetArchiveLink(org, repo, archiveFormat, &opts)
	if err != nil {
		log.Errorf("Could not get archive URL: %s", err.Error())
		return nil, "", err
	}
	b, _, err := c.client.Repositories.GetBranch(org, repo, branch)
	if err != nil {
		return nil, "", err
	}
	commit := *b.Commit.SHA
	return archiveURL, commit, nil
}

// CheckBranch checks if the repo supplied exists and the branch exists for the
// supplied repo. Returns a boolean
func (c *Client) checkRepoAndBranch(org, repo, branch string) error {
	if !c.checkGithubRepo(org, repo) {
		return fmt.Errorf("Github repo not found: %s", repo)
	}
	branches, _, err := c.client.Repositories.ListBranches(org, repo, &github.ListOptions{})
	if err != nil {
		return fmt.Errorf("Could not fetch branches for %s: %s", repo, err.Error())
	}
	for _, b := range branches {
		if branch == *b.Name {
			return nil
		}
	}
	return fmt.Errorf("No branch named %s found in repo %s", branch, repo)
}

// GetGithubUsers returns the usernames for all users in the github organization
// This can then be used in the assignee section of !git issue. This is useful if you don't
// know the github username of the person you'd like to assign the issue to.
func (c *Client) GetGithubUsers(org string) (out string) {
	// Get Org members
	members, resp, err := c.client.Organizations.ListMembers(org, nil)
	if err != nil {
		out = "Error accessing Org membership"
	}
	log.Debug("Github user status: " + resp.Status)
	if resp.StatusCode != 200 {
		out = "Bad responds from Github: " + resp.Status
	}
	s := []string{"*Here's a list of all " + org + " Github usernames:*"}
	for _, r := range members {
		githubUsername := github.Stringify(r.Login)
		log.Debug("Github Username: " + githubUsername)
		s = append(s, githubUsername)
	}
	out = strings.Join(s, "\n")
	return
}

// CreateGithubIssue creates issues in github for the supplied repo
func (c *Client) CreateGithubIssue(org, repo, issue string) (out string) {

	// Check if repo exists
	if !c.checkGithubRepo(org, repo) {
		out = "PANIC: `" + repo + "` Repository Does Not Exist"
		return
	}

	// Creates issueRequest message based on supplied issue string
	issueMsg := github.IssueRequest{
		Title: github.String(issue),
		Body:  github.String("Issue created by the Deckard Chatbot Plugin"),
	}
	// Create issue
	i, resp, err := c.client.Issues.Create(org, repo, &issueMsg)
	if err != nil {
		out = fmt.Sprintf("Error occurred when creating issue: %s", err.Error())
		return
	}
	// Check returned status code
	issueStatusCode := resp.StatusCode
	if issueStatusCode != 201 {
		out = "PANIC: Issue was not created!\nResponse code: " + resp.Status
		return
	}
	issueNumber := *i.Number
	issueURL := *i.HTMLURL
	log.Debugf("Issue URL: %s", issueURL)
	log.Debugf("Issue number: %d", issueNumber)
	log.Debugf("Create issue status code: %d", issueStatusCode)

	out = fmt.Sprintf("*Issue # %d has been created successfully*\n%s", issueNumber, issueURL)

	return
}

// Octocat is a wrapper around github Client octocat
// prints an ASCII octocat
func (c *Client) Octocat(message string) string {
	octocat, _, _ := c.client.Octocat(message)
	return octocat
}
