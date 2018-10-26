package entrypoint

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
)

func some_super_cool_hash(s string) int {
	return int(math.Mod(float64(len(s)), float64(num_chains)))
}

type NodeStat struct {
	N    Node
	Stat int
}

type Status struct {
	NumGoRoutines int
}

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

	// ping each node in parallel and choose which one
	c := make(chan int)
	for _, node := range nodes {
		go getStats(node, c)
	}

	nodeStats := make(map[Node]int, len(nodes))
	for _, node := range nodes {
		stat := <-c
		nodeStats[node] = stat
	}
	// get min node
	var minStat = 10000000
	var minNode = nodes[0]
	for node, stat := range nodeStats {
		if stat < minStat {
			minStat = stat
			minNode = node
		}
	}

	js, err := json.Marshal(minNode)
	if err != nil {
		http.Error(writer, "could not form json response", 500)
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}

func getStats(node Node, c chan int) {
	client := &http.Client{}
	reqStatus, _ := http.NewRequest("GET", fmt.Sprintf("http://%s:%s/status-check", node.Host, strconv.Itoa(node.Port)), nil)
	respStatusCheck, _ := client.Do(reqStatus)
	var data Status
	body, _ := ioutil.ReadAll(respStatusCheck.Body)
	json.Unmarshal(body, &data)
	c <- data.NumGoRoutines
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
