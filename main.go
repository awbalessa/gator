package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/awbalessa/gator/internal/config"
	"github.com/awbalessa/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("error: command required")
	}
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to open database connection: %v", err)
	}
	if db == nil {
		log.Fatalf("database connection is nil")
	}
	dbQueries := database.New(db)
	s := &state{
		cfg: cfg,
		db:  dbQueries,
	}
	cmds := commands{
		cmdToHandler: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handleLogin)
	cmds.register("register", handleRegister)
	cmds.register("reset", handleReset)
	cmds.register("users", handleUsers)
	cmds.register("agg", handleAgg)
	cmds.register("addfeed", middlewareLoggedIn(handleAddFeed))
	cmds.register("feeds", handleFeeds)
	cmds.register("follow", middlewareLoggedIn(handleFollow))
	cmds.register("following", middlewareLoggedIn(handleFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handleUnfollow))
	cmds.register("browse", middlewareLoggedIn(handleBrowse))
	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	if err = cmds.run(s, cmd); err != nil {
		log.Fatalf("failed to run command: %v", err)
	}
}
