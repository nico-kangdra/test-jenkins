package main

import (
    "fmt"
    "net/http"
    "os"
)

// Handler fungsi untuk merespons permintaan web
func handler(w http.ResponseWriter, r *http.Request) {
    hostname, err := os.Hostname()
    if err != nil {
        hostname = "unknown"
    }
    fmt.Fprintf(w, "Hello, World! from %s", hostname)
}

func main() {
    // Menghubungkan rute "/" dengan fungsi handler
    http.HandleFunc("/", handler)

    fmt.Println("Server berjalan di http://0.0.0.0:8080")
    // Memulai web server pada port 8080
    http.ListenAndServe(":8080", nil)
}
