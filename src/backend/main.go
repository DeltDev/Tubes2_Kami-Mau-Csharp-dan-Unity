package main

import (
	"fmt"
	"net/http"
)

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

		// Debug
        fmt.Printf("Start: %s, Finish: %s\n", start, finish)

		// Memproses data (ubah ini)
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Data received successfully"))
    })

    // Menyalakan server
    fmt.Println("Server is running on port 8080")
    http.ListenAndServe(":8080", nil)
}
