package main

import (
	"fmt"
	"log"

	"github.com/awbalessa/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	err = cfg.SetUser("aziz")
	if err != nil {
		log.Fatalf("error setting user to aziz: %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("failed to read config a second time: %v", err)
	}

	fmt.Println(cfg)
}
