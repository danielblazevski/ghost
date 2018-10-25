package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const prefix = "doge://"

type Node struct {
	Host string
	Port int
}

func Cp(from string, to string) error {
	fromRemote := strings.HasPrefix(from, prefix)
	toRemote := strings.HasPrefix(to, prefix)

	if !fromRemote && toRemote {
		return cpLocalToRemote(from, to)
	} else if fromRemote && !toRemote {
		return cpRemoteToLocal(from, to)
	} else {
		return errors.New("could not cp the files, either both remote or both local")
	}
}

func cpLocalToRemote(local string, remote string) error {
	// first figure out which node to hit
	client1 := &http.Client{}
	fmt.Println(local)

	// hits a "gateway" server to determine which server to upload to
	// TO-DO: Change name of "gateway"
	req1, _ := http.NewRequest("GET", "http://localhost:8090/upload", nil)
	q1 := req1.URL.Query()
	q1.Add("file", strings.Replace(remote, prefix, "", -1))
	req1.URL.RawQuery = q1.Encode()
	resp1, _ := client1.Do(req1)
	var data Node
	body, _ := ioutil.ReadAll(resp1.Body)
	json.Unmarshal(body, &data)

	localFile, err := os.Open(local)
	if err != nil {
		return err
	}
	fmt.Println(data)
	//Get the Content-Type of the file
	//Create a buffer to store the header of the file in
	fileHeader := make([]byte, 512)
	//Copy the headers into the fileHeader buffer
	localFile.Read(fileHeader)
	//Get content type of file
	fileContentType := http.DetectContentType(fileHeader)

	//Get the file size
	fileStat, _ := localFile.Stat()                    //Get info from file
	fileSize := strconv.FormatInt(fileStat.Size(), 10) //Get file size as a string

	client := &http.Client{}

	req, _ := http.NewRequest("POST", fmt.Sprintf("http://localhost:%d/upload", data.Port), localFile)
	q := req.URL.Query()
	q.Add("dest", strings.Replace(remote, prefix, "", -1))
	q.Add("next", "8081")
	req.URL.RawQuery = q.Encode()
	//Send the headers
	req.Header.Set("Content-Disposition", "attachment; local="+local)
	req.Header.Set("Content-Type", fileContentType)
	req.Header.Set("Content-Length", fileSize)
	localFile.Seek(0, 0)
	resp, err := client.Do(req)
	log.Println(resp)
	if err != nil {
		return err
	}
	return nil
}

func cpRemoteToLocal(local string, remote string) error {
	// first figure out which node to hit
	client1 := &http.Client{}
	req1, _ := http.NewRequest("GET", "http://localhost:8090/download", nil)
	q1 := req1.URL.Query()
	q1.Add("file", strings.Replace(remote, prefix, "", -1))
	req1.URL.RawQuery = q1.Encode()
	resp1, _ := client1.Do(req1)
	var data Node
	body, _ := ioutil.ReadAll(resp1.Body)
	json.Unmarshal(body, &data)

	// now hit that node
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/download", data.Port), nil)
	q := req.URL.Query()
	q.Add("file", remote)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		bytes, _ := ioutil.ReadAll(resp.Body)
		return errors.New(fmt.Sprintf("Failed to fetch file, reason: %s ", string(bytes)))
	}

	defer resp.Body.Close()

	f, err := os.Create(local)
	if err != nil {
		return err
	}
	buf := make([]byte, 512000)
	io.CopyBuffer(f, resp.Body, buf)
	stats, _ := f.Stat()
	log.Printf("Saved remote file %s to %s \nbytes downloaded: %d \n", remote, f.Name(), stats.Size())
	return nil
}