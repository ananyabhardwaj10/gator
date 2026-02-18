package main
import (
	"errors"
	"fmt"
	"time"
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/google/uuid"
	"github.com/ananyabhardwaj10/gator/internal/config"
	"github.com/ananyabhardwaj10/gator/internal/database"
)

type state struct {
	db *database.Queries
	cfg *config.Config  `json:"cfg"`
}

type command struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
}

type commands struct {
	CommandNames map[string]func(*state, command) error `json:"command_names"`
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return errors.New("please pass a username.")
	}
	
	_, err := s.db.GetUser(context.Background(), cmd.Args[0]) 
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user does not exist")
		}
		return err 
	}

	err = s.cfg.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Println("user has been set.")
	return nil 
}

func (c *commands) run(s *state, cmd command) error {
	f, exists := c.CommandNames[cmd.Name]
	if !exists {
		return errors.New("unknown command.") 
	}

	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.CommandNames[name] = f
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return errors.New("please pass a name.")
	}

	id := uuid.New()
	now := time.Now()
	name := cmd.Args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID: id,
		CreatedAt: now,
		UpdatedAt: now,
		Name: name,
	})
	_ = user 
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == "23505" {
			return fmt.Errorf("user already exists")
		}
		return err 
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return err 
	}
	return nil 
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUser(context.Background())
	if err != nil {
		return fmt.Errorf("Error encountered while resetting the database") 
	}

	return nil 
}

func handlerGetUsers(s *state, cmd command) error {
	usersList, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error encountered: %v", err)
	}

	for _, user := range usersList {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil 
}