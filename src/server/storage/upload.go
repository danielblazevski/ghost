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

type NextClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func HandleUploadStorage(writer http.ResponseWriter,
	request *http.Request,
	nextNodeClient NextClient,
	nextService string,
	nextPort string) {

	baseLocation := "/ghost/files"
	query := request.URL.Query()
	dest := query.Get("dest")

	if dest == "" {
		http.Error(writer, "Post 'file' not specified in url", 400)
	}

	log.Println("upload request to: " + dest)

	fileMainPath := fmt.Sprintf("%s/%s", baseLocation, dest)
	versionPtr, err := util.GetVersionedFile(writer, fileMainPath)
	if err != nil {
		log.Println("Could not create new file")
		http.Error(writer, "Could not fetch version.", 500)
		return
	}
	version := *versionPtr + 1

	file, err := os.Create(fmt.Sprintf("%s/#%s", fileMainPath, strconv.Itoa(version)))
	if err != nil {
		log.Println("Could not create new file")
		http.Error(writer, "Could not create file.", 500)
		return
	}
	defer file.Close()

	buf := make([]byte, 512000)
	io.CopyBuffer(file, request.Body, buf)

	if version > 1 {
		err = os.Remove(fmt.Sprintf("%s/#%s", fileMainPath, strconv.Itoa(version-1)))
		if err != nil {
			log.Println("Could not delete old file!")
			http.Error(writer, "Could not delete old version.", 500)
			return
		}
	}

	if nextService == "none" {
		return
	}

	req, _ := http.NewRequest("POST", fmt.Sprintf("http://%s:%s/upload", nextService, nextPort), file)
	q := req.URL.Query()
	q.Add("dest", dest)
	req.URL.RawQuery = q.Encode()
	fileContentType, fileSize := util.GetHeaderInfo(file)
	req.Header.Set("Content-Disposition", "attachment; filename="+dest)
	req.Header.Set("Content-Type", fileContentType)
	req.Header.Set("Content-Length", fileSize)
	file.Seek(0, 0)
	_, err = nextNodeClient.Do(req)
	if err != nil {
		log.Println(err)
		http.Error(writer, "Could not replicate file", 500)
	}
	return
}
