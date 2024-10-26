package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
)

var (
	db     *MongoDB = &MongoDB{}
	config *Config  = DefaultConfig
	// instanceID    uint16   = GetInstanceID()
	// instanceCount uint16   = GetInstanceCount()
)

func init() {
	var err error

	if err = config.ReadFile("config.yml"); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			panic(err)
		}

		if err = config.WriteFile("config.yml"); err != nil {
			log.Fatalf("Failed to write config file: %v", err)
		}
	}

	if err := db.Connect(config.MongoDB); err != nil {
		panic(err)
	}

	log.Println("Successfully connected to MongoDB")
}

func main() {
	defer db.Close()

	go StartRunner()

	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	<-s
}
