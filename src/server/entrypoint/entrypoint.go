package entrypoint

import (
	"encoding/json"
	"fmt"
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

	fmt.Println("Client requests: " + filename)
	chainIndex := some_super_cool_hash(filename)
	nodes := chains[chainIndex].nodes
	roundRobin = int(math.Mod(float64(roundRobin+1), float64(len(nodes))))
	fmt.Println(roundRobin)
	node := nodes[roundRobin]
	fmt.Println(node.Host)
	fmt.Println(node.Port)
	fmt.Println(node)
	js, err := json.Marshal(node)
	fmt.Println(js)
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
	fmt.Println("Client requests: " + filename)
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
