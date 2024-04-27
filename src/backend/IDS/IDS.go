package IDS

import (
	"backend/scraper"
	"fmt"
	"sync"
)
func IDS(startPage string, endPage string) ([]string,int) {
	if startPage == endPage { //cek apakah awal dan akhirnya sama
		return []string{startPage},1
	}
	links := scraper.GetLinksArr(startPage)
	if len(links) == 0 { //handling error: halaman tidak ada di wikipedia
		return []string{},0
	}
	path := []string{}
	temp:= 0
	total := 0
	for iteration := 1; iteration <= 8; iteration += 2 { //iterasi fragmentasi IDS multithreading (fragmentasi: IDS dipecah menjadi asumsi 1-9 degree)
		fmt.Println("Fragmen dari: ",iteration," sampai: ",iteration+1)
		path,temp= IDSFragment(startPage, endPage, iteration, iteration+1)
		total += temp
		if path != nil { //kalau sudah ketemu path dari fragmen, jangan dilanjutkan IDSnya
			break
		}
	}

	return path,total
}

func DLS(src string, target string, limit int, visited map[string]bool, stopExplore chan bool, visCount* int) ([]string, bool) {
	visited[src] = true //tandai sudah pernah divisit
	*visCount++
	if src == target { //kalau halaman yang divisit sekarang sama dengan halaman yang dicari
		ret := []string{src} //masukin ke path
		return ret, true    //artinya sudah ketemu pathnya
	}

	if limit <= 0 { //kalau limitnya sudah habis DAN src dan targetnya beda
		return nil, false //tidak ketemu
	}

	links := scraper.GetLinksArr(src) //dapatkan semua link yang ada di halaman yang sedang dikunjungi
	// links = scraper.RemoveRedundant(links) //hilangkan semua link yang redundant
	// links = scraper.RemoveRedundanthashtag(links) //hilangkan semua link yang mengandung # (karena hanya merupakan redirect ke halaman yang sama)
	links = scraper.SortStringsBySim(target,links) //urutkan semua link berdasarkan kemiripan dengan target
	fmt.Println(target)
	fmt.Println(links)
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
			subPath, found := DLS(nextLink, target, limit-1, visited, stopExplore, visCount) //kunjungi node selanjutnya dan kurangi limit dengan 1 dan dapatkan nilai subpath dan nilai sudah ketemu path atau belum
			if found {                                                             //kalau ketemu
				return append([]string{src}, subPath...), true //tambahkan nama halaman yang sedang dikunjungi sekarang ke subpath dan tandai pathnya ketemu
			}
		}

		delete(visited, nextLink)
	}
	return nil, false //tidak ketemu pathnya
}

func IDSFragment(startPage string, endPage string, startIdx int, endIdx int) ([]string,int) {
	var wg sync.WaitGroup
	wg.Add(2)

	ch := make(chan []string, 2)
	stopExplore := make(chan bool, 1)
	total:=0
	
	for iteration := startIdx; iteration <= endIdx; iteration++ { //tambah kedalaman terus sampai ketemu pathnya
		initVal :=0
		go func(d int) {
			defer wg.Done()
			path, found:= DLS(startPage, endPage, d, map[string]bool{}, stopExplore, &initVal)
			if found {
				ch <- path
				stopExplore <- true
				total += initVal
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
				return path, total
			}
		}
	}
	return nil,total //return path yang udah ketemu
}