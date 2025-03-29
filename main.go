package main

import (
	"os"
	"fmt"
	"log"
	"github.com/josequiceno2000/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	appState := state{CfgPointer: &cfg}
	
	cmdRegistry := commands{}
	cmdRegistry.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("Error: not enough arguments provided")
		os.Exit(1)
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	cmd := command{Name: cmdName, Arguments: cmdArgs}

	err = cmdRegistry.run(&appState, cmd)
	if err != nil {
		log.Fatalf("Command error: %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading updated config: %v", err)
	}

	fmt.Println("Config after command:")
	fmt.Printf("DB URL: %s\n", cfg.DBURL)
	fmt.Printf("Current User: %s\n", cfg.CurrentUsername)
}