package main

import (
	"fmt"
	"github.com/josequiceno2000/gator/internal/config"
)

type state struct {
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
	err := s.CfgPointer.SetUser(username)
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