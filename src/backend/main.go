package main

import (
    "fmt"
    "net/http"
	"strings"

	"golang.org/x/net/html"
)


func bfs(startPage string, endPage string) []string {

	path := [][]string{{startPage}}
	queue := []string{startPage}
	visited := make(map[string]bool)
	visited[startPage] = false
	if startPage == endPage {
		fmt.Println("Found the end page!")
		fmt.Println("Path: ", startPage)
		return []string{startPage}
	}
	for len(queue) > 0 {
		currentPage := queue[0]
		queue = queue[1:]
		if !visited[currentPage] {
			visited[currentPage] = true
			links := getLinks(currentPage)
			for _, link := range links {
				if !visited[link] {
					if link == endPage {
						fmt.Println("Found the end page!")
						fmt.Println("Path: ", currentPage, " -> ", link)
						for i:=0; i < len(path); i++{
							if path[i][len(path[i])-1] == currentPage {
								temp := make([]string, len(path[i]))
								copy(temp, path[i])
								temp = append(temp, link)
								path = append(path, temp)
								return path[len(path)-1]
							}
						}
					}

					for i:=0; i < len(path); i++{
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
		fmt.Println("Error: ", err)
		return nil
	}
	defer resp.Body.Close()
	links := []string{}
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
						if strings.HasPrefix(a.Val, "/wiki/") && !strings.Contains(a.Val, ":"){
						havekHref = true
						}
					}
					if a.Key == "title" {
						// tidak boleh ada titik dua
						str	:= a.Val
						if !strings.Contains(str, ":"){
							haveTitle = true
						}
					}
				}
				if havekHref && haveTitle {
					var pranala string = strings.TrimPrefix(t.Attr[0].Val, "/wiki/")
					links = append(links, pranala)
				}
			}
		}
	}
}

func main() {
    // Membuat server untuk frontend
    fs := http.FileServer(http.Dir("../frontend"))
    http.Handle("/", fs)

    // Proses mengambil data dari form
    http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
        // Mengecek data dari form
        err := r.ParseForm()
        if err != nil {
            http.Error(w, "Failed to parse form data", http.StatusBadRequest)
            return
        }

        // Mengambil value dan meletakkannya pada variable
        start := r.Form.Get("start")
        finish := r.Form.Get("finish")
		algorithm := r.Form.Get("algorithm")

		// Debug
		fmt.Printf("Start: %s, Finish: %s, Algorithm: %s\n", start, finish, algorithm)

        if algorithm == "TES" {
            start := start
            // daftar:= getLinks(start)
            // fmt.Println(daftar)
			// fmt.Fprintf(w, "%v", daftar)
			path := bfs(start, finish)
			fmt.Fprintf(w, "%v", path)


        }

		// Memproses data (ubah ini)
        w.WriteHeader(http.StatusOK)
        // w.Write([]byte("Data received successfully"))
    })

    // Menyalakan server
    fmt.Println("Server is running on port 8080")
    http.ListenAndServe(":8080", nil)
}