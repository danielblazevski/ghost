package storage

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

// TO-DO: grab latst version!
func HandleDownloadStorage(writer http.ResponseWriter, request *http.Request) {
	//First of check if Get is set in the URL
	filename := request.URL.Query().Get("file")
	if filename == "" {
		//Get not set, send a 400 bad request
		http.Error(writer, "Get 'file' not specified in url.", 400)
		return
	}
	fmt.Println("Client requests: " + filename)

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		//File not found, send 404
		http.Error(writer, "File not found.", 404)
		return
	}

	fileHeader := make([]byte, 512)
	file.Read(fileHeader)
	fileContentType := http.DetectContentType(fileHeader)

	fileStat, _ := file.Stat()
	fileSize := strconv.FormatInt(fileStat.Size(), 10)

	writer.Header().Set("Content-Disposition", "attachment; filename="+filename)
	writer.Header().Set("Content-Type", fileContentType)
	writer.Header().Set("Content-Length", fileSize)

	file.Seek(0, 0)
	buf := make([]byte, 512000)
	io.CopyBuffer(writer, file, buf)
	return
}
