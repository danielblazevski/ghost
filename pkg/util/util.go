package util

import (
	"net/http"
	"os"
	"strconv"
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
