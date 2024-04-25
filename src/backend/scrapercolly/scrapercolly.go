package scrapercolly

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

func CollyGetLinks(url string) []string {
	url = "https://en.wikipedia.org/wiki/" + url
	c := colly.NewCollector()

	links := []string{}
	alreadyAdded := make(map[string]bool)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		if strings.HasPrefix(link, "/wiki/") && !strings.Contains(link, ":") && (link != "/wiki/Main_Page") {
			if !alreadyAdded[link] {
				links = append(links, strings.TrimPrefix(link, "/wiki/"))
				alreadyAdded[link] = true
			}
		}
	})

	err := c.Visit(url)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	return links
}
