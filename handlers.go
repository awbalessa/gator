package main

import (
	"fmt"

	"github.com/awbalessa/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("username required")
	}

	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Username has been set to %s\n", s.cfg.GetUser())
	return nil
}

type commands struct {
	cmdToHandler map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmdToHandler[name] = f
}
func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.cmdToHandler[cmd.name]
	if !ok {
		return fmt.Errorf("command not registered")
	}
	return handler(s, cmd)
}
