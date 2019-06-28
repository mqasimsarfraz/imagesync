package main

import (
	imagesync "github.com/MQasimSarfraz/imagesync"
	"log"
	"os"
)

func main() {
	err := imagesync.Execute()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
