package util

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func GetHeaderInfo(file *os.File) (string, string) {
	//Get the Content-Type of the file
	//Create a buffer to store the header of the file in
	fileHeader := make([]byte, 512)
	//Copy the headers into the fileHeader buffer
	file.Read(fileHeader)
	//Get content type of file
	fileContentType := http.DetectContentType(fileHeader)
	//Get the file size
	fileStat, _ := file.Stat()                         //Get info from file
	fileSize := strconv.FormatInt(fileStat.Size(), 10) //Get file size as a string
	file.Seek(0, 0)
	return fileContentType, fileSize
}

func GetVersionedFile(writer http.ResponseWriter, fileMainPath string) (*int, error) {
	// check if file exists -- really just check if directory exists
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
		version = currentVersion
	} else {
		err := os.MkdirAll(fileMainPath, 0777)
		if err != nil {
			log.Println(err)
			http.Error(writer, "Could not write new directory", 500)
			return nil, err
		}
		version = 0
	}
	return &version, nil
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
