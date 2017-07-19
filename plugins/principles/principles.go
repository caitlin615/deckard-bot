// Package principles searches the Handwriting.io team's Engineering principles
// document and returns the best fit to the search query
package principles

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/handwritingio/deckard-bot/github"
	"github.com/handwritingio/deckard-bot/log"
	"github.com/handwritingio/deckard-bot/message"

	"github.com/renstrom/fuzzysearch/fuzzy"
)

// Plugin holds the list of principles
type Plugin struct {
	List []*Principle
}

// Principle contains the title and description of an engineering principle
type Principle struct {
	Title       string
	Description string
	Number      int
}

var (
	// rePrinciple is the regexp variables for logic in HandleMessage
	rePrinciple        = regexp.MustCompile(`(?i)^!principle`)
	rePrincipleAll     = regexp.MustCompile(`(?i)^!principles?$`)
	rePrincipleNum     = regexp.MustCompile(`(?i)^!principles?\s+(\d+)`)
	rePrincipleKeyword = regexp.MustCompile(`(?i)^!principles?\s+(.+)`)
	rePrincipleFormat  = regexp.MustCompile(`(?i)^(\d+\.)\s+\*\*(\w.+)\*\*\s+(.+)`)

	// Default values for weighting the fuzzy search to prefer the title
	// The higher the value, the less important it will be in returning a match
	titleFactor = 1
	descFactor  = 3

	// Specifies the location in github where we store our Engineering Principles
	// https://github.com/handwritingio/principles/blob/master/EngineeringPrinciples.md
	principleOrg      = "handwritingio"
	principleRepo     = "principles"
	principleFilename = "EngineeringPrinciples.md"
)

// Usage returns the Plugin's usage
func (p *Plugin) Usage() string {
	return "`!principle` will give you some Engineering principles\n" +
		"`!principles` will list all\n" +
		"`!principle {digit} will list the specified principle" +
		"`!principle {keyword} will search for a principle based on the keyword"
}

// Command lists the base commands to use the plugin
func (p *Plugin) Command() []string {
	return []string{"!principle"}
}

// OnInit handles all actions that should occur when the plugin starts
func (p *Plugin) OnInit() error {
	// Get Engineering principles from Github
	data, err := getPrinciples()
	if err != nil {
		return fmt.Errorf("Error getting principles: %s", err.Error())
	}
	// Build PrincipleList data into struct
	principleList := buildPrinciples(data)
	if len(principleList) == 0 {
		return fmt.Errorf("Error building principles: no principles found")
	}
	// Set principles to Plugin struct
	p.List = principleList
	return nil
}

// Name returns the name of the plugin
func (p *Plugin) Name() string {
	return "Engineering Principles"
}

// Regexp returns the regexp of a message that should be handled by this plugin
func (p Plugin) Regexp() *regexp.Regexp {
	return rePrinciple
}

// HandleMessage takes a message.Basic in, checks to see if it needs to do anything
// with it (using regexp), and returns a message.Basic with a response
func (p *Plugin) HandleMessage(in message.Basic) (out message.Basic) {

	if len(p.List) == 0 {
		out.Text = "Sorry, there are no principles loaded at this time."
		return
	}
	switch {
	// Matches the command to list all (`!principle` or `!principles`)
	case rePrincipleAll.MatchString(in.Text):

		var s []string
		for i, principle := range p.List {
			num := i + 1
			s = append(s, fmt.Sprintf("%d. *%s*: %s", num, principle.Title, principle.Description))
		}
		out.Text = strings.Join(s, "\n")

	// Matches the command to list a numbered principle (`!principle 5`)
	case rePrincipleNum.MatchString(in.Text):
		chunks := rePrincipleNum.FindStringSubmatch(in.Text)
		if len(chunks) != 2 { // invalid command
			out.Text = p.Usage()
			return
		}
		num, _ := strconv.Atoi(chunks[1])
		if num >= len(p.List) {
			out.Text = "Sorry, the principle you requested does not exist"
			return
		}
		principle := p.List[num-1]
		out.Text = fmt.Sprintf("%d. *%s*: %s", principle.Number, principle.Title, principle.Description)

	// Matches the command to list a principle based on a keyword (`!principle code`)
	case rePrincipleKeyword.MatchString(in.Text):
		chunks := rePrincipleKeyword.FindStringSubmatch(in.Text)
		if len(chunks) != 2 { // invalid command
			out.Text = p.Usage()
			return
		}
		keyword := chunks[1]
		matchPrinciple := p.fuzzySearch(keyword)
		if matchPrinciple == nil {
			out.Text = fmt.Sprintf("Sorry, no principles match keyword `%s`", keyword)
			return
		}
		out.Text = fmt.Sprintf("%d. *%s*: %s", matchPrinciple.Number, matchPrinciple.Title, matchPrinciple.Description)
	default:
		out.Text = p.Usage()
	}

	return out
}

// fuzzySearch takes a keyword as a string, completes a search through titles and descriptions
// and returns the highest ranked result as a Principle. The rank is based on Levenshtein distance.
func (p *Plugin) fuzzySearch(keyword string) *Principle {

	log.Debugf("Starting searching for keyword: `%s`", keyword)
	ranks := make(map[int]*Principle)
	for _, v := range p.List {

		titleRank := fuzzy.RankMatchFold(keyword, v.Title)
		descRank := fuzzy.RankMatchFold(keyword, v.Description)

		// Neither return a match, so skip it
		if titleRank == -1 && descRank == -1 {
			continue
		}
		// Add principle with a weighted title rank as long as the original rank isn't -1 (no match)
		if titleRank != -1 {
			ranks[titleRank*titleFactor] = v
		}
		// Add principle with a weighted description rank as long as the original rank isn't -1 (no match)
		if descRank != -1 {
			ranks[descRank*descFactor] = v
		}
	}
	// Sort ranked search results
	var keys []int
	for k := range ranks {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	// Returns nil if no matches were found
	if len(keys) == 0 {
		log.Debugf("No matches were found for keyword: `%s`", keyword)
		return nil
	}
	// Get the top rank and the top ranked principle
	topRank := keys[0]
	topPrinciple := ranks[topRank]
	log.Debugf("Top ranked principle: Rank %d, %v", topRank, topPrinciple)
	log.Debugf("Done with search for keyword `%s`", keyword)
	return topPrinciple
}

// getPrinciples returns the data from from the EngineeringPrinciples.md file in Github
func getPrinciples() ([]byte, error) {
	githubClient := github.NewClient("")
	contents, _, err := githubClient.GetFile(principleOrg, principleRepo, principleFilename)
	if err != nil {
		log.Warnf("Error encountered getting file contents: %s", err.Error())
		return nil, err
	}
	return contents, err
}

// buildPrinciples takes the principle data in bytes and returns a list of Principle
func buildPrinciples(data []byte) (list []*Principle) {
	// Split principles by double newline (seperator between principles)
	principles := strings.Split(string(data), "\n\n")

	// Loop through all principles, check that it matches regex, add to Principle list
	for _, pr := range principles {
		prString := strings.Replace(pr, "\n", " ", -1)
		if !rePrincipleFormat.MatchString(prString) {
			continue
		}
		prChunks := rePrincipleFormat.FindStringSubmatch(prString)

		var tmp Principle
		tmp.Title = prChunks[2]
		tmp.Description = prChunks[3]
		list = append(list, &tmp)
	}
	for i, principle := range list {
		principle.Number = i + 1
	}

	return
}
