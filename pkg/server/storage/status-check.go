package storage

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
)

type Status struct {
	NumGoRoutines int
}

func HandleStatusCheck(writer http.ResponseWriter, request *http.Request) {
	numGoRoutines := runtime.NumGoroutine()
	log.Printf("number of goroutines %d", numGoRoutines)

	status := Status{NumGoRoutines: numGoRoutines}
	js, err := json.Marshal(status)
	if err != nil {
		http.Error(writer, "could not form json response", 500)
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}
