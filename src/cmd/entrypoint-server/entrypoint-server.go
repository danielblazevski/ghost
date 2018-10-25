package main

import (
	"fmt"
	"net/http"
	"os"
	"server/entrypoint"
)

var storage_nodes_client = &http.Client{}

func main() {
	port := os.Args[1]
	fmt.Println("Starting http file sever")
	http.HandleFunc("/upload", entrypoint.HandleUpload)

	http.HandleFunc("/download", entrypoint.HandleDownload)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		fmt.Println(err)
	}
}
