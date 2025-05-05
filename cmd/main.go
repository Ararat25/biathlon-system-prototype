package main

import (
	"biathlon-system-prototype/config"
	"biathlon-system-prototype/internal"
	"log"
	"os"
)

func main() {
	loadConfig, err := config.LoadConfig("../configExample.json")
	if err != nil {
		log.Fatalln(err)
	}

	file, err := os.Open("../src/eventsExample")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	err = internal.PrintReport(loadConfig, file)
	if err != nil {
		log.Fatalln(err)
	}
}
