package main 

import (
	"context"
	"errors"
	"database/sql"
	"os"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/josequiceno2000/gator/internal/config"
	"github.com/josequiceno2000/gator/internal/database"
	_ "github.com/lib/pq"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Arguments) == 0 {
		return errors.New("register command requires a username argument")
	}

	username := cmd.Arguments[0]
	userID := uuid.New()
	now := time.Now().UTC()

	user, err := s.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID: userID,
		CreatedAt: now,
		UpdatedAt: now,
		Name: username,
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user with name '%s' already exists", username)
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	err = s.CfgPointer.SetUser(username)
	if err != nil {
		return fmt.Errorf("failed to set current user in config: %w", err)
	}

	fmt.Printf("User '%s' registered successfully.\n", username)
	log.Printf("Registered user: %+v\n", user)

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.DB.DeletaAllUsers((context.Background()))
	if err != nil {
		return fmt.Errorf("reset: failed to delete all users: %w", err)
	}
	fmt.Println("reset: all users deleted successfully")
	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	// Open db connection
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("Error opening db connectio: %v", err)
	}
	defer db.Close()

	// Create database queries instance
	dbQueries := database.New(db)

	appState := state{DB:dbQueries, CfgPointer: &cfg}
	
	cmdRegistry := commands{}
	cmdRegistry.register("login", handlerLogin)
	cmdRegistry.register("register", handlerRegister)
	cmdRegistry.register("reset", handlerReset)

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