package main

import (
	"github.com/kristofaranyos/tech-challenge-time/cmd/core"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting time tracker service.")

	c, err := core.NewApiCore()
	if err != nil {
		log.Fatalf("Couldn't start server: %v", err)
	}
	defer c.Close()

	if err := c.Run(); err != nil {
		log.Fatalf("Web server error: %v", err)
	}
}
