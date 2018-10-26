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
	"strings"
	"util"
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
	log.Printf("copy local file %s to remote %s", local, remote)
	// first figure out which node to hit
	clientEntrypoint := &http.Client{}

	// hits a "entrypoint" server to determine which server to upload to
	reqEntryPoint, _ := http.NewRequest("GET", "http://localhost:8090/upload", nil)
	data := makeRequestEntrypoint(clientEntrypoint, reqEntryPoint, remote)

	localFile, err := os.Open(local)
	defer localFile.Close()
	if err != nil {
		return err
	}

	// upload to node obtained from entrypoint.
	client := &http.Client{}
	reqStroage, _ := http.NewRequest("POST", fmt.Sprintf("http://localhost:%d/upload", data.Port), localFile)
	queryStorage := reqStroage.URL.Query()
	queryStorage.Add("dest", strings.Replace(remote, prefix, "", -1))
	reqStroage.URL.RawQuery = queryStorage.Encode()

	fileContentType, fileSize := util.GetHeaderInfo(localFile)
	reqStroage.Header.Set("Content-Disposition", "attachment; local="+local)
	reqStroage.Header.Set("Content-Type", fileContentType)
	reqStroage.Header.Set("Content-Length", fileSize)
	resp, err := client.Do(reqStroage)
	log.Println(resp)
	if err != nil {
		return err
	}
	return nil
}

func cpRemoteToLocal(remote string, local string) error {
	log.Printf("copy remote file %s to local %s", local, remote)

	// first figure out which node to hit
	clientEntrypoint := &http.Client{}
	reqEntryPoint, _ := http.NewRequest("GET", "http://localhost:8090/download", nil)
	data := makeRequestEntrypoint(clientEntrypoint, reqEntryPoint, remote)

	// now hit that node
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/download", data.Port), nil)
	q := req.URL.Query()
	q.Add("file", strings.Replace(remote, prefix, "", -1))
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
		log.Println("could not create local file!")
		return err
	}
	buf := make([]byte, 512000)
	io.CopyBuffer(f, resp.Body, buf)
	stats, _ := f.Stat()
	log.Printf("Saved remote file %s to %s \nbytes downloaded: %d \n", remote, f.Name(), stats.Size())
	return nil
}

func makeRequestEntrypoint(client *http.Client, req *http.Request, remote string) Node {
	queryEntryPoint := req.URL.Query()
	queryEntryPoint.Add("file", strings.Replace(remote, prefix, "", -1))
	req.URL.RawQuery = queryEntryPoint.Encode()
	respEntryPoint, _ := client.Do(req)
	var data Node
	body, _ := ioutil.ReadAll(respEntryPoint.Body)
	json.Unmarshal(body, &data)
	return data
}
