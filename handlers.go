package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/awbalessa/gator/internal/config"
	"github.com/awbalessa/gator/internal/database"
	"github.com/awbalessa/gator/internal/feeds"
	"github.com/google/uuid"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("username required")
	}

	if _, err := s.db.GetUser(context.Background(), cmd.args[0]); errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("user does not exist: %v", err)
	} else if err != nil {
		return fmt.Errorf("error getting user: %v", err)
	}

	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User has been set to %s\n", s.cfg.GetUser())
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("username required")
	}
	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err == nil {
		return fmt.Errorf("User %s already registered\n", cmd.args[0])
	} else if errors.Is(err, sql.ErrNoRows) {
		userParams := database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.args[0],
		}
		newUser, err := s.db.CreateUser(context.Background(), userParams)
		if err != nil {
			return fmt.Errorf("error creating user: %w", err)
		}
		s.cfg.SetUser(newUser.Name)
		fmt.Printf("User has been created:\n %+v\n", newUser)
		return nil
	} else {
		return fmt.Errorf("error occurred: %v", err)
	}
}

func handlerReset(s *state, _ command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		return fmt.Errorf("error resetting database: %v", err)
	}
	s.cfg.SetUser("")
	fmt.Println("Database reset successfully")
	return nil
}

func handleUsers(s *state, _ command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error getting users: %v", err)
	}
	for _, user := range users {
		if user.Name == s.cfg.CurrentUsername {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func handleAgg(s *state, _ command) error {
	exampleURL := "https://www.wagslane.dev/index.xml"
	feed, err := feeds.FetchFeed(context.Background(), exampleURL)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", *feed)
	return nil
}

func handleAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("feed name and url required")
	}
	if s.cfg.GetUser() == "" {
		return fmt.Errorf("you must be logged in to add a feed")
	}
	user, err := s.db.GetUser(context.Background(), s.cfg.GetUser())
	if err != nil {
		return fmt.Errorf("error getting user: %v", err)
	}

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return fmt.Errorf("error creating feed: %v", err)
	}

	fmt.Printf("Feed added successfully:\n %+v\n", feed)
	return nil
}

func handleFeeds(s *state, _ command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds: %v", err)
	}

	for i := range feeds {
		user, err := s.db.GetUserByID(context.Background(), feeds[i].UserID)
		if err != nil {
			return fmt.Errorf("error getting user by ID %v", err)
		}
		fmt.Printf("Feed #%d:\nName: %s\nURL: %s\nOwner: %s\n", i+1, feeds[i].Name, feeds[i].Url, user.Name)
	}
	return nil
}

func handleFollow(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("url required")
	}

	feed, err := s.db.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("error getting feed: %v", err)
	}
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
