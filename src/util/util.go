package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func GetVersionedFile(writer http.ResponseWriter, fileMainPath string) (*int, error) {
	// check if file exists -- really just check if directory exists
	//fileMainPath := fmt.Sprintf("%s/%s", baseLocation, dest)
	exists, err := exists(fileMainPath)
	if err != nil {
		return nil, err
	}

	var version int
	if exists {
		files, err := ioutil.ReadDir(fileMainPath)
		if err != nil {
			http.Error(writer, "could not read in directory", 500)
			return nil, err
		}
		latestfile := files[len(files)-1]
		splitted := strings.Split(latestfile.Name(), "#")
		currentVersion, _ := strconv.Atoi(splitted[len(splitted)-1])
		version = currentVersion + 1
	} else {
		err := os.MkdirAll(fileMainPath, 0777)
		if err != nil {
			log.Println(err)
			http.Error(writer, "Could not write new directory", 500)
			return nil, err
		}
		version = 1
	}
	fmt.Println(fmt.Sprintf("%s/%s", fileMainPath, strconv.Itoa(version)))
	return &version, nil
	//return os.Create(fmt.Sprintf("%s/%s", fileMainPath, strconv.Itoa(version)))
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
