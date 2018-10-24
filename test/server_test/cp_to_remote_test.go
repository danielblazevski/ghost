package server_test

import (
	"go-storage/server/storage"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type StubClientSuccess struct{}

func (sc StubClientSuccess) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{}, nil
}

// this works, but saves in baseLocation = "/Users/danielblazevski/gofun"
// To-DO clean up so that can save in current directory, maybe have baseLocation
// be a param to the cp to remote
// even better, shouldn't save

func TestToRemoteGoodScenario(t *testing.T) {

	localFile, _ := os.Open("testfile.txt")
	request := httptest.NewRequest("POST", "http://localhost:8080/upload", localFile)
	q := request.URL.Query()
	q.Add("dest", "testfile-dir.txt")
	q.Add("next", "8081")
	request.URL.RawQuery = q.Encode()

	writer := httptest.NewRecorder()
	stub_client := &StubClientSuccess{}
	storage.HandleUploadStorage(writer, request, stub_client)

	resp := writer.Result()
	expected_status_code := 200
	if resp.StatusCode != expected_status_code {
		t.Errorf("status code was not 200, it was %d", resp.StatusCode)
	}
}
