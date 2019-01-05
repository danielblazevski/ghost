package storage

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDownloadHandler(t *testing.T) {

	reqStroage, err := http.NewRequest("GET", "http://localhost/download", nil)

	queryStorage := reqStroage.URL.Query()
	queryStorage.Add("file", "test_file.txt")
	reqStroage.URL.RawQuery = queryStorage.Encode()

	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	// cleaner way to read test resources?
	HandleDownloadStorage(res, reqStroage, "../../../test")

	exp := "some text."
	act := res.Body.String()
	if exp != act {
		t.Fatalf("Expected = %s was not actual =  %s", exp, act)
	}
}
