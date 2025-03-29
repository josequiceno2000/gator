package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/josequiceno2000/gator/internal/config"
	"github.com/josequiceno2000/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	DB *database.Queries
	CfgPointer *config.Config
}

type command struct {
	Name string
	Arguments []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("username is required")
	}

	username := cmd.Arguments[0]

	// Check if the user exists in the database
	_, err := s.DB.GetUser(context.Background(), username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("login: user '%s' does not exist", username)
		}
		return fmt.Errorf("login: failed to check user instance: %w", err)
	}


	err = s.CfgPointer.SetUser(username)
	if err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}
	
	fmt.Printf("Username has been set to %v", username)

	return nil
}

type commands struct {
	Handlers map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	if c.Handlers == nil {
		c.Handlers = make(map[string]func(*state, command) error)
	}
	c.Handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.Handlers[cmd.Name]
	if !ok {
		return fmt.Errorf("command not found: %s", cmd.Name)
	}
	return handler(s, cmd)
}