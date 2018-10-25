package main

import (
	"fmt"
	"net/http"
	"os"
	"server/storage"
)

var nextNodeClient = &http.Client{}

func main() {
	port := os.Args[1]
	nextService := os.Args[2]
	nextPort := os.Args[3]

	fmt.Println("Starting http file sever")

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		storage.HandleUploadStorage(w, r, nextNodeClient, nextService, nextPort)
	})

	http.HandleFunc("/download", storage.HandleDownloadStorage)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		fmt.Println(err)
	}
}
