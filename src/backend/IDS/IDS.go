package IDS

import(
	"fmt"
	"sync"
	"backend/scraper"
)
func IDS(startPage string, endPage string) []string {
	if startPage == endPage { //cek apakah awal dan akhirnya sama
		return []string{startPage}
	}
	//debug]]
	links := scraper.GetLinksArr(startPage)
	if len(links) == 0 { //handling error: halaman tidak ada di wikipedia
		return []string{}
	}
	// fmt.Println(links)
	path := []string{}

	for iteration := 0; iteration <= 8; iteration += 3 { //iterasi fragmentasi
		path = IDSFragment(startPage, endPage, iteration, iteration+2)
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

	links := scraper.GetLinksArr(src) //dapatkan semua link yang ada di halaman yang sedang dikunjungi
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

func IDSFragment(startPage string, endPage string, startIdx int, endIdx int) []string {
	var wg sync.WaitGroup
	wg.Add(3)

	ch := make(chan []string, 3)
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

	for i := 0; i < 3; i++ {
		path := <-ch

		if path != nil {
			if path[0] == startPage && path[len(path)-1] == endPage {
				return path
			}
		}
	}
	return nil //return path yang udah ketemu
}