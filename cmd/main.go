package main

import (
	"biathlon-system-prototype/config"
	"biathlon-system-prototype/internal"
	"log"
	"os"
)

func init() {
	config.LoadEnvVars()
}

func main() {
	loadConfig, err := config.LoadConfig(config.ConfigPath)
	if err != nil {
		log.Fatalln(err)
	}

	file, err := os.Open(config.EventsPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	err = internal.PrintReport(loadConfig, file)
	if err != nil {
		log.Fatalln(err)
	}
}
