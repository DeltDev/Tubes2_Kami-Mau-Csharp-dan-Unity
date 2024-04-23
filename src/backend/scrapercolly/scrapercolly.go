package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

type void struct{}

var member void

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
	)

	c.OnHTML(".mw-body-content", func(e *colly.HTMLElement) {
		links := e.ChildAttrs("a", "href")
		for _, link := range links {

			if validLink(link) {
				fmt.Println(link)
			}
		}
	})

	c.Visit("https://en.wikipedia.org/Adolf_Hitler")

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
