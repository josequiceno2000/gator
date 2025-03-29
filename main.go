package main

import (
	"fmt"
	"log"
	"github.com/josequiceno2000/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	err = cfg.SetUser("jose")
	if err != nil {
		log.Fatalf("Error setting user: %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("Error reading updated config: %v", err)
	}

	fmt.Println("Config after setting user:")
	fmt.Printf("DB URL: %s\n", cfg.DBURL)
	fmt.Printf("Current User: %s\n", cfg.CurrentUsername)
}