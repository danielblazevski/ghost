package main

import (
	"go-storage/client"
	"log"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		log.Println("No action specified!")
		return
	}
	action := os.Args[1]

	switch action {
	case "cp":
		if len(os.Args) < 4 {
			log.Println("Source and destination files not specified!")
			return
		}
		from := os.Args[2]
		to := os.Args[3]
		client.Cp(from, to)
	default:
		log.Printf("%s is not a valid action! \n", action)
	}

}
