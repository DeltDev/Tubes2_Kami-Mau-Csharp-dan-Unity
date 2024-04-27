package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

func GetLinksArr(url string) []string {
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
func removeRedundantMaphashtag(data map[string][]string) map[string][]string {
	result := make(map[string][]string)

	for key, arr := range data {
		result[key] = removeRedundanthashtag(arr)
	}

	return result
}

func removeRedundanthashtag(array []string) []string {
	result := make([]string, 0)
	seen := make(map[string]int)

	for _, word := range array {
		parts := strings.Split(word, "#")
		baseWord := parts[0]

		if idx, ok := seen[baseWord]; ok {
			if len(word) < len(result[idx]) {
				result[idx] = word
			}
		} else {
			seen[baseWord] = len(result)
			result = append(result, word)
		}
	}

	return result
}

func removeRedundantMap(data map[string][]string) map[string][]string {
	result := make(map[string][]string)

	for key, arr := range data {
		result[key] = removeRedundant(arr)
	}

	return result
}

func removeRedundant(array []string) []string {
	result := make([]string, 0)
	seen := make(map[string]int)

	for _, word := range array {
		parts := strings.Split(word, "_")
		baseWord := parts[0]

		if idx, ok := seen[baseWord]; !ok || len(word) < len(array[idx]) {
			seen[baseWord] = len(result)
			result = append(result, word)
		}
	}

	return result
}
