package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/olleboll/images/api"
	"github.com/olleboll/images/store"
)

func main() {

	godotenv.Load()

	// Dirty shortcut to make the database ready for connections
	if os.Getenv("SLOW_START") == "true" {
		time.Sleep(2 * time.Second)
	}

	imageStore, err := store.Connect()

	if err != nil {
		log.Fatal("Could not connect to db")
		return
	}

	api.Run(imageStore)
}
