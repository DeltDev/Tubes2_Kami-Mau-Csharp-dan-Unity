package scraper

import (
	"fmt"
	"strings"
	"sync"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	
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

func RemoveRedundantMaphashtag(data map[string][]string) map[string][]string {
	result := make(map[string][]string)

	for key, arr := range data {
		result[key] = RemoveRedundanthashtag(arr)
	}

	return result
}

func RemoveRedundanthashtag(array []string) []string {
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

func RemoveRedundantMap(data map[string][]string) map[string][]string {
	result := make(map[string][]string)

	for key, arr := range data {
		result[key] = RemoveRedundant(arr)
	}

	return result
}

func RemoveRedundant(array []string) []string {
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

type SafeLinksMap struct {
	sync.Map
	mux sync.Mutex // Mutex for SafeLinksMap
}

func (s *SafeLinksMap) StoreLinks(key string, value []string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.Store(key, value)
}

func (s *SafeLinksMap) LoadLinks(key string) ([]string, bool) {
	value, ok := s.Load(key)
	if !ok {
		return nil, false
	}
	return value.([]string), true
}

func GetLinksMap(juduls []string) map[string][]string {
	results := SafeLinksMap{}

	c := colly.NewCollector(
		colly.MaxDepth(1),
	)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 4000})

	isExist := SafeLinksMap{}

	q, _ := queue.New(
		22,
		&queue.InMemoryQueueStorage{MaxSize: 1000000},
	)

	for _, judul := range juduls {
		q.AddURL("https://en.wikipedia.org/wiki/" + judul)
		results.StoreLinks(judul, make([]string, 0)) // Initialize with empty slice
	}

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.HasPrefix(link, "/wiki/") && !strings.Contains(link, ":") && link != "/wiki/Main_Page" {
			if _, ok := isExist.Load(link); !ok {
				isExist.Store(link, true)
				link := strings.TrimPrefix(link, "/wiki/")
				judul := e.Request.URL.String()
				judul = strings.TrimPrefix(judul, "https://en.wikipedia.org/wiki/")
				value, ok := results.LoadLinks(judul)
				if ok {
					links := value
					if !Contains(links, link) { // Check if link already exists
						links = append(links, link)
						results.StoreLinks(judul, links)
					}
				}
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("visiting", r.URL)
	})

	err := q.Run(c)

	if err != nil {
		return nil
	}

	c.Wait()

	hasil := make(map[string][]string)
	results.Range(func(key, value interface{}) bool {
		hasil[key.(string)] = append([]string{}, value.([]string)...)
		return true
	})

	return hasil
}

// Helper function to check if a slice contains a string
func Contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}