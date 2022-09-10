package main

import (
	"log"
	"os"

	"github.com/mqasimsarfraz/imagesync"
)

func main() {
	if err := imagesync.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
