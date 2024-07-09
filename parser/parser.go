package parser

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

var (
	regExLinks    = regexp.MustCompile(`<a.*?href="/(.*?)[#"]`)
	regExTitle    = regexp.MustCompile(`(?s)<title.*?>(.*?)</title>`)
	regExDesc     = regexp.MustCompile("(?s)<meta.*?name=\"description\".*?content=\"(.*?)\"")
	regExKeywords = regexp.MustCompile("(?s)<meta.*?name=\"keywords\".*?content=\"(.*?)\"")
)

// PagesData represents a slice of PageData.
type PagesData []PageData

// PageData represents a data from HTML page.
type PageData struct {
	URL        string `json:"path"`
	StatusCode int    `json:"status code"`
	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Keywords   string `json:"keywords"`
}

// Parser represents a parser for the page.
type Parser struct {
	Client http.Client
}

// New creates a new Parser.
func New() Parser {
	return Parser{
		//Client: client,
	}
}

// ParseResponse parses the URL and returns the data from the page.
func (p *Parser) ParseResponse(resp *http.Response) (PageData, []string, error) {
	if resp.StatusCode != http.StatusOK {
		return PageData{
			URL:        resp.Request.URL.String(),
			StatusCode: resp.StatusCode,
		}, nil, fmt.Errorf("returned status: %s, url: %#v", resp.Status, resp.Request.URL.String())
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return PageData{
			URL:        resp.Request.URL.String(),
			StatusCode: resp.StatusCode,
		}, nil, fmt.Errorf("can't read response body, url: %#v. Error: %s", resp.Request.URL.String(), err)
	}
	defer resp.Body.Close()

	contentString := string(content)

	links := p.getLinks(contentString)
	title := p.getTitle(contentString)
	desc := p.getDescription(contentString)
	keywords := p.getKeywords(contentString)

	return PageData{
		URL:        resp.Request.URL.String(),
		StatusCode: resp.StatusCode,
		Title:      title,
		Desc:       desc,
		Keywords:   keywords,
	}, p.unique(links), nil
}

func (p *Parser) getLinks(content string) []string {
	matches := regExLinks.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return nil
	}

	var links []string
	for _, m := range matches {
		links = append(links, m[1])
	}

	return links
}

func (p *Parser) getTitle(content string) string {
	matches := regExTitle.FindStringSubmatch(content)
	if len(matches) == 0 {
		return ""
	}

	return strings.TrimSpace(matches[1])
}

func (p *Parser) getDescription(content string) string {
	matches := regExDesc.FindStringSubmatch(content)
	if len(matches) == 0 {
		return ""
	}

	return strings.TrimSpace(matches[1])
}

func (p *Parser) getKeywords(content string) string {
	matches := regExKeywords.FindStringSubmatch(content)
	if len(matches) == 0 {
		return ""
	}

	return strings.TrimSpace(matches[1])
}

func (p *Parser) unique(intSlice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
