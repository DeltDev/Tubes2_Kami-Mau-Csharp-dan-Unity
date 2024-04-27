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
	for iteration := 1; iteration <= 8; iteration += 2 { //iterasi fragmentasi IDS multithreading (fragmentasi: IDS dipecah menjadi asumsi 1-8 degree dengan membangkitkan 2 buah thread secara bersamaan)
		fmt.Println("Fragmen dari: ",iteration," sampai: ",iteration+1)
		path,temp= IDSFragment(startPage, endPage, iteration, iteration+1) //cari path dari kedalaman iteration sampai iteration+1 secara bersamaan dengan IDS terfragmentasi
		total += temp //jumlahkan link yang dikunjungi
		if path != nil { //kalau sudah ketemu path dari fragmen, jangan dilanjutkan IDSnya
			break
		}
	}

	return path,total //kembalikan path dan total link yang dikunjungi
}

func DLS(src string, target string, limit int, visited map[string]bool, stopExplore chan bool, visCount* int) ([]string, bool) {
	visited[src] = true //tandai sudah pernah divisit
	*visCount++ //tambah jumlah halaman yang dikunjungi dengan 1 (halaman saat ini)
	if src == target { //kalau halaman yang divisit sekarang sama dengan halaman yang dicari
		ret := []string{src} //masukin ke path
		return ret, true    //artinya sudah ketemu pathnya
	}

	if limit <= 0 { //kalau limitnya sudah habis DAN src dan targetnya beda
		return nil, false //tidak ketemu
	}

	links := scraper.GetLinksArr(src) //dapatkan semua link yang ada di halaman yang sedang dikunjungi
	links = scraper.SortStringsBySim(target,links) //urutkan semua link berdasarkan kemiripan dengan target
	fmt.Println(target) //debug
	fmt.Println(links)
	for _, nextLink := range links {         //iterasi ke semua link yang ada di halaman yang sedang dikunjungi
		if visited[nextLink] { //lewati link yang sudah dikunjungi
			continue
		}
		visited[nextLink] = true //tandai link selanjutnya dengan true
		select {
		case explored := <-stopExplore: //jika ada sinyal untuk menghentikan algoritma IDS
			if explored { 
				return nil, true //hentikan semua algoritma IDS
			}
		default:
			subPath, found := DLS(nextLink, target, limit-1, visited, stopExplore, visCount) //kunjungi node selanjutnya dan kurangi limit dengan 1 dan dapatkan nilai subpath dan nilai sudah ketemu path atau belum
			if found { //kalau ketemu
				return append([]string{src}, subPath...), true //tambahkan nama halaman yang sedang dikunjungi sekarang ke subpath dan tandai pathnya ketemu
			}
		}

		delete(visited, nextLink) //Backtrack visited
	}
	return nil, false //tidak ketemu pathnya
}

func IDSFragment(startPage string, endPage string, startIdx int, endIdx int) ([]string,int) { //thread IDS terfragmentasi
	var wg sync.WaitGroup
	wg.Add(2)

	ch := make(chan []string, 2)
	stopExplore := make(chan bool, 1)
	total:=0
	
	for iteration := startIdx; iteration <= endIdx; iteration++ { //tambah kedalaman terus dari startIDX sampai endIdx sampai ketemu pathnya
		initVal :=0
		go func(d int) { //buat goroutine
			defer wg.Done()
			path, found:= DLS(startPage, endPage, d, map[string]bool{}, stopExplore, &initVal) //dapatkan path dan apakah berhasil ditemukan
			if found { // jika ditemukan
				ch <- path //masukkan path ke channel
				stopExplore <- true //hentikan semua goroutine DLS
				total += initVal //jumlahkan banyak link yang divisit
				return
			}
			ch <- nil
		}(iteration)
	}

	wg.Wait()

	for i := 0; i < 2; i++ {
		path := <-ch

		if path != nil { //jika path tidak kosong
			if path[0] == startPage && path[len(path)-1] == endPage { //periksa jika link pertama adalah startPage dan link terakhir adalah endPage
				return path, total //return path yang udah ketemu dan banyak link total yang dikunjungi
			}
		}
	}
	return nil,total //return path yang tidak ditemukan
}