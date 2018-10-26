package storage

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"util"
)

func HandleDownloadStorage(writer http.ResponseWriter, request *http.Request) {
	baseLocation := "/ghost/files"

	//First of check if Get is set in the URL
	filename := request.URL.Query().Get("file")
	if filename == "" {
		//Get not set, send a 400 bad request
		http.Error(writer, "Get 'file' not specified in url.", 400)
		return
	}
	log.Println("Download request: " + filename)
	mainPath := fmt.Sprintf("%s/%s", baseLocation, filename)

	versionPtr, err := util.GetVersionedFile(writer, mainPath)
	if err != nil {
		log.Println("Could not create new file")
		http.Error(writer, "Could not fetch version.", 500)
		return
	}
	version := *versionPtr
	versionedFile := fmt.Sprintf("%s/#%s", mainPath, strconv.Itoa(version))

	file, err := os.Open(versionedFile)
	defer file.Close()
	if err != nil {
		//File not found, send 404
		log.Printf("could not open file %s", versionedFile)
		http.Error(writer, "File not found.", 404)
		return
	}

	fileContentType, fileSize := util.GetHeaderInfo(file)
	writer.Header().Set("Content-Disposition", "attachment; filename="+filename)
	writer.Header().Set("Content-Type", fileContentType)
	writer.Header().Set("Content-Length", fileSize)

	file.Seek(0, 0)
	buf := make([]byte, 512000)
	io.CopyBuffer(writer, file, buf)
	return
}
