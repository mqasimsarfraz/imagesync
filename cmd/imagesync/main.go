package main

import (
	imagesync "github.com/MQasimSarfraz/imagesync"
	"log"
)

func main() {
	err := imagesync.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
