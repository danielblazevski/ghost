package fileversion

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

// methods for bookeeping of versions of files, with some thread safety.
// when a user uploads 'foo.txt' to 'doge://cloud/bar.txt' the path on the server is
//  /ghost/files/cloud/bar.txt/#10, where #10 is the version number.  

func makeNewFileVersion(writer http.ResponseWriter,
	fileMainPath string,
	filename string,
	version int) error {

	dirExists, err := exists(fileMainPath)
	if err != nil {
		return err
	}

	if !dirExists {
		err := os.MkdirAll(fileMainPath, 0777)
		if err != nil {
			log.Println(err)
			http.Error(writer, "Could not write new directory", 500)
			return err
		}
	}

	fileExists, _ := exists(fmt.Sprintf("%s/%s", fileMainPath, filename))
	if err != nil {
		return err
	}
	if !fileExists {
		log.Println("about to make new file:")
		log.Println(fmt.Sprintf("%s/%s", fileMainPath, filename))
		f, err := os.Create(fmt.Sprintf("%s/%s", fileMainPath, filename))
		if err != nil {
			log.Print(err)
			http.Error(writer, "Could not create new file", 500)
			return err
		}
		_, err = f.WriteString(fmt.Sprintf("%d\n", version))
		if err != nil {
			log.Print(err)
			http.Error(writer, "Could not write to file", 500)
			return err
		}
		f.Close()

	}
	return nil
}

func ReadVersionFromFilename(writer http.ResponseWriter, fileMainPath string, filename string) (*int, error) {
	exists, err := exists(fmt.Sprintf("%s/%s", fileMainPath, filename))
	if err != nil {
		return nil, err
	}
	if !exists {
		err = fmt.Errorf("file does not exist")
		return nil, err
	} else {
		f, _ := os.Open(fmt.Sprintf("%s/%s", fileMainPath, filename))
		f.Seek(0, 0)
		byteSlice := make([]byte, 16)
		bytesRead, _ := f.Read(byteSlice)
		s := byteSlice[:bytesRead]
		version, _ := strconv.Atoi(strings.Replace(string(s), "\n", "", 1))
		return &version, nil
	}
	return nil, err
}

// read, creates if empty.
func UpsertVersionFromFilename(writer http.ResponseWriter, fileMainPath string, filename string) (*int, error) {
	exists, err := exists(fileMainPath)
	if err != nil {
		return nil, err
	}
	if !exists {
		// Only concurrent threads that make it here will be blocked
		var mutex = &sync.Mutex{}
		mutex.Lock()
		if !exists {
			version := 0
			makeNewFileVersion(writer, fileMainPath, filename, version)
			mutex.Unlock()
			return &version, nil
		} else {
			// grab version.  race conditin edge case
			// TODO: deal with this.
		}
	} else {
		f, _ := os.Open(fmt.Sprintf("%s/%s", fileMainPath, filename))
		f.Seek(0, 0)
		byteSlice := make([]byte, 16)
		bytesRead, _ := f.Read(byteSlice)
		s := byteSlice[:bytesRead]
		version, _ := strconv.Atoi(strings.Replace(string(s), "\n", "", 1))
		return &version, nil
	}

	return nil, err
}

func UpdatevVersionFromFilename(writer http.ResponseWriter,
	fileMainPath string,
	filename string,
	version int) (*int, error) {

	exists, err := exists(fmt.Sprintf("%s/%s", fileMainPath, filename))
	if err != nil {
		log.Println("exists err!")
		return nil, err
	}

	if !exists {
		// Only concurrent threads that make it here will be blocked
		// TODO: really only need to lock on the file, not this block of code
		var mutex = &sync.Mutex{}
		mutex.Lock()
		if !exists {
			// make a new file with the input version.  Cannot assume that version is 1
			makeNewFileVersion(writer, fileMainPath, filename, version)
			mutex.Unlock()
			return &version, err
		} else {
			// grab version. race condition edge case
			// TODO: deal with this
		}
	} else {
		// main case (most often, not new file).  grab version and update
		f, _ := os.OpenFile(fmt.Sprintf("%s/%s", fileMainPath, filename), os.O_RDWR, 0755)
		syscall.Flock(int(f.Fd()), syscall.LOCK_EX)
		f.Seek(0, 0)
		byteSlice := make([]byte, 16)
		bytesRead, _ := f.Read(byteSlice)
		s := byteSlice[:bytesRead]
		version, _ := strconv.Atoi(strings.Replace(string(s), "\n", "", 1))
		log.Println(fmt.Sprintf("version reading in %d for file %s", version, filename))
		versionNew := version + 1
		f.Seek(0, 0)
		versionNewString := strconv.Itoa(versionNew) + "\n"
		_, err = f.Write([]byte(versionNewString))
		if err != nil {
			log.Println("Could not write new line to file!")
			log.Println(err)
		}
		syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		f.Close()
		return &versionNew, err
	}
	return nil, err
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
