package main

import (
	"fmt"
	"ghost/pkg/server/storage"
	"log"
	"net/http"
	"os"
)

var nextNodeClient = &http.Client{}

func main() {
	port := os.Args[1]
	nextService := os.Args[2]
	nextPort := os.Args[3]

	baseLocation = "/ghost/files"

	log.Println("Starting http file sever")

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		storage.HandleUploadStorage(w, r, nextNodeClient, nextService, nextPort, baseLocation)
	})

	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		storage.HandleDownloadStorage(w, r, baseLocation)
	})

	http.HandleFunc("/status-check", storage.HandleStatusCheck)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Println(err)
	}
}
