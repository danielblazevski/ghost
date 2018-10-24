package main

import (
	"fmt"
	"go-storage/server/gateway"
	"net/http"
	"os"
)

var storage_nodes_client = &http.Client{}

func main() {
	port := os.Args[1]
	fmt.Println("Starting http file sever")
	http.HandleFunc("/upload", gateway.HandleUpload)

	http.HandleFunc("/download", gateway.HandleDownload)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		fmt.Println(err)
	}
}
