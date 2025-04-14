package main

import (
	"log"
	"os"

	"github.com/awbalessa/gator/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("error: command required")
	}
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}
	s := &state{
		cfg: cfg,
	}
	cmds := commands{
		cmdToHandler: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	if err = cmds.run(s, cmd); err != nil {
		log.Fatalf("failed to run command: %v", err)
	}
}
