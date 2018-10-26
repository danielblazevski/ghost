package entrypoint

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
)

func some_super_cool_hash(s string) int {
	return int(math.Mod(float64(len(s)), float64(num_chains)))
}

var roundRobin = 0

// these return node for client to ping
func HandleDownload(writer http.ResponseWriter, request *http.Request) {
	filename := request.URL.Query().Get("file")
	if filename == "" {
		http.Error(writer, "Get 'file' not specified in url.", 400)
		return
	}

	log.Println("Client requests to download: " + filename)
	chainIndex := some_super_cool_hash(filename)
	nodes := chains[chainIndex].nodes
	roundRobin = int(math.Mod(float64(roundRobin+1), float64(len(nodes))))
	node := nodes[roundRobin]
	js, err := json.Marshal(node)
	if err != nil {
		http.Error(writer, "could not form json response", 500)
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}

func HandleUpload(writer http.ResponseWriter, request *http.Request) {
	filename := request.URL.Query().Get("file")
	if filename == "" {
		http.Error(writer, "Get 'file' not specified in url.", 400)
		return
	}
	log.Println("Client request to upload: " + filename)
	chainIndex := some_super_cool_hash(filename)
	nodes := chains[chainIndex].nodes
	node := nodes[0]
	js, err := json.Marshal(node)
	if err != nil {
		http.Error(writer, "could not form json response", 500)
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}
