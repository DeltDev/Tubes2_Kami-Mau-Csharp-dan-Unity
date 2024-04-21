package scrapercolly

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

type void struct{}

var member void

func ScrapeWikiLinks(initLink string) map[string]void {
	linkSet := make(map[string]void)
	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
	)

	c.OnHTML(".mw-body-content", func(e *colly.HTMLElement) {
		links := e.ChildAttrs("a", "href")
		for _, link := range links {

			if validLink(link) {
				linkSet[link] = member
			}
		}
	})

	c.Visit("https://en.wikipedia.org" + initLink)
	for k := range linkSet {
		fmt.Println(k)
	}

	return linkSet
}

func validLink(link string) bool {
	if !strings.HasPrefix(link, "/wiki/") {
		return false
	}
	if strings.Contains(link[len("/wiki/"):], ":") {
		return false
	}
	return true
}
