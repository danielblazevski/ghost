package main

import (
	"fmt"
	"ghost/server/storage"
	"log"
	"net/http"
	"os"
)

var nextNodeClient = &http.Client{}

func main() {
	port := os.Args[1]
	nextService := os.Args[2]
	nextPort := os.Args[3]

	log.Println("Starting http file sever")

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		storage.HandleUploadStorage(w, r, nextNodeClient, nextService, nextPort)
	})

	http.HandleFunc("/download", storage.HandleDownloadStorage)
	http.HandleFunc("/status-check", storage.HandleStatusCheck)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Println(err)
	}
}
