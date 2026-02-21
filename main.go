package main

import _ "github.com/lib/pq"
import (
	"fmt"
	"log"
	"os"
	"context"

	"database/sql"
	"github.com/ananyabhardwaj10/gator/internal/config"
	"github.com/ananyabhardwaj10/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user_name := s.cfg.CurrentUserName
		user, err := s.db.GetUser(context.Background(), user_name)
		if err != nil {
			return err 
		}
		err = handler(s, cmd, user)
		if err != nil {
			return err 
		}

		return nil 
	}
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)
	
	s := state{
		db:dbQueries,
		cfg:&cfg, 
	}

	cmds := commands{
		CommandNames: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	args := os.Args

	if len(args) < 2 {
		fmt.Println("not enough arguments. exiting the program")
		os.Exit(1)
	} 
	
	cmd := command{
		Name: args[1],
		Args: args[2:],
	}

	err = cmds.run(&s, cmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}