package main

import (
	"log"
	"os"

	imagesync "github.com/MQasimSarfraz/imagesync"
)

func main() {
	err := imagesync.Execute()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
