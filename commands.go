package main
import (
	"errors"
	"fmt"

	"github.com/ananyabhardwaj10/gator/internal/config"
)

type state struct {
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
	
	err := s.cfg.SetUser(cmd.Args[0])
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