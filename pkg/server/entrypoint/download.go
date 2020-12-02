package entrypoint

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Status struct {
	NumGoRoutines int
}


// pings each node storing the file, finds the node that has the least traffic and
// downloads from there 

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
	//c := make(chan int)
	channelMap := make(map[Node](chan int))
	for _, node := range nodes {
		channelMap[node] = make(chan int)
		go getStats(node, channelMap[node])
	}

	nodeStats := make(map[Node]int, len(nodes))
	for _, node := range nodes {
		stat := <-channelMap[node]
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

// toy stat to get a sense of which storage node is the busiest.  For simplicity, use the number of currently executing
// goroutines as a proxy for how busy a node is. 
func getStats(node Node, c chan int) {
	client := &http.Client{}
	reqStatus, _ := http.NewRequest("GET", fmt.Sprintf("http://%s:%s/status-check", node.Host, strconv.Itoa(node.Port)), nil)
	respStatusCheck, _ := client.Do(reqStatus)
	var data Status
	body, _ := ioutil.ReadAll(respStatusCheck.Body)
	json.Unmarshal(body, &data)
	c <- data.NumGoRoutines
}
