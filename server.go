package main

import (
	"log"
	"net/http"
	"time"

	"./src/downloader"
	"./src/routes"
)

// sync files when the server starts
func init() {
	downloader.SyncFiles()
}

func main() {
	ticker := time.NewTicker(5 * time.Minute)
	go func(ticker *time.Ticker) {
		for {
			select {
			case <-ticker.C:
				downloader.SyncFiles()
			}
		}
	}(ticker)

	http.Handle("/", http.HandlerFunc(routes.Index))
	http.Handle("/download", http.HandlerFunc(routes.Downloader))
	log.Printf("Serving %s on HTTP port: %s\n", "/static", ":8001")
	log.Fatal(http.ListenAndServe(":8001", nil))
}
