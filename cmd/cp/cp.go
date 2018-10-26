package main

import (
	"ghost/client"
	"log"
	"os"
)

func main() {

	if len(os.Args) < 3 {
		log.Println("Source and destination files not specified!")
		return
	}
	from := os.Args[1]
	to := os.Args[2]
	err := client.Cp(from, to)
	if err != nil {
		log.Println(err)
	}

}
