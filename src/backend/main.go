package main

import (
	"backend/scrapercolly"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

type Response struct {
	Path   []string
	PathLink []string
	Degree int
	Duration time.Duration
}

func bfs(startPage string, endPage string) []string {

	path := [][]string{{startPage}}
	queue := []string{startPage}
	visited := make(map[string]bool)
	visited[startPage] = false
	if startPage == endPage {
		// fmt.Println("Found the end page!")
		// fmt.Println("Path: ", startPage)
		return []string{startPage}
	}
	for len(queue) > 0 {
		currentPage := queue[0]
		queue = queue[1:]
		if !visited[currentPage] {
			visited[currentPage] = true
			links := scrapercolly.CollyGetLinks(currentPage)
			// fmt.Println("links: ", links)
			for _, link := range links {
				if !visited[link] {
					if link == endPage {
						// fmt.Println("Found the end page!")
						for i := 0; i < len(path); i++ {
							if path[i][len(path[i])-1] == currentPage {
								temp := make([]string, len(path[i]))
								copy(temp, path[i])
								temp = append(temp, link)
								path = append(path, temp)
								return path[len(path)-1]
							}
						}
					}

					for i := 0; i < len(path); i++ {
						if path[i][len(path[i])-1] == currentPage {
							// fmt.Println("currentPage: ", currentPage)
							temp := make([]string, len(path[i]))
							copy(temp, path[i])
							temp = append(temp, link)
							// fmt.Println("tempAkhir: ", temp)
							path = append(path, temp)
							break
						}
					}

					queue = append(queue, link)
				}
			}
		}
	}
	return []string{}
}

func getLinks(url string) []string {
	url = "https://en.wikipedia.org/wiki/" + url
	resp, err := http.Get(url)
	if err != nil {
		// fmt.Println("Error: ", err)
		return nil
	}
	defer resp.Body.Close()
	links := []string{}
	alreadyAdded := make(map[string]bool)
	z := html.NewTokenizer(resp.Body)
	for {
		//contoh tt adalah EndTag atau bisa jadi StartTag
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return links
		case html.StartTagToken:
			// contoh t adalah <a href="/wiki/Category:Go_(programming_language)" title="Category:Go (programming language)">
			// atau <ul> atau <li> dan lain lain
			t := z.Token()
			// contoh t.Attr adalah [{ href //en.m.wikipedia.org/w/index.php?title=Go_(programming_language)&mobileaction=toggle_view_mobile} { class noprint stopMobileRedirectToggle}]
			// atau
			// contoh t.Data adalah a  atau ul atau li atau div, dan lain lain
			// buat boolean checkHref
			havekHref := false
			haveTitle := false
			if t.Data == "a" {
				for _, a := range t.Attr {
					// fmt.Println(a)
					// contoh a itu adalah { title AngularJS} atau { href /wiki/AngularJS} atau { href https://developer.wikimedia.org} atau dll
					if a.Key == "href" {
						if strings.HasPrefix(a.Val, "/wiki/") && !strings.Contains(a.Val, ":") && (a.Val != "/wiki/Main_Page") {
							havekHref = true
						}
					}
					if a.Key == "title" {
						// tidak boleh ada titik dua
						str := a.Val
						if !strings.Contains(str, ":") {
							haveTitle = true
						}
					}
				}
				if havekHref && haveTitle {
					var pranala string = strings.TrimPrefix(t.Attr[0].Val, "/wiki/")
					if !alreadyAdded[pranala] {
						links = append(links, pranala)
						alreadyAdded[pranala] = true
					}
				}
			}
		}
	}
}

func main() {
	// Membuat server untuk frontend
	// sekaligus inisialisasi awal empty array
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var paths []string

		data := Response{
			Path: paths,
			PathLink: paths,
			Degree: 0,
			Duration: 0,
		}

		tmpl, err := template.ParseFiles("../frontend/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Proses mengambil data dari form
	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		// Mengecek data dari form
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		// Mengambil value dan meletakkannya pada variable
		start := url.QueryEscape(strings.ReplaceAll(r.Form.Get("start"), " ", "_"))
		finish := url.QueryEscape(strings.ReplaceAll(r.Form.Get("finish"), " ", "_"))
		algorithm := r.Form.Get("algorithm")

		// Debug
		fmt.Printf("Start: %s, Finish: %s, Algorithm: %s\n", start, finish, algorithm)

		// Mencari hasil
		var path []string
		var pathLink []string
		startTime := time.Now()
		if algorithm == "BFS" {
			pathLink = bfs(start, finish)
		} else if algorithm == "IDS" {
			pathLink = IDS(start, finish)
		}
		endTime := time.Now()

		// Judul
		for _, link := range pathLink {
			decodedLink, err := url.QueryUnescape(link)
			if err != nil {
				fmt.Println("Error decoding link:", err)
				return
			}
			path = append(path, decodedLink)
		}

		// Degree
		degree := len(pathLink) - 1
		
		// Duration
		duration := endTime.Sub(startTime)

		// Debug
		fmt.Println(path)
		// fmt.Println("Duration:", duration)

		// Passing ke HTML
		data := Response{
			Path: path,
			PathLink: pathLink,
			Degree: degree,
			Duration: duration,
		}

		tmpl, err := template.ParseFiles("../frontend/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Menyalakan server
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func IDS(startPage string, endPage string) []string {
	if startPage == endPage{ //cek apakah awal dan akhirnya sama
		return []string{startPage}
	}
	//debug]]
	links := scrapercolly.CollyGetLinks(startPage)
	if len(links) == 0 { //handling error: halaman tidak ada di wikipedia
		return []string{}
	}
	// fmt.Println(links)
	path := []string{}

	for iteration := 0; iteration<=8; iteration+=2{
		path = IDSFragment(startPage,endPage,iteration,iteration+1);
		if path != nil {
			break
		}
	}

	return path
}

func DLS(src string, target string, limit int, visited map[string]bool, stopExplore chan bool) ([]string, bool) {
	fmt.Println("Halaman yang dikunjungi sekarang: ", src, "Halaman tujuan: ", target, "Batas kedalaman iterasi: ", limit)
	visited[src] = true
	if src == target { //kalau halaman yang divisit sekarang sama dengan halaman yang dicari
		ret := []string{src} //masukin ke path
		return ret, true     //artinya sudah ketemu pathnya
	}

	if limit <= 0 { //kalau limitnya sudah habis DAN src dan targetnya beda
		return nil, false //tidak ketemu
	}

	links := scrapercolly.CollyGetLinks(src) //dapatkan semua link yang ada di halaman yang sedang dikunjungi
	for _, nextLink := range links {         //iterasi ke semua link yang ada di halaman yang sedang dikunjungi
		if visited[nextLink] {
			continue
		}
		visited[nextLink] = true
		select {
		case explored := <-stopExplore:
			if explored {
				return nil, true
			}
		default:
			subPath, found := DLS(nextLink, target, limit-1, visited, stopExplore) //kunjungi node selanjutnya dan kurangi limit dengan 1 dan dapatkan nilai subpath dan nilai sudah ketemu path atau belum
			if found {                                                             //kalau ketemu
				return append([]string{src}, subPath...), true //tambahkan nama halaman yang sedang dikunjungi sekarang ke subpath dan tandai pathnya ketemu
			}
		}

		delete(visited, nextLink)
	}
	return nil, false //tidak ketemu pathnya
}

func IDSFragment(startPage string, endPage string, startIdx int, endIdx int) ([]string){
	var wg sync.WaitGroup
	wg.Add(2)

	ch := make(chan []string, 2)
	stopExplore := make(chan bool, 1)
	for iteration := startIdx; iteration <= endIdx; iteration++ { //tambah kedalaman terus sampai ketemu pathnya
		go func(d int) {
			defer wg.Done()
			path, found := DLS(startPage, endPage, d, map[string]bool{}, stopExplore)
			if found {
				ch <- path
				stopExplore <- true
				return
			}
			ch <- nil
		}(iteration)
	}

	wg.Wait()

	for i := 0; i < 2; i++ {
		path := <-ch

		if path != nil {
			if path[0] == startPage && path[len(path)-1] == endPage {
				return path
			}
		}
	}
	return nil //return path yang udah ketemu
}