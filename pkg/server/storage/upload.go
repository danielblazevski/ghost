package storage

import (
	"fmt"
	"ghost/pkg/fileversion"
	"ghost/pkg/util"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

// can use the interface for testing.  Only out of the box call.
// to really test, should not do any file system IO.
// should look into testing complex void Go methods w/o mocking...

type NextClient interface {
	Do(req *http.Request) (*http.Response, error)
}

//  when a user uploads 'foo.txt' to 'doge://cloud/bar.txt' the path on the server is
//  /ghost/files/cloud/bar.txt/#10, where #10 is the version number.
//  we also make files /ghost/files/cloud/bar.txt/latest-upload-start and 
// /ghost/files/cloud/bar.txt/latest-upload-complete to keep track of version numbers 
// to ensure atomic writes and safe deletion of older versions

func HandleUploadStorage(writer http.ResponseWriter,
	request *http.Request,
	nextNodeClient NextClient,
	nextService string,
	nextPort string,
	baseLocation string) {

	query := request.URL.Query()
	dest := query.Get("dest")

	if dest == "" {
		http.Error(writer, "Post 'file' not specified in url", 400)
	}

	log.Println("upload request to: " + dest)

	fileMainPath := fmt.Sprintf("%s/%s", baseLocation, dest)

	// TODO: update to read in latestUploadStart in file and increment
	versionPrevPtr, err := fileversion.ReadOrCreateVersionFromFilename(writer, fileMainPath, "latest-upload-start")
	if err != nil {
		log.Println("Could not create new file")
		http.Error(writer, "Could not fetch version.", 500)
		return
	}
	versionPrev := *versionPrevPtr
	versionPtr, _ := fileversion.UpdatevVersionFromFilename(writer, fileMainPath, "latest-upload-start", versionPrev)
	version := *versionPtr + 1

	// TODO since using version files don't need to explicity name files based on version
	file, err := os.Create(fmt.Sprintf("%s/#%s", fileMainPath, strconv.Itoa(version)))
	if err != nil {
		log.Println("Could not create new file")
		http.Error(writer, "Could not create file.", 500)
		return
	}
	defer file.Close()

	buf := make([]byte, 512000)
	io.CopyBuffer(file, request.Body, buf)

	if version > 2 {
		err = os.Remove(fmt.Sprintf("%s/#%s", fileMainPath, strconv.Itoa(version-1)))
		if err != nil {
			log.Println("Could not delete old file!")
			http.Error(writer, "Could not delete old version.", 500)
			return
		}
	}
	_, err = fileversion.UpdatevVersionFromFilename(writer,
		fileMainPath,
		"latest-upload-complete",
		version)
	// TODO: update lastUploadComplete in file

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
		return
	}

}
