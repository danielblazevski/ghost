package entrypoint

import (
	"encoding/json"
	"log"
	"net/http"
)

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
