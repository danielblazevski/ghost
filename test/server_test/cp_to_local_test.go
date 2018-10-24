package server_test

import (
	"go-storage/server/storage"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestToLocalGoodScenario(t *testing.T) {

	request := httptest.NewRequest("GET", "http://localhost:8080/download", nil)
	q := request.URL.Query()
	q.Add("file", "testfile.txt")
	request.URL.RawQuery = q.Encode()
	writer := httptest.NewRecorder()
	storage.HandleDownloadStorage(writer, request)

	resp := writer.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	expected_status_code := 200
	expected_body := "123\n"

	if resp.StatusCode != expected_status_code {
		t.Errorf("status code was not 200, it was %d", resp.StatusCode)
	}
	if string(body) != expected_body {
		t.Errorf("Body is not 123 as expected, it is: %s", string(body))
	}
}

func TestToLocalFileNotExist(t *testing.T) {

	request := httptest.NewRequest("GET", "http://localhost:8080/download", nil)
	q := request.URL.Query()
	q.Add("file", "invalid.txt")
	request.URL.RawQuery = q.Encode()
	writer := httptest.NewRecorder()
	storage.HandleDownloadStorage(writer, request)

	resp := writer.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	expected_status_code := 404
	expected_body := "File not found.\n"

	if resp.StatusCode != expected_status_code {
		t.Errorf("status code was not 400, it was %d", resp.StatusCode)
	}
	if string(body) != expected_body {
		t.Errorf("Body is not 'File not found.' as expected, it is: %s", string(body))
	}
}
