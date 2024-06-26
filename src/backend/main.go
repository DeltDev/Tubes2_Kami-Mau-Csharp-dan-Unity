package main

import (
	"backend/BFS"
	"backend/IDS"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Response struct {
	Path     []string
	PathLink []string
	Degree   int
	Duration time.Duration
	File     int
}

func main() {
	// Membuat server untuk frontend
	// sekaligus inisialisasi awal empty array
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var paths []string

		data := Response{
			Path:     paths,
			PathLink: paths,
			Degree:   0,
			Duration: 0,
			File:     0,
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
		var nFileVisited int
		startTime := time.Now()
		if algorithm == "BFS" {
			pathLink, nFileVisited = BFS.BFS(start, finish)
		} else if algorithm == "IDS" {
			pathLink, nFileVisited = IDS.IDS(start, finish)

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
		fmt.Println("Jumlah File yang dikunjungi BFS: ", nFileVisited)
		// fmt.Println("Duration:", duration)

		// Passing ke HTML
		data := Response{
			Path:     path,
			PathLink: pathLink,
			Degree:   degree,
			Duration: duration,
			File:     nFileVisited,
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
