package main

import (
	"fmt"
	"ghost/server/entrypoint"
	"log"
	"net/http"
	"os"
)

var storage_nodes_client = &http.Client{}

func main() {
	port := os.Args[1]
	log.Println("Starting http file sever")
	http.HandleFunc("/upload", entrypoint.HandleUpload)

	http.HandleFunc("/download", entrypoint.HandleDownload)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Println(err)
	}
}
