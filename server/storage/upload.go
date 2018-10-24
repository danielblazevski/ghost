package storage

import (
	"fmt"
	"go-storage/util"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

// probably should have a "gateway server" that keeps track of which node is head, replica and tail.
// and forwards the write request to the head, if it is still alive.  If not alive, it will make next
// repliate the head.  Similary w/ tail.

// each node will still have <download> and <upload> routes only.

const baseLocation = "/Users/danielblazevski/gofun"

type NextClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// TO-DO: need to have a tail node so we know when to
// stop replicating.  Code only replicates once.
// Little annoying that next node writes to same path and deletes what
// the first node wrote if runnign on localhost

func HandleUploadStorage(writer http.ResponseWriter,
	request *http.Request,
	nextNodeClient NextClient) {

	query := request.URL.Query()
	dest := query.Get("dest")
	nextNode := query.Get("next")

	if dest == "" {
		http.Error(writer, "Post 'file' not specified in url", 400)
	}

	fileMainPath := fmt.Sprintf("%s/%s", baseLocation, dest)
	versionPtr, err := util.GetVersionedFile(writer, fileMainPath)
	if err != nil {
		log.Println("Could not create new file")
		http.Error(writer, "Could not fetch version.", 500)
		return
	}
	version := *versionPtr

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
	// turn this into a method "passToNextNode"
	// pass to neighboring node if next is not null
	if nextNode == "" {
		return
	}

	fileHeader := make([]byte, 512)
	file.Read(fileHeader)
	fileContentType := http.DetectContentType(fileHeader)

	fileStat, _ := file.Stat()
	fileSize := strconv.FormatInt(fileStat.Size(), 10)

	req, _ := http.NewRequest("POST", fmt.Sprintf("http://server2:%s/upload", nextNode), file)
	q := req.URL.Query()
	q.Add("dest", dest)
	req.URL.RawQuery = q.Encode()
	//Send the headers
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
