package main

import (
	"fmt"
	"go-storage/server/storage"
	"net/http"
	"os"
)

var nextNodeClient = &http.Client{}

func main() {
	port := os.Args[1]
	fmt.Println("Starting http file sever")

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		storage.HandleUploadStorage(w, r, nextNodeClient)
	})

	http.HandleFunc("/download", storage.HandleDownloadStorage)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		fmt.Println(err)
	}
}
